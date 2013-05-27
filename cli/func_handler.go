package cli

type FuncHandler func(env *Environment) error

func (f FuncHandler) Execute(env *Environment) error {
	return f(env)
}
