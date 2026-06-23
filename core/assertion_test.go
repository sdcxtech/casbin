package core_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sdcxtech/casbin/core"
)

func TestAssertion(t *testing.T) {
	const (
		action   = "get"
		resource = "order"
	)

	a, err := core.NewAssertionSchema("sub, obj, act")
	require.NoError(t, err)

	_, err = a.CreateAssertion([]string{"charlie", resource, action})
	require.NoError(t, err)

	_, err = a.CreateAssertion([]string{"charlie", resource, action, "foobar"})
	require.Error(t, err)

	_, err = core.NewAssertionSchema("sub")
	require.Error(t, err)

	_, err = core.NewAssertionSchema("sub,,act")
	require.Error(t, err)
}
