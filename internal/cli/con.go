package cli

import (
	"fmt"
	"io"
)

func Outln(os Os, args ...interface{}) {
	pln(os.Stdout(), args...)
}

func Errln(os Os, args ...interface{}) {
	pln(os.Stderr(), args...)
}

func pln(w io.Writer, args ...interface{}) {
	must(fmt.Fprintln(w, args...))
}

func must(_ interface{}, err error) {
	if err != nil {
		panic(err)
	}
}
