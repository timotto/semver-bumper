package gitrepo

import (
	"fmt"
	. "github.com/timotto/semver-bumper/pkg/config"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
)

type Gitrepo struct {
	conf *Options
	repo *git.Repository
}

func NewGitRepo(conf *Options, path string) (*Gitrepo, error) {
	var err error

	r := &Gitrepo{conf: conf}
	r.repo, err = git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open git: %w", err)
	}

	return r, nil
}

func (g Gitrepo) LatestTaggedRelease() (*semver.Version, error) {
	versions, err := g.versionTags(true)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, nil
	}

	return versions.Latest().Tag, nil
}

func (g Gitrepo) LatestTaggedPrerelease() (*semver.Version, error) {
	versions, err := g.versionTags(false)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return g.LatestTaggedRelease()
	}

	return versions.Latest().Tag, nil
}

func (g Gitrepo) versionTags(strict bool) (collection, error) {
	iter, err := g.repo.Tags()
	if err != nil {
		return nil, fmt.Errorf("cannot list git tags: %w", err)
	}

	c := g.newCollector(strict)
	if err := iter.ForEach(c.collect); err != nil {
		return nil, err
	}

	sort.Sort(collection(c.Result))

	return c.Result, nil
}
