package cli

import (
	"fmt"
	"reflect"
	"strings"
)

type Command struct {
	names    []string
	env_vars MultiHandler
	flags    MultiHandler
	args     MultiHandler
	handler  Handler
}

func NewCommand(names ...string) *Command {
	return &Command{names: names}
}

func (c *Command) Names() []string {
	return c.names
}

func (c *Command) Bind(v interface{}) *Command {
	return c.BindValue(reflect.ValueOf(v))
}

func (c *Command) BindValue(v reflect.Value) *Command {
	rv := v
	rt := rv.Type()

	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rt.Elem()))
		}

		if rt.Implements(typ_Handler) {
			c.Handle(rv.Interface().(Handler))
		}

		rv = reflect.Indirect(rv)
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		panic("Expected a struct when binding a command")
	}

	for i, j := 0, rt.NumField(); i < j; i++ {
		f := rt.Field(i)
		fv := rv.Field(i)

		if f.PkgPath != "" {
			continue
		}

		if tag := f.Tag.Get("env"); tag != "" {
			names := strings.Split(tag, ",")
			c.Var(names...).BindValue(fv)
		}

		if tag := f.Tag.Get("flag"); tag != "" {
			names := strings.Split(tag, ",")
			c.Flag(names...).BindValue(fv)
		}

		// if tag := f.Tag.Get("arg"); tag != "" {
		//   names := strings.Split(tag, ",")
		//   c.Arg().Bind(fv)
		// }
	}

	return c
}

func (c *Command) Handle(h Handler) *Command {
	c.handler = h
	return c
}

func (c *Command) HandleFunc(f FuncHandler) *Command {
	c.handler = f
	return c
}

func (c *Command) Var(names ...string) *EnvVariable {
	f := NewEnvVariable(names...)
	c.env_vars = append(c.env_vars, f)
	return f
}

func (c *Command) Flag(names ...string) *Flag {
	f := NewFlag(names...)
	c.flags = append(c.flags, f)
	return f
}

func (c *Command) Execute(env *Environment) error {
	args := env.Args()

	if args.Len() == 0 {
		return &MissingCommandError{}
	}

	name := args.At(0)
	if len(c.names) > 0 {
		found := false
		for _, n := range c.names {
			if name == n {
				found = true
				break
			}
		}
		if !found {
			return &UnrecognizedCommandError{name}
		}
	}

	args.Skip(1)

	err := c.env_vars.Execute(env)
	if err != nil {
		return err
	}

	err = c.flags.Execute(env)
	if err != nil {
		return err
	}

	err = c.args.Execute(env)
	if err != nil {
		return err
	}

	if env.Args().Len() != 0 {
		return fmt.Errorf("unexpected arguments: %s", env.Args())
	}

	if c.handler == nil {
		return nil
	}

	return c.handler.Execute(env)
}

type UnrecognizedCommandError struct {
	Command string
}

func (err *UnrecognizedCommandError) Error() string {
	return fmt.Sprintf("unrecognized command `%s`", err.Command)
}

type MissingCommandError struct {
}

func (err *MissingCommandError) Error() string {
	return "missing command"
}
