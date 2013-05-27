package cli

import (
	"reflect"
)

type Handler interface {
	Execute(env *Environment) error
}

type NamedHandler interface {
	Handler
	Names() []string
}

var typ_Handler = reflect.TypeOf((*Handler)(nil)).Elem()
