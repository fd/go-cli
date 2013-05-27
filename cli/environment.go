package cli

import (
	"os"
)

type Environment struct {
	args Arguments
}

func NewEnvironment(args []string) *Environment {
	return &Environment{args: Arguments(args)}
}

func (e *Environment) Args() *Arguments {
	return &e.args
}

func (e *Environment) Var(key string) string {
	return os.Getenv(key)
}
