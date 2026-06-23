package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sdcxtech/casbin/core"
)

func TestRoleType(t *testing.T) {
	c, err := core.RoleTypeFromLine("_, _,_")
	require.NoError(t, err)
	assert.Equal(t, core.RoleTypeWithDomain, c)

	c, err = core.RoleTypeFromLine("_, _")
	require.NoError(t, err)
	assert.Equal(t, core.RoleTypeWithoutDomain, c)

	_, err = core.RoleTypeFromLine("_,_,_,_")
	assert.Error(t, err)
}
