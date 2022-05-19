package casbin

import (
	"fmt"

	"github.com/sdcxtech/casbin/core"
	"github.com/sdcxtech/casbin/core/graph"
)

type applyEnforceOptionFunc func(c *enforceConfig) error

func (f applyEnforceOptionFunc) apply(c *enforceConfig) error {
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
			err = fmt.Errorf("apply enforce option: %w", err)

			return err
		}
	}

	return nil
}

// EnforceOption is an option for calling enforce on an enforcer.
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

// UseMatcher use the specified matcher instead of the default matcher named 'm'.
func UseMatcher(name string) EnforceOption {
	return applyEnforceOptionFunc(func(c *enforceConfig) (err error) {
		if name == "" {
			err = ErrEmptyMatcherName

			return
		}
		c.matcher = name

		return
	})
}

// WithPolicies inject per-calling raw policies.
func WithRawPolicies(policies [][]string) EnforceOption {
	return applyEnforceOptionFunc(func(c *enforceConfig) (err error) {
		c.rawPolicies = policies

		return
	})
}

// WithPolicies inject per-calling policies.
func WithPolicies(policies []core.Assertion) EnforceOption {
	return applyEnforceOptionFunc(func(c *enforceConfig) (err error) {
		c.policies = policies

		return
	})
}
