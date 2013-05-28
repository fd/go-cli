package cli

import (
	"strings"
)

type Environment struct {
	args Arguments
	env  map[string]string
}

func NewEnvironment(args []string, env []string) *Environment {
	return &Environment{args: Arguments(args), env: parse_env(env)}
}

func (e *Environment) Args() *Arguments {
	return &e.args
}

func (e *Environment) Var(key string) (string, bool) {
	val, p := e.env[key]
	return val, p
}

func parse_env(in []string) map[string]string {
	m := make(map[string]string, len(in))

	for _, l := range in {
		parts := strings.SplitN(l, "=", 2)
		if len(parts) == 1 {
			m[parts[0]] = ""
		} else {
			m[parts[0]] = parts[1]
		}
	}

	return m
}
