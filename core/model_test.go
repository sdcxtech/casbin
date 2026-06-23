package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sdcxtech/casbin/core"
)

type mockAssertionIterator struct {
	data  [][]string
	index int
}

func (it *mockAssertionIterator) Next() (ok bool, key string, vals []string) {
	ok = true
	if it.index >= len(it.data) {
		ok = false

		return
	}

	line := it.data[it.index]
	key = line[0]
	vals = line[1:]
	it.index++

	return
}

func (it *mockAssertionIterator) Error() (err error) {
	return nil
}

func TestModelLoad(t *testing.T) {
	const (
		resource = "order"
		user     = "alice"
	)

	policy, err := core.NewAssertionSchema("sub, obj, act")
	require.NoError(t, err)

	request, err := core.NewAssertionSchema("sub, dom, obj, act")
	require.NoError(t, err)

	rolesSchema := make(core.RolesSchema)
	rolesSchema["g"] = core.RoleSchema{
		Type:            core.RoleTypeWithDomain,
		DomainMatchFunc: core.RoleDomainMatchEqual,
	}
	rolesSchema["g1"] = core.RoleSchema{
		Type: core.RoleTypeWithoutDomain,
	}

	matchers, err := core.MatchersConfig{
		Roles: rolesSchema,
		Define: map[string]string{
			"m": "g(r.sub, p.sub, r.dom) && r.obj == p.obj && r.act == p.act",
		},
	}.New()
	require.NoError(t, err)

	m := core.NewModel(policy, request, rolesSchema, nil, matchers)

	itr := &mockAssertionIterator{
		data: [][]string{
			{"p", user, resource, "get"},
			{"p", user, resource, "set"},
			{"g", user, "admin", "console"},
			{"g1", user, "admin"},
		},
	}

	policies, _, err := m.Load(itr)
	require.NoError(t, err)
	assert.Len(t, policies, 2)
}
