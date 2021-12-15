package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sdcxtech/casbin/core"
)

func TestAssertion(t *testing.T) {
	a, err := core.NewAssertionSchema("sub, obj, act")
	assert.NoError(t, err)

	_, err = a.CreateAssertion([]string{"charlie", "order", "get"})
	assert.NoError(t, err)

	_, err = a.CreateAssertion([]string{"charlie", "order", "get", "foobar"})
	assert.Error(t, err)

	_, err = core.NewAssertionSchema("sub")
	assert.Error(t, err)

	_, err = core.NewAssertionSchema("sub,,act")
	assert.Error(t, err)
}
