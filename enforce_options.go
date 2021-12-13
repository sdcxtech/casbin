package casbin

import (
	"fmt"

	"github.com/sdcxtech/casbin/core/graph"
)

type ApplyEnforceOptionFunc func(c *enforceConfig) error

func (f ApplyEnforceOptionFunc) apply(c *enforceConfig) error {
	return f(c)
}

func newEnforceConfig(options ...EnforceOption) (enforceConfig, error) {
	var c enforceConfig
	c.matcher = "m"
	c.withRoleGraphs = make(map[string][]*graph.Graph)

	err := applyEnforceConfigOptions(&c, options...)
	return c, err
}

func applyEnforceConfigOptions(c *enforceConfig, options ...EnforceOption) error {
	for _, o := range options {
		if err := o.apply(c); err != nil {
			return err
		}
	}
	return nil
}

type EnforceOption interface {
	apply(*enforceConfig) error
}

type enforceOptionWithRoleGraphsImpl struct {
	gKey   string
	graphs []*graph.Graph
}

func (o enforceOptionWithRoleGraphsImpl) apply(c *enforceConfig) error {
	if c.withRoleGraphs == nil {
		c.withRoleGraphs = make(map[string][]*graph.Graph)
	}
	c.withRoleGraphs[o.gKey] = o.graphs
	return nil
}

// WithRoleGraphs inject per-calling role mapping graphs.
//
// Find in these grpahs first, then the global role mapping graph.
func WithRoleGraphs(gKey string, graphs ...*graph.Graph) EnforceOption {
	return enforceOptionWithRoleGraphsImpl{
		gKey:   gKey,
		graphs: graphs,
	}
}

func UseMatcher(name string) EnforceOption {
	return ApplyEnforceOptionFunc(func(c *enforceConfig) (err error) {
		if name == "" {
			err = fmt.Errorf("matcher name can not be empty")
			return
		}
		c.matcher = name
		return
	})
}
