package cli

import (
	"io"
	"os"
)

type OS struct {
}

func (o OS) Args() []string {
	return os.Args
}

func (o OS) Stdout() io.Writer {
	return os.Stdout
}

func (o OS) Stderr() io.Writer {
	return os.Stderr
}
