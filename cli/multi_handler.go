package cli

type MultiHandler []Handler

func (m MultiHandler) Execute(env *Environment) error {
	for _, h := range m {
		err := h.Execute(env)
		if err != nil {
			return err
		}
	}
	return nil
}
