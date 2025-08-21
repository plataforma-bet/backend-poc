package config

import (
	"fmt"

	"github.com/ardanlabs/conf/v3"
)

type Opts struct {
	prefix string
}

type Opt func(*Opts)

func Prefix(prefix string) Opt {
	return func(opts *Opts) {
		opts.prefix = prefix
	}
}

type LazyConfig[T any] func() (T, error)

func Config[T any](opts ...Opt) LazyConfig[T] {
	return func() (T, error) {
		opt := &Opts{}
		for _, o := range opts {
			o(opt)
		}

		var config T

		if _, err := conf.Parse(opt.prefix, &config); err != nil {
			return config, fmt.Errorf("error parsing config: %w", err)
		}

		return config, nil
	}
}
