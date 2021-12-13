package load

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVIterator(t *testing.T) {
	csvContent := `

    p1, admin, order, *
    p2, admin, router, *
\t
    g, superadmin, admin
    `
	reader := strings.NewReader(csvContent)

	itr := NewCSVIterator(reader)

	var ok bool
	var key string
	var vals []string
	ok, key, vals = itr.Next()
	assert.True(t, ok)
	assert.Equal(t, "p1", key)
	assert.Len(t, vals, 3)

	ok, key, vals = itr.Next()
	assert.True(t, ok)
	assert.Equal(t, "p2", key)
	assert.Len(t, vals, 3)

	ok, key, vals = itr.Next()
	assert.True(t, ok)
	assert.Equal(t, "g", key)
	assert.Len(t, vals, 2)

	ok, key, vals = itr.Next()
	assert.False(t, ok)
	assert.Empty(t, key)
	assert.Empty(t, vals)
	assert.NoError(t, itr.Error())
}
