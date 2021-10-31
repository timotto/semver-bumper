package config

import (
	"fmt"
	"github.com/jessevdk/go-flags"
)

func FromOsArgs(osArgs []string) (*Options, string, error) {
	opts := Options{}
	args, err := flags.
		NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash).
		ParseArgs(osArgs)
	if err != nil {
		return nil, "", err
	}
	args = args[1:]

	var gitRepoPath string
	switch len(args) {
	case 0:
		gitRepoPath = "."
	case 1:
		gitRepoPath = args[0]
	default:
		return nil, "", fmt.Errorf("too many arguments: %v", args[1:])
	}

	return &opts, gitRepoPath, nil
}
