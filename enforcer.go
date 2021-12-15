package casbin

import (
	"fmt"

	"github.com/sdcxtech/casbin/core"
	"github.com/sdcxtech/casbin/core/graph"
)

// Enforcer check if a user has requested permission.
//
// Default to use matcher with name 'm', can be changed by option.
type Enforcer struct {
	model        *core.Model
	policies     core.Policies
	roleMappings core.RoleMappings
}

type enforceConfig struct {
	withRoleGraphs map[string][]*graph.Graph
	matcher        string
}

// Enforce check if the requested permissions assertion is allowed.
func (e *Enforcer) Enforce(requestValues []string, options ...EnforceOption) (allow bool, err error) {
	conf, err := newEnforceConfig(options...)
	if err != nil {
		return
	}

	matcher, err := e.model.Matchers().Get(conf.matcher)
	if err != nil {
		return
	}

	rVar, err := e.model.Request().CreateAssertion(requestValues)
	if err != nil {
		err = fmt.Errorf("convert request vals: %w", err)

		return
	}

	prg, err := matcher.Program(
		map[string]interface{}{
			"r": rVar,
		},
		e.roleMappings,
		conf.withRoleGraphs,
	)
	if err != nil {
		err = fmt.Errorf("new program: %w", err)

		return
	}

	vars := make(map[string]interface{}, 1)

	policyEval := func(policy core.Assertion) (allow bool, err error) {
		vars["p"] = policy

		result, _ /*details*/, _err := prg.Eval(vars)

		if _err != nil {
			err = _err

			return
		}

		allow, ok := result.Value().(bool)
		if !ok {
			err = fmt.Errorf(
				"%w: eval result should be bool type, but got %s",
				core.ErrUnexpectedEvalResult,
				result.Type().TypeName(),
			)

			return
		}

		return
	}

	allow, err = e.model.Effector().Execute(policyEval, e.policies)
	if err != nil {
		err = fmt.Errorf("effector execute: %w", err)

		return
	}

	return allow, err
}

// NewEnforcer new an enforcer by a model and a policy iterator.
func NewEnforcer(
	model *core.Model,
	policyItr core.LoadIterator,
) (enforcer *Enforcer, err error) {
	policies, roles, err := model.Load(policyItr)
	if err != nil {
		err = fmt.Errorf("new enforcer: %w", err)

		return
	}

	enforcer = &Enforcer{
		model:        model,
		policies:     policies,
		roleMappings: roles,
	}

	return
}
