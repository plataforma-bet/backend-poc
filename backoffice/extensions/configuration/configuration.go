package configuration

import (
	"fmt"

	"github.com/ardanlabs/conf/v3"
)

type ConfigOpt struct {
	prefix string
}

type ConfigOpts func(*ConfigOpt)

func Load[T any](opts ...ConfigOpts) func() (T, error) {
	return func() (T, error) {
		opt := &ConfigOpt{}
		for _, o := range opts {
			o(opt)
		}

		var config T

		if _, err := conf.Parse(opt.prefix, &config); err != nil {
			return config, fmt.Errorf("error parsing configuration: %w", err)
		}

		return config, nil
	}
}
