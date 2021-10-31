package cli

import "github.com/timotto/semver-bumper/pkg/bumper"

func (rt runtime) run() error {
	if err := rt.beforeResult(); err != nil {
		return err
	}

	version, commits, err := bumper.Bump(rt.opts, rt.repo, rt.esti)
	if err != nil {
		return err
	}

	if err := rt.onResult(version, commits); err != nil {
		return err
	}

	return nil
}
