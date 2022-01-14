package casbin_test

import (
	"strings"
	"testing"

	"github.com/sdcxtech/casbin"
	"github.com/sdcxtech/casbin/core"
	"github.com/sdcxtech/casbin/effector"
	"github.com/sdcxtech/casbin/load"
	"github.com/stretchr/testify/assert"
)

func TestEnforcer(t *testing.T) {
	policy, err := core.NewAssertionSchema("sub, obj, act")
	assert.NoError(t, err)

	request, err := core.NewAssertionSchema("sub, obj, act")
	assert.NoError(t, err)

	rolesSchema := make(core.RolesSchema)
	rolesSchema["g"] = core.RoleSchema{
		Type: core.RoleTypeWithoutDomain,
	}

	matchers, err := core.MatchersConfig{
		Roles: rolesSchema,
		Define: map[string]string{
			"m": "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act",
		},
	}.New()
	assert.NoError(t, err)

	m := core.NewModel(policy, request, rolesSchema, effector.NewAllowOverride(), matchers)

	csv := `
    p, staff, order, get
    g, admin, staff
    `
	reader := strings.NewReader(csv)

	itr := load.NewCSVIterator(reader)

	enforcer, err := casbin.NewEnforcer(m, itr)
	assert.NoError(t, err)

	allow, err := enforcer.Enforce(casbin.Request("admin", "order", "get"))
	assert.NoError(t, err)
	assert.True(t, allow)
}
