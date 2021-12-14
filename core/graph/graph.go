package graph

import (
	"github.com/sdcxtech/casbin/internal/sets"
)

type DomainMatchFunc func(domain string) bool

type graphEdge struct {
	src    string
	dst    string
	domain string
}

type Graph struct {
	edges map[string][]*graphEdge
}

func (g *Graph) HasLink(
	src []string,
	dst string,
	domainMatch DomainMatchFunc,
) (hasLink bool, reached []string) {
	foundNodes := sets.NewString(src...)
	if foundNodes.Has(dst) {
		hasLink = true

		return
	}

	srcLen := len(src)
	jobs := make([]string, srcLen, srcLen+16)
	copy(jobs, src)

	for i := 0; i < len(jobs); i++ {
		src := jobs[i]
		edges, ok := g.edges[src]

		if ok {
			for _, edge := range edges {
				if domainMatch != nil && !domainMatch(edge.domain) {
					continue
				}

				next := edge.dst
				if !foundNodes.Has(next) {
					if next == dst {
						hasLink = true

						return
					}

					foundNodes.Insert(next)
					jobs = append(jobs, next)
				}
			}
		}
	}

	reached = foundNodes.UnsortedList()

	return hasLink, reached
}

func (g *Graph) AddEdge(src, dst, domain string) {
	edge := &graphEdge{
		src:    src,
		dst:    dst,
		domain: domain,
	}
	edges := g.edges[edge.src]
	edges = append(edges, edge)
	g.edges[edge.src] = edges
}

func New() (g *Graph) {
	g = &Graph{
		edges: make(map[string][]*graphEdge),
	}

	return
}

func HasLink(
	src, dst string,
	domainMatch DomainMatchFunc,
	graphs ...*Graph,
) (hasLink bool) {
	if src == dst {
		return true
	}

	if len(graphs) == 0 {
		return
	}

	from := []string{src}

	for _, g := range graphs {
		hasLink, from = g.HasLink(from, dst, domainMatch)
		if hasLink {
			return
		}
	}

	return
}
