package load

import "github.com/sdcxtech/casbin/core"

type loadModelConfig struct {
	extensionFuncs []core.ExtensionFunc
}

type LoadModelOption interface {
	apply(*loadModelConfig) error
}

type LoadModelOptionFunc func(c *loadModelConfig) error

func (f LoadModelOptionFunc) apply(c *loadModelConfig) error {
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

func ExtensionFuncs(funcs []core.ExtensionFunc) LoadModelOption {
	return LoadModelOptionFunc(func(c *loadModelConfig) error {
		c.extensionFuncs = funcs
		return nil
	})
}
