package bumper

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing/object"
	. "github.com/timotto/semver-bumper/pkg/model"
)

type (
	Config interface {
		BumpPrerelease() bool
		ShouldFakePrerelease() (string, bool)
		InitialVersionValue() *semver.Version
	}

	GitRepo interface {
		LatestTaggedRelease() (*semver.Version, error)
		LatestTaggedPrerelease() (*semver.Version, error)
		CommitMessagesSince(v *semver.Version) ([]*object.Commit, error)
	}

	Estimator interface {
		BumpLevelFrom(commitMessages []string) BumpLevel
		NextPrerelease(pre string) (string, error)
	}
)

func Bump(conf Config, repo GitRepo, esti Estimator) (*semver.Version, []*object.Commit, error) {
	nextRelease, commits, err := bumpRelease(repo, esti)
	if err != nil {
		return nil, nil, err
	}

	if !conf.BumpPrerelease() {
		if nextRelease == nil {
			return conf.InitialVersionValue(), commits, nil
		}

		return nextRelease, commits, nil
	}

	latestPrerelease, ok, err := fakePrerelease(conf)
	if err != nil {
		return nil, nil, err
	}

	if !ok {
		latestPrerelease, err = repo.LatestTaggedPrerelease()
		if err != nil {
			return nil, nil, err
		}
	}

	if latestPrerelease == nil {
		if nextRelease == nil {
			nextRelease = conf.InitialVersionValue()
		}

		return prerelease1(esti, commits, nextRelease)
	}

	if nextReleaseIsGreaterThanLastPrerelease(nextRelease, latestPrerelease) {
		return prerelease1(esti, commits, nextRelease)
	}

	return bumpPrerelease(esti, latestPrerelease, commits)
}

func bumpRelease(repo GitRepo, esti Estimator) (*semver.Version, []*object.Commit, error) {
	latestRelease, err := repo.LatestTaggedRelease()
	if err != nil {
		return nil, nil, err
	}

	commits, err := repo.CommitMessagesSince(latestRelease)
	if err != nil {
		return nil, nil, err
	}

	if latestRelease == nil {
		return nil, commits, err
	}

	nextRelease := bump(latestRelease, esti.BumpLevelFrom(messagesFrom(commits)))

	return &nextRelease, commits, nil
}

func fakePrerelease(conf Config) (*semver.Version, bool, error) {
	prerelease, ok := conf.ShouldFakePrerelease()
	if !ok {
		return nil, false, nil
	}

	version, err := semver.StrictNewVersion(prerelease)
	if err != nil {
		return nil, true, fmt.Errorf("cannot parse given prerelease version: %w", err)
	}

	return version, true, nil
}

func prerelease1(esti Estimator, commits []*object.Commit, v *semver.Version) (*semver.Version, []*object.Commit, error) {
	if nextPrerelease, err := esti.NextPrerelease(""); err != nil {
		return nil, nil, err
	} else if ver, err := v.SetPrerelease(nextPrerelease); err != nil {
		return nil, nil, err
	} else {
		return &ver, commits, nil
	}
}

func bumpPrerelease(esti Estimator, latestPrerelease *semver.Version, commits []*object.Commit) (*semver.Version, []*object.Commit, error) {
	pre, err := esti.NextPrerelease(latestPrerelease.Prerelease())
	if err != nil {
		return nil, nil, err
	}

	nextPrerelease, err := latestPrerelease.SetPrerelease(pre)
	if err != nil {
		return nil, nil, err
	}

	return &nextPrerelease, commits, nil
}

func bump(v *semver.Version, lvl BumpLevel) semver.Version {
	switch lvl {
	case BumpLevelMajor:
		return v.IncMajor()

	case BumpLevelMinor:
		return v.IncMinor()

	case BumpLevelPatch:
		return v.IncPatch()

	default:
		return *v
	}
}

func messagesFrom(commits []*object.Commit) []string {
	var result []string

	for _, commit := range commits {
		result = append(result, commit.Message)
	}

	return result
}

func nextReleaseIsGreaterThanLastPrerelease(nextRelease, lastPrerelease *semver.Version) bool {
	if nextRelease == nil {
		return false
	}

	lastPrereleaseReleased, _ := lastPrerelease.SetPrerelease("")

	return nextRelease.GreaterThan(&lastPrereleaseReleased)
}
