package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sdcxtech/casbin/core"
)

func TestRoleType(t *testing.T) {
	c, err := core.RoleTypeFromLine("_, _,_")
	assert.NoError(t, err)
	assert.Equal(t, core.RoleTypeWithDomain, c)

	c, err = core.RoleTypeFromLine("_, _")
	assert.NoError(t, err)
	assert.Equal(t, core.RoleTypeWithoutDomain, c)

	_, err = core.RoleTypeFromLine("_,_,_,_")
	assert.Error(t, err)
}
