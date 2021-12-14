package core

import (
	"fmt"
	"strings"
)

type RolesSchema map[string]RoleType

// Model is casbin model schema.
type Model struct {
	policy   AssertionSchema
	request  AssertionSchema
	roles    RolesSchema
	effector Effector
	matchers Matchers
}

// NewModel constructes a new Model.
func NewModel(
	policy, request AssertionSchema,
	roles RolesSchema,
	effector Effector,
	matchers Matchers,
) (m *Model) {
	m = &Model{
		policy:   policy,
		request:  request,
		roles:    roles,
		effector: effector,
		matchers: matchers,
	}
	return
}

// Policy returns the policy assertion schema.
func (m *Model) Policy() AssertionSchema {
	return m.policy
}

// Request returns the request assertion schema.
func (m *Model) Request() AssertionSchema {
	return m.request
}

// Request returns the matcherx.
func (m *Model) Matchers() Matchers {
	return m.matchers
}

// Effector returns the effector.
func (m *Model) Effector() Effector {
	return m.effector
}

// Load load and returns the policies and role mappings by an assertion iterator.
//
// Load would check the loaded data if it is matched with the model.
func (m *Model) Load(itr AssertionIterator) (
	policies Policies, roleMappings RoleMappings, err error,
) {
	policies = make([]Assertion, 0)
	roleMappings = make(map[string]*RoleMapping)

	for key := range m.roles {
		roleMappings[key] = NewRoleMapping(key)
	}

	for {
		ok, key, vals := itr.Next()
		if !ok {
			err = itr.Error()
			break
		}

		if key == "p" {
			o, _err := m.policy.CreateAssertion(vals)
			if _err != nil {
				err = fmt.Errorf("load policy: %w", _err)
				return
			}
			policies = append(policies, o)
		} else {
			rType, ok := m.roles[key]
			if !ok {
				err = fmt.Errorf("unknown assertion key: %s", key)
				return
			}
			rg := roleMappings[key]

			var src, dst, domain string
			if rType == RoleTypeWithDomain {
				if len(vals) != 3 {
					err = fmt.Errorf("invalid role assertion: %s: %s", key, strings.Join(vals, ","))
					return
				}
				src = vals[0]
				dst = vals[1]
				domain = vals[2]
			} else {
				if len(vals) != 2 {
					err = fmt.Errorf("invalid role assertion: %s: %s", key, strings.Join(vals, ","))
					return
				}
				src = vals[0]
				dst = vals[1]
			}

			rg.graph.AddEdge(src, dst, domain)
		}
	}

	return
}
