package core

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sdcxtech/casbin/core/graph"
)

func TestGenerateG(t *testing.T) {
	fn := generateG(nil, nil)
	assert.True(t, fn("a", "a", ""))
	assert.True(t, fn("a", "a", ""))

	g := graph.New()
	g.AddEdge("a", "b", "x")

	fn = generateG([]*graph.Graph{g}, nil)
	assert.True(t, fn("a", "b", ""))
	assert.False(t, fn("a", "c", ""))

	fn = generateG([]*graph.Graph{g}, DefaultRoleDomainMatch)
	assert.True(t, fn("a", "b", "x"))
	assert.False(t, fn("a", "b", ""))
}
