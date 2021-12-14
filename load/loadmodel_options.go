package load

import "github.com/sdcxtech/casbin/core"

type loadModelConfig struct {
	extensionFuncs []core.ExtensionFunc
}

// An option configures a new model when load it.
type LoadModelOption interface {
	apply(*loadModelConfig) error
}

type loadModelOptionFunc func(c *loadModelConfig) error

func (f loadModelOptionFunc) apply(c *loadModelConfig) error {
	return f(c)
}

func newLoadModelConfig(options ...LoadModelOption) (loadModelConfig, error) {
	var c loadModelConfig

	err := applyLoadModelConfigOptions(&c, options...)
	return c, err
}

func applyLoadModelConfigOptions(c *loadModelConfig, options ...LoadModelOption) error {
	for _, o := range options {
		if err := o.apply(c); err != nil {
			return err
		}
	}
	return nil
}

func ExtensionFuncs(funcs ...core.ExtensionFunc) LoadModelOption {
	return loadModelOptionFunc(func(c *loadModelConfig) error {
		c.extensionFuncs = funcs
		return nil
	})
}
