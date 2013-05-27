package cli

import (
	"testing"
)

func TestFlag(t *testing.T) {
	var (
		value string
		flag  = NewFlag("--bar", "-b").Bind(&value)
		env   *Environment
		err   error
	)

	value = ""
	env = NewEnvironment([]string{"hello", "--foo", "--bar", "qux", "--baz"})
	err = flag.Execute(env)
	if value != "qux" {
		t.Fatalf("expected value to be `qux` but it was `%s`", value)
	}
	if err != nil {
		t.Fatalf("expected error to be nil: %s", err)
	}
	t.Logf("args: %+v", env.Args())

	value = ""
	env = NewEnvironment([]string{"hello", "--foo", "-b", "qux", "--baz"})
	err = flag.Execute(env)
	if value != "qux" {
		t.Fatalf("expected value to be `qux` but it was `%s`", value)
	}
	if err != nil {
		t.Fatalf("expected error to be nil: %s", err)
	}
	t.Logf("args: %+v", env.Args())

	value = ""
	env = NewEnvironment([]string{"hello", "--foo", "--bar=qux", "--baz"})
	err = flag.Execute(env)
	if value != "qux" {
		t.Fatalf("expected value to be `qux` but it was `%s`", value)
	}
	if err != nil {
		t.Fatalf("expected error to be nil: %s", err)
	}
	t.Logf("args: %+v", env.Args())

	value = ""
	env = NewEnvironment([]string{"hello", "--foo", "--bar"})
	err = flag.Execute(env)
	if value != "" {
		t.Fatalf("expected value to be `` but it was `%s`", value)
	}
	if err == nil {
		t.Fatal("expected an error but it was nil")
	}
	t.Logf("args: %+v", env.Args())
}
