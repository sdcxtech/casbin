package core //nolint: testpackage

import (
	"testing"

	"github.com/google/cel-go/common/types"
	"github.com/stretchr/testify/assert"
)

func TestRoleMapping(t *testing.T) {
	g := NewRoleMapping("g", nil)
	overload := g.GenerateGFunc()
	result := overload.Binary(types.String("a"), types.String("b"))
	link, ok := result.Value().(bool)
	assert.True(t, ok)
	assert.False(t, link)

	g.domainMatch = RoleDomainMatchEqual
	overload = g.GenerateGFunc()
	result = overload.Function(types.String("a"), types.String("b"), types.String("x"))
	link, ok = result.Value().(bool)
	assert.True(t, ok)
	assert.False(t, link)
}

func TestRoleMappings(t *testing.T) {
	rg := RoleMappings{}
	{
		funcs := rg.GenerateGFuncs(nil)
		assert.Equal(t, 0, len(funcs))
	}

	rg["g"] = NewRoleMapping("g", RoleDomainMatchEqual)

	{
		funcs := rg.GenerateGFuncs(nil)
		assert.Equal(t, 1, len(funcs))
	}
}
