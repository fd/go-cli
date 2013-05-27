package cli

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

var flag_exp = regexp.MustCompile("^([-][a-zA-Z0-9]|[-]{2}[a-zA-Z0-9][a-zA-Z0-9-]*)(?:[=](.+))?$")

type Flag struct {
	names         []string
	expects_value bool
	value         reflect.Value
}

func NewFlag(names ...string) *Flag {
	return &Flag{names: names}
}

func (f *Flag) Bind(v interface{}) *Flag {
	return f.BindValue(reflect.ValueOf(v))
}

func (f *Flag) BindValue(v reflect.Value) *Flag {
	f.value = v
	if f.terminal_type().Kind() != reflect.Bool {
		f.expects_value = true
	}
	return f
}

type flag_handler_context struct {
	name               string
	still_need_a_value bool
}

func (flag *Flag) Execute(env *Environment) error {
	var (
		ctx flag_handler_context
	)

	err := env.Args().Consume(func(arg string) (bool, error) {
		if ctx.still_need_a_value {
			ctx.still_need_a_value = false
			return flag.handle_value(&ctx, arg)
		}

		m := flag_exp.FindStringSubmatch(arg)
		if len(m) == 0 {
			return false, nil
		}

		name_matched := false
		ctx.name = m[1]
		value := m[2]

		for _, n := range flag.names {
			if n == ctx.name {
				name_matched = true
				break
			}
		}

		if !name_matched {
			return false, nil
		}

		if !flag.expects_value {
			if value != "" {
				return false, fmt.Errorf("flag `%s` doesn't accept a value", ctx.name)
			} else {
				return flag.handle_flag(&ctx)
			}
		}

		if value == "" {
			ctx.still_need_a_value = true
			return true, nil
		}

		return flag.handle_value(&ctx, value)
	})
	if err != nil {
		return err
	}

	if ctx.still_need_a_value {
		return fmt.Errorf("flag `%s` expected a value", ctx.name)
	}

	return nil
}

func (flag *Flag) handle_value(ctx *flag_handler_context, val string) (bool, error) {
	switch flag.value.Kind() {

	case reflect.Ptr:
		if flag.value.IsNil() {
			flag.value.Set(reflect.New(flag.value.Type().Elem()))
		}
		flag.value = flag.value.Elem()
		return flag.handle_value(ctx, val)

	case reflect.String:
		flag.value.SetString(val)
		return true, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := strconv.ParseInt(val, 10, flag.value.Type().Bits())
		if err != nil {
			return false, err
		}
		flag.value.SetInt(x)
		return true, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := strconv.ParseUint(val, 10, flag.value.Type().Bits())
		if err != nil {
			return false, err
		}
		flag.value.SetUint(x)
		return true, nil

	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(val, flag.value.Type().Bits())
		if err != nil {
			return false, err
		}
		flag.value.SetFloat(x)
		return true, nil

	default:
		return false, fmt.Errorf("flag: `%s` unsupported value type %s", flag.value.Type())

	}

	panic("not reached")
}

func (flag *Flag) handle_flag(ctx *flag_handler_context) (bool, error) {
	switch flag.value.Kind() {

	case reflect.Ptr:
		if flag.value.IsNil() {
			flag.value.Set(reflect.New(flag.value.Type().Elem()))
		}
		flag.value = flag.value.Elem()
		return flag.handle_flag(ctx)

	case reflect.Bool:
		flag.value.SetBool(true)
		return true, nil

	default:
		return false, fmt.Errorf("flag: `%s` unsupported value type %s", flag.value.Type())

	}

	panic("not reached")
}

func (f *Flag) terminal_type() reflect.Type {
	t := f.value.Type()

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}
