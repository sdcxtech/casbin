package core

import (
	"fmt"

	"github.com/sdcxtech/casbin/core/graph"
)

func generateG(graphs []*graph.Graph, domainMatch func(string, string) bool) func(r, p, domain string) bool {
	memorized := map[string]bool{}

	return func(r, p, domain string) bool {
		key := fmt.Sprintf("%s:%s:%s", r, p, domain)

		v, ok := memorized[key]
		if ok {
			return v
		}

		if len(graphs) == 0 {
			v = r == p
			memorized[key] = v

			return v
		}

		if domainMatch == nil {
			v = graph.HasLink(r, p, nil, graphs...)
		} else {
			v = graph.HasLink(
				r,
				p,
				func(d string) bool {
					return domainMatch(domain, d)
				},
				graphs...,
			)
		}

		memorized[key] = v

		return v
	}
}
