package main

import (
	"fmt"
	"github.com/fd/go-cli/cli"
	"os"
)

func init() {
	cli.Register(Root{})
	cli.Register(App{})
	cli.Register(AppList{})
	cli.Register(AppCreate{})
	cli.Register(AppDestroy{})
}

type Root struct {
	cli.Root
	cli.Arg0

	Help    bool `flag:"-h,--help"`
	Verbose bool `flag:"-v,--verbose" env:"VERBOSE"`
	Debug   bool `flag:"--debug" env:"DEBUG"`

	cli.Manual `
    Usage:   example CMD ...ARGS
    Summary: Example cli tool.
  `
}

type App struct {
	Root
	cli.Arg0 `name:"application"`

	cli.Manual `
    Usage:   example app CMD ...ARGS
    Summary: Manage the applications.
  `
}

type AppList struct {
	App
	cli.Arg0 `name:"list"`

	cli.Manual `
    Usage:   example app list
    Summary: List the applications.
  `
}

type AppCreate struct {
	App
	cli.Arg0 `name:"create"`
	Region   string `flag:"--region" env:"REGION"`
	Name     string `arg`

	cli.Manual `
    Usage:   example app create NAME
    Summary: Create an application.
    .Region: The region to create the app in
    .Name:   The *name* of the application
  `
}

type AppDestroy struct {
	App
	cli.Arg0 `name:"destroy"`
	Name     string `arg`

	cli.Manual `
    Usage:   example app destroy NAME
    Summary: Destroy an application.
    .Name:   The *name* of the application
  `
}

func (cmd *AppList) Execute() error {
	if cmd.Help {
		fmt.Printf("CMD: %+v\n", cmd)
	}
	return nil
}

func (cmd *AppCreate) Execute() error {
	if cmd.Help {
		cmd.Manual.Open()
	}
	return nil
}

func (cmd *AppDestroy) Execute() error {
	return nil
}

func main() {
	err := cli.Main(os.Args, os.Environ())
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}
