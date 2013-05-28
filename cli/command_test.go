package cli

import (
	"testing"
)

func TestCommand(t *testing.T) {
	var (
		cmd = NewCommand("echo", "say")

		verbose bool
		_       = cmd.Flag("-v", "--verbose").Bind(&verbose)

		person string
		_      = cmd.Flag("-h", "--hello").Bind(&person)

		env *Environment
		err error
	)

	verbose, person = false, ""
	env = NewEnvironment([]string{"say", "-v", "--hello", "Anaïs"}, nil)
	err = cmd.Execute(env)
	if !verbose {
		t.Fatalf("expected verbose to be `true` but it was `%v`", verbose)
	}
	if person != "Anaïs" {
		t.Fatalf("expected person to be `Anaïs` but it was `%s`", person)
	}
	if err != nil {
		t.Fatalf("expected error to be nil: %s", err)
	}
	t.Logf("args: %+v", env.Args())
}

func TestCommandBind(t *testing.T) {
	var (
		c struct {
			Verbose  bool   `flag:"-v,--verbose"`
			Greeting string `env:"GREETING"`
			Person   string `flag:"-h,--hello"`

			Manual `
        Usage:   echo --hello John
        Summary: A test command

        .Verbose:  A Short description of the --verbose option
        .Greeting:
          A Multi line description of
          the GREETING variable.
          This text is automatically rewrapped.

        This is a test command.
      `
		}
		cmd = NewCommand("echo", "say").Bind(&c)

		env *Environment
		err error
	)

	env = NewEnvironment(
		[]string{"say", "-v", "--hello", "Anaïs"},
		[]string{"GREETING=Hi"},
	)
	err = cmd.Execute(env)
	if err != nil {
		t.Fatalf("expected error to be nil: %s", err)
	}
	if !c.Verbose {
		t.Fatalf("expected verbose to be `true` but it was `%v`", c.Verbose)
	}
	if c.Greeting != "Hi" {
		t.Fatalf("expected greeting to be `Hi` but it was `%s`", c.Greeting)
	}
	if c.Person != "Anaïs" {
		t.Fatalf("expected person to be `Anaïs` but it was `%s`", c.Person)
	}
	t.Logf("args: %+v", env.Args())
}
