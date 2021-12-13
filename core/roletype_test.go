package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleType(t *testing.T) {
	c, err := RoleTypeFromLine("_, _,_")
	assert.NoError(t, err)
	assert.Equal(t, RoleTypeWithDomain, c)

	c, err = RoleTypeFromLine("_, _")
	assert.NoError(t, err)
	assert.Equal(t, RoleTypeWithoutDomain, c)

	_, err = RoleTypeFromLine("_,_,_,_")
	assert.Error(t, err)
}
