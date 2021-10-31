package main

import (
	. "github.com/timotto/semver-bumper/internal/cli"
	"os"
)

func main() {
	must(Run(&OS{}))
}

func must(err error) {
	if err != nil {
		os.Exit(1)
	}
}
