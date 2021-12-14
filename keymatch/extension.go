package keymatch

import (
	"fmt"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/sdcxtech/casbin/core"
)

// Key match function signature.
type Func func(key1, key2 string) (matched bool, err error)

func ToExtensionFunc(funcName string, fn Func) core.ExtensionFunc {
	return core.ExtensionFunc{
		Decl: decls.NewFunction(funcName, decls.NewParameterizedOverload(
			fmt.Sprintf("key_match_%s", funcName),
			[]*exprpb.Type{decls.String, decls.String},
			decls.Bool,
			[]string{"key1", "key2"},
		)),
		Overload: &functions.Overload{
			Operator: funcName,
			Binary:   celKeyMatchFunc(fn),
		},
	}
}

func celKeyMatchFunc(fn Func) functions.BinaryOp {
	return func(lhs ref.Val, rhs ref.Val) ref.Val {
		if lhs.Type().TypeName() != types.StringType.TypeName() {
			return types.NewErr("invalid arguments: key must be string type: key 1")
		}

		if rhs.Type().TypeName() != types.StringType.TypeName() {
			return types.NewErr("invalid arguments: key must be string type: key 2")
		}

		matched, err := fn(lhs.Value().(string), rhs.Value().(string))
		if err != nil {
			return types.NewErr("key match func: %s", err.Error())
		}

		return types.Bool(matched)
	}
}
