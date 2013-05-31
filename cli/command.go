package cli

import (
	"reflect"
)

/*

Example command:
  type Cmd struct {
    Root|ParentCommand
    Arg0

    FlagsAndEnvs
    Args
  }

*/

type Command interface {
	Execute() error
}

type Arg0 string

type Root struct {
}

type Manual struct {
	usage      string
	summary    string
	paragraphs []paragraph_t
}

type paragraph_t struct {
	Header string
	Body   string
}

var (
	typ_Command = reflect.TypeOf((*Command)(nil)).Elem()
	typ_Arg0    = reflect.TypeOf((*Arg0)(nil)).Elem()
	typ_Root    = reflect.TypeOf((*Root)(nil)).Elem()
	typ_Manual  = reflect.TypeOf((*Manual)(nil)).Elem()
)
