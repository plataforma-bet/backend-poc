package env

import (
	"errors"
	"os"
)

var (
	dev  = "dev"
	prod = "prd"
	sdx  = "sdx"
)

type Environment struct {
	value string
}

func (e *Environment) String() string {
	return e.value
}

func (e *Environment) UnmarshalText(text []byte) error {
	env, err := ParseValue(string(text))
	if err != nil {
		return err
	}

	e.value = env.value

	return nil
}

func ParseValue(value string) (Environment, error) {
	switch value {
	case dev:
		return Environment{
			value: dev,
		}, nil
	case prod:
		return Environment{
			value: prod,
		}, nil
	case sdx:
		return Environment{
			value: sdx,
		}, nil
	default:
		return Environment{}, errors.New("invalid environment")
	}
}

func FromEnv() (Environment, error) {
	env, present := os.LookupEnv("ENV")
	if !present {
		env = dev
	}

	return ParseValue(env)
}

func (e *Environment) IsDev() bool {
	return e.value == dev
}
