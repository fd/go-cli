package cli

type (
	Arguments        []string
	ArgumentConsumer func(s string) (bool, error)
)

func (a *Arguments) Consume(f ArgumentConsumer) error {
	var (
		s = *a
		n Arguments
	)

	for _, arg := range s {
		consume, err := f(arg)
		if err != nil {
			return err
		}

		if !consume {
			n = append(n, arg)
		}
	}

	*a = n
	return nil
}

func (a *Arguments) Len() int {
	return len(*a)
}

func (a *Arguments) At(idx int) string {
	if idx >= a.Len() {
		return ""
	}
	return (*a)[idx]
}

func (a *Arguments) Skip(n int) {
	if n > a.Len() {
		n = a.Len()
	}

	s := *a
	s = s[n:]
	*a = s
}
