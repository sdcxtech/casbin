package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssertion(t *testing.T) {
	a, err := NewAssertionSchema("sub, obj, act")
	assert.NoError(t, err)

	_, err = a.CreateAssertion([]string{"charlie", "order", "get"})
	assert.NoError(t, err)
	_, err = a.CreateAssertion([]string{"charlie", "order", "get", "foobar"})
	assert.Error(t, err)

	_, err = NewAssertionSchema("sub,obj")
	assert.Error(t, err)
	_, err = NewAssertionSchema("sub,,act")
	assert.Error(t, err)
}
