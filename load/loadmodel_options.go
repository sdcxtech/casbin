package load

import (
	"fmt"

	"github.com/sdcxtech/casbin/core"
)

type modelConfig struct {
	extensionFuncs []core.ExtensionFunc
}

// An option configures a new model when load it.
type ModelOption interface {
	apply(*modelConfig) error
}

type loadModelOptionFunc func(c *modelConfig) error

func (f loadModelOptionFunc) apply(c *modelConfig) error {
	return f(c)
}

func newLoadModelConfig(options ...ModelOption) (modelConfig, error) {
	var c modelConfig

	err := applyLoadModelConfigOptions(&c, options...)

	return c, err
}

func applyLoadModelConfigOptions(c *modelConfig, options ...ModelOption) error {
	for _, o := range options {
		if err := o.apply(c); err != nil {
			err = fmt.Errorf("apply model option: %w", err)

			return err
		}
	}

	return nil
}

func ExtensionFuncs(funcs ...core.ExtensionFunc) ModelOption {
	return loadModelOptionFunc(func(c *modelConfig) error {
		c.extensionFuncs = funcs

		return nil
	})
}
