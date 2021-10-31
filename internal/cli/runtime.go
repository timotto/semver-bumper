package cli

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/timotto/semver-bumper/pkg/bumper"
	. "github.com/timotto/semver-bumper/pkg/config"
	"github.com/timotto/semver-bumper/pkg/estimator"
	"github.com/timotto/semver-bumper/pkg/gitrepo"
	"io"
)

type runtime struct {
	os   Os
	opts *Options
	repo bumper.GitRepo
	esti bumper.Estimator
}

//counterfeiter:generate . Os
type Os interface {
	Args() []string
	Stdout() io.Writer
	Stderr() io.Writer
}

func newRuntime(os Os) (*runtime, error) {
	opts, gitRepoPath, err := readOptions(os)
	if err != nil {
		return nil, err
	}

	repo, err := gitrepo.NewGitRepo(opts, gitRepoPath)
	if err != nil {
		return nil, err
	}

	esti := estimator.NewEstimator(opts)

	rt := &runtime{
		os:   os,
		opts: opts,
		repo: repo,
		esti: esti,
	}
	return rt, nil
}

func readOptions(os Os) (*Options, string, error) {
	opts, gitRepoPath, err := readCliArgs(os)
	if err != nil {
		return nil, "", err
	}

	return opts, gitRepoPath, opts.Valid()
}

func readCliArgs(os Os) (*Options, string, error) {
	opts, gitRepoPath, err := FromOsArgs(os.Args())
	if err != nil {
		return nil, "", err
	}

	return opts, gitRepoPath, maybeReadConfigFile(gitRepoPath, opts)
}

func maybeReadConfigFile(gitRepoPath string, opts *Options) error {
	filename, err := SearchConfigFile(gitRepoPath, opts.ConfigFile)
	if err != nil {
		return err
	}

	if filename == "" {
		return nil
	}

	cfg, err := FromFile(filename)
	if err != nil {
		return err
	}

	opts.SetMissingFrom(cfg)

	return nil
}
