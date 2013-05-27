package cli

import (
	"reflect"
	"strings"
)

type Group struct {
	names    []string
	flags    MultiHandler
	commands []NamedHandler
}

func NewGroup(names ...string) *Group {
	return &Group{names: names}
}

func (g *Group) Bind(v interface{}) *Group {
	return g.BindValue(reflect.ValueOf(v))
}

func (g *Group) BindValue(v reflect.Value) *Group {
	rv := v
	rt := rv.Type()

	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rt).Elem())
		}
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		panic("Expected a struct when binding a group")
	}

	for i, j := 0, rt.NumField(); i < j; i++ {
		f := rt.Field(i)
		fv := rv.Field(i)

		if f.PkgPath != "" {
			continue
		}

		if tag := f.Tag.Get("flag"); tag != "" {
			names := strings.Split(tag, ",")
			g.Flag(names...).BindValue(fv)
		}

		if tag := f.Tag.Get("cmd"); tag != "" {
			names := strings.Split(tag, ",")
			g.Command(names...).BindValue(fv)
		}

		if tag := f.Tag.Get("grp"); tag != "" {
			names := strings.Split(tag, ",")
			g.Group(names...).BindValue(fv)
		}
	}

	return g
}

func (g *Group) Names() []string {
	return g.names
}

func (g *Group) Flag(names ...string) *Flag {
	f := NewFlag(names...)
	g.flags = append(g.flags, f)
	return f
}

func (g *Group) Command(names ...string) *Command {
	c := NewCommand(names...)
	g.commands = append(g.commands, c)
	return c
}

func (g *Group) Group(names ...string) *Group {
	c := NewGroup(names...)
	g.commands = append(g.commands, c)
	return c
}

func (g *Group) Execute(env *Environment) error {
	args := env.Args()

	if args.Len() == 0 {
		return &MissingCommandError{}
	}

	name := args.At(0)
	if len(g.names) > 0 {
		found := false
		for _, n := range g.names {
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

	err := g.flags.Execute(env)
	if err != nil {
		return err
	}

	if args.Len() == 0 {
		return &MissingCommandError{}
	}

	name = args.At(0)
	for _, c := range g.commands {
		for _, n := range c.Names() {
			if name == n {
				return c.Execute(env)
			}
		}
	}

	return &UnrecognizedCommandError{name}
}
