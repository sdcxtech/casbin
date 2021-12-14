package core

import (
	"fmt"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"

	"github.com/sdcxtech/casbin/core/graph"
)

type ManagerDomainMatchFunc func(request string, mapping string) (matched bool)

type RoleMapping struct {
	key         string
	graph       *graph.Graph
	domainMatch ManagerDomainMatchFunc
}

func (rg *RoleMapping) Key() string {
	return rg.key
}

func NewRoleMapping(key string) *RoleMapping {
	return &RoleMapping{
		key:         key,
		graph:       graph.New(),
		domainMatch: nil,
	}
}

func (rg *RoleMapping) GenerateGFunc(moreGraphs ...*graph.Graph) (overload *functions.Overload) {
	graphs := make([]*graph.Graph, 0, len(moreGraphs)+1)
	graphs = append(graphs, moreGraphs...)
	graphs = append(graphs, rg.graph)

	fn := generateG(graphs, rg.domainMatch)
	if rg.domainMatch == nil {
		f := func(lhs, rhs ref.Val) ref.Val {
			v1, ok := lhs.(types.String)
			if !ok {
				return types.NewErr("first arg: wrong args type: expect string got %v", lhs.Type().TypeName())
			}
			v2, ok := rhs.(types.String)
			if !ok {
				return types.NewErr("second arg: wrong args type: expect string got %v", rhs.Type().TypeName())
			}

			matched := fn(string(v1), string(v2), "")

			return types.Bool(matched)
		}
		overload = &functions.Overload{
			Operator: rg.key,
			Binary:   f,
		}
	} else {
		f := func(values ...ref.Val) ref.Val {
			v1, ok := values[0].(types.String)
			if !ok {
				return types.NewErr("first arg: wrong args type: expect string got %v", values[0].Type().TypeName())
			}
			v2, ok := values[1].(types.String)
			if !ok {
				return types.NewErr("second arg: wrong args type: expect string got %v", values[1].Type().TypeName())
			}
			v3, ok := values[2].(types.String)
			if !ok {
				return types.NewErr("third arg: wrong args type: expect string got %v", values[2].Type().TypeName())
			}

			matched := fn(string(v1), string(v2), string(v3))
			return types.Bool(matched)
		}
		overload = &functions.Overload{
			Operator: rg.key,
			Function: f,
		}
	}
	return
}

type RoleMappings map[string]*RoleMapping

func (rg RoleMappings) GenerateGFuncs(
	moreGraphs map[string][]*graph.Graph,
) (funcs []*functions.Overload) {
	if len(rg) == 0 {
		return
	}

	funcs = make([]*functions.Overload, 0, len(rg))
	for _, rg := range rg {
		key := rg.Key()
		graphs := moreGraphs[key]
		overload := rg.GenerateGFunc(graphs...)
		funcs = append(funcs, overload)
	}
	return
}

func (rg RoleMappings) SetDomainMatchFuncion(key string, domainMatch ManagerDomainMatchFunc) (err error) {
	g, ok := rg[key]
	if !ok {
		err = fmt.Errorf("not found role group with key: %s", key)
		return
	}
	g.domainMatch = domainMatch
	return
}

func RoleDomainMatchEqual(request, mapping string) bool {
	return request == mapping
}
