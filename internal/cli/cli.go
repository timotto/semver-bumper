package cli

func Run(os Os) error {
	if err := run(os); err != nil {
		Errln(os, err.Error())
		return err
	}

	return nil
}

func run(os Os) error {
	rt, err := newRuntime(os)
	if err != nil {
		return err
	}

	return rt.run()
}
