package gitrepo

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"io"
)

func (g Gitrepo) CommitMessagesSince(v *semver.Version) ([]*object.Commit, error) {
	if v == nil {
		return g.commitMessagesSince(nil)
	}

	tags, err := g.versionTags(true)
	if err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if tag.IsVersion(v) {
			return g.commitMessagesSince(&tag)
		}
	}

	return nil, nil
}

func (g Gitrepo) commitMessagesSince(tag *taggedCommit) ([]*object.Commit, error) {
	opts := &git.LogOptions{
		Order:      git.LogOrderCommitterTime,
		PathFilter: g.FiltersAccept,
	}

	if tag != nil {
		opts.Since = &tag.Ref.Committer.When
	}

	iter, err := g.repo.Log(opts)
	if err != nil {
		if tag == nil && err == plumbing.ErrReferenceNotFound {
			// bare / no commits
			return nil, nil
		}

		return nil, fmt.Errorf("cannot get log: %w", err)
	}

	return collectMessages(iter, tag)
}

type commitMessageCollector struct {
	stop *taggedCommit

	commits []*object.Commit
}

func collectMessages(iter object.CommitIter, stop *taggedCommit) ([]*object.Commit, error) {
	c := &commitMessageCollector{stop: stop}
	err := c.run(iter)

	return c.commits, err
}

func (c *commitMessageCollector) run(iter object.CommitIter) error {
	err := iter.ForEach(c.collect)
	switch err {
	case io.EOF:
		fallthrough
	case nil:
		return nil

	default:
		return err
	}
}

func (c *commitMessageCollector) collect(commit *object.Commit) error {
	if c.stop != nil && c.stop.IsCommit(commit) {
		return io.EOF
	}

	c.commits = append(c.commits, commit)

	return nil
}
