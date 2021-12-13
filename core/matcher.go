package core

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/interpreter/functions"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/sdcxtech/casbin/core/graph"
)

type ExtensionFunc struct {
	Decl     *exprpb.Decl
	Overload *functions.Overload
}

type Matcher struct {
	ast   *cel.Ast
	env   *cel.Env
	funcs []*functions.Overload
}

func (m *Matcher) Program(
	vars map[string]interface{},
	roleMappings RoleMappings,
	moreGraphs map[string][]*graph.Graph,
) (prg cel.Program, err error) {
	funcs := roleMappings.GenerateGFuncs(moreGraphs)
	funcs = append(funcs, m.funcs...)
	prg, err = m.env.Program(
		m.ast,
		cel.EvalOptions(cel.OptOptimize),
		cel.Functions(funcs...),
		cel.Globals(vars),
	)
	return
}

type Matchers struct {
	matchers map[string]Matcher
}

func (m Matchers) Get(key string) (matcher Matcher, err error) {
	matcher, ok := m.matchers[key]
	if !ok {
		err = fmt.Errorf("%w: %s", ErrNotFoundMatcher, key)
		return
	}
	return
}

type MatchersConfig struct {
	Roles          RolesSchema
	Define         map[string]string
	ExtensionFuncs []ExtensionFunc
}

func (c MatchersConfig) New() (m Matchers, err error) {
	if c.Roles == nil || c.Define == nil {
		err = fmt.Errorf("must give roles schema and matcher define")
		return
	}
	if len(c.Define) == 0 {
		err = ErrAtLeastOneMatcher
		return
	}

	m.matchers = make(map[string]Matcher)

	_decls := []*exprpb.Decl{
		decls.NewVar(
			"r",
			decls.NewMapType(decls.String, decls.String),
		),
		decls.NewVar(
			"p",
			decls.NewMapType(decls.String, decls.String),
		),
	}

	funcs := make([]*functions.Overload, 0, len(c.ExtensionFuncs))
	for _, fn := range c.ExtensionFuncs {
		_decls = append(_decls, fn.Decl)
		funcs = append(funcs, fn.Overload)
	}

	for key, rType := range c.Roles {
		if rType == RoleTypeWithDomain {
			_decls = append(
				_decls,
				decls.NewFunction(key, decls.NewParameterizedOverload(
					fmt.Sprintf("g_%s", key),
					[]*exprpb.Type{decls.String, decls.String, decls.String},
					decls.Bool,
					[]string{"from", "to", "domain"},
				)),
			)
		} else if rType == RoleTypeWithoutDomain {
			_decls = append(
				_decls,
				decls.NewFunction(key, decls.NewParameterizedOverload(
					fmt.Sprintf("g_%s", key),
					[]*exprpb.Type{decls.String, decls.String},
					decls.Bool,
					[]string{"from", "to"},
				)),
			)
		}
	}

	env, err := cel.NewEnv(
		cel.Declarations(_decls...),
	)
	if err != nil {
		err = fmt.Errorf("new cel env: %w", err)
		return
	}

	for key, v := range c.Define {
		ast, iss := env.Compile(v)
		if iss.Err() != nil {
			err = fmt.Errorf("compile match expression: %w", iss.Err())
			return
		}
		matcher := Matcher{
			env:   env,
			ast:   ast,
			funcs: funcs,
		}
		m.matchers[key] = matcher
	}

	return
}
