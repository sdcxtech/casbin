package graph_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sdcxtech/casbin/core/graph"
)

func TestSimple(t *testing.T) {
	g := graph.New()
	g.AddEdge("a", "b", "")
	g.AddEdge("b", "c", "")
	g.AddEdge("x", "y", "")

	fn := func(src, dst string) (has bool) {
		has, _ = g.HasLink([]string{src}, dst, nil)

		return
	}

	assert.True(t, fn("a", "b"))
	assert.True(t, fn("b", "c"))
	assert.True(t, fn("a", "c"))
	assert.False(t, fn("a", "y"))
	assert.False(t, fn("a", "x"))
	assert.True(t, fn("a", "a"))
}

func TestComposed(t *testing.T) {
	g1 := graph.New()
	g1.AddEdge("a", "b", "")
	g1.AddEdge("b", "c", "")

	g2 := graph.New()
	g2.AddEdge("x", "y", "")
	g2.AddEdge("c", "y", "")

	assert.True(t, graph.HasLink("a", "b", nil, g1, g2))
	assert.True(t, graph.HasLink("a", "c", nil, g1, g2))
	assert.True(t, graph.HasLink("x", "y", nil, g1, g2))
	assert.True(t, graph.HasLink("a", "y", nil, g1, g2))

	assert.False(t, graph.HasLink("a", "b", nil))
	assert.True(t, graph.HasLink("a", "a", nil))
}

func TestComposedMore(t *testing.T) {
	g1 := graph.New()
	g1.AddEdge("a", "b", "")

	g2 := graph.New()
	g2.AddEdge("b", "c", "")

	g3 := graph.New()
	g3.AddEdge("c", "d", "")

	assert.True(t, graph.HasLink("a", "d", nil, g1, g2, g3))
	assert.False(t, graph.HasLink("c", "a", nil, g1, g2, g3))
}

func TestWithDomainMatch(t *testing.T) {
	g := graph.New()
	g.AddEdge("a", "b", "x")
	g.AddEdge("b", "c", "x")
	g.AddEdge("b", "d", "y")

	gEqMatch := func(reqDom string) graph.DomainMatchFunc {
		return func(domain string) bool {
			return reqDom == domain
		}
	}

	fn := func(src, dst string) (has bool) {
		has, _ = g.HasLink([]string{src}, dst, gEqMatch("x"))

		return
	}

	assert.True(t, fn("a", "b"))
	assert.True(t, fn("a", "c"))
	assert.True(t, fn("b", "c"))
	assert.False(t, fn("c", "a"))
	assert.False(t, fn("a", "d"))
}
