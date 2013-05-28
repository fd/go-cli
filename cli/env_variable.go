package cli

import (
	"fmt"
	"reflect"
	"strconv"
)

type EnvVariable struct {
	names []string
	value reflect.Value
}

func NewEnvVariable(names ...string) *EnvVariable {
	return &EnvVariable{names: names}
}

func (e *EnvVariable) Bind(v interface{}) *EnvVariable {
	return e.BindValue(reflect.ValueOf(v))
}

func (e *EnvVariable) BindValue(v reflect.Value) *EnvVariable {
	e.value = v
	return e
}

func (e *EnvVariable) Execute(env *Environment) error {
	var (
		name  string
		value string
		p     bool
	)

	for _, name = range e.names {
		value, p = env.Var(name)
		if p {
			break
		}
	}

	if value == "" {
		return nil
	}

	return e.handle_value(value)
}

func (e *EnvVariable) handle_value(val string) error {
	switch e.value.Kind() {

	case reflect.Ptr:
		if e.value.IsNil() {
			e.value.Set(reflect.New(e.value.Type().Elem()))
		}
		e.value = e.value.Elem()
		return e.handle_value(val)

	case reflect.Bool:
		x := val == "true" || val == "t" || val == "1"
		e.value.SetBool(x)
		return nil

	case reflect.String:
		e.value.SetString(val)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := strconv.ParseInt(val, 10, e.value.Type().Bits())
		if err != nil {
			return err
		}
		e.value.SetInt(x)
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := strconv.ParseUint(val, 10, e.value.Type().Bits())
		if err != nil {
			return err
		}
		e.value.SetUint(x)
		return nil

	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(val, e.value.Type().Bits())
		if err != nil {
			return err
		}
		e.value.SetFloat(x)
		return nil

	default:
		return fmt.Errorf("flag: `%s` unsupported value type %s", e.value.Type())

	}

	panic("not reached")
}
