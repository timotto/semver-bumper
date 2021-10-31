package gitrepo

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing"
)

type (
	collector struct {
		strict bool
		repo   *git.Repository
		prefix string
		Result []taggedCommit
	}
	taggedCommit struct {
		Tag *semver.Version
		Ref *object.Commit
	}
	collection []taggedCommit
)

func (g Gitrepo) newCollector(strict bool) *collector {
	return &collector{
		strict: strict,
		repo:   g.repo,
		prefix: g.conf.TagPrefix,
	}
}

func (c *collector) collect(ref *plumbing.Reference) error {
	ok, tag := c.hasTagPrefix(ref.Name().Short())
	if !ok {
		return nil
	}

	v, err := semver.StrictNewVersion(tag)
	if err != nil {
		return fmt.Errorf("failed to parse version [%v]: %w", tag, err)
	}

	if c.strict {
		if v.Prerelease() != "" || v.Metadata() != "" {
			return nil
		}
	}

	commit, err := c.resolveCommit(ref, v)
	if err != nil {
		return fmt.Errorf("failed to resolve commit for %v: %w", tag, err)
	}

	item := taggedCommit{
		Tag: v,
		Ref: commit,
	}
	c.Result = append(c.Result, item)

	return nil
}

func (c collector) resolveCommit(ref *plumbing.Reference, v *semver.Version) (*object.Commit, error) {
	tagObject, err := c.repo.TagObject(ref.Hash())
	switch err {
	case plumbing.ErrObjectNotFound:
		return c.repo.CommitObject(ref.Hash())
	case nil:
		return tagObject.Commit()
	default:
		return nil, fmt.Errorf("failed to resolve tag %v: %w", v, err)
	}
}

func (c collector) hasTagPrefix(tag string) (bool, string) {
	if len(c.prefix) == 0 {
		return true, tag
	}

	if !strings.HasPrefix(tag, c.prefix) {
		return false, ""
	}

	return true, strings.TrimPrefix(tag, c.prefix)
}

func (t taggedCommit) IsVersion(v *semver.Version) bool {
	return t.Tag.Equal(v)
}

func (t taggedCommit) IsCommit(commit *object.Commit) bool {
	return t.Ref.Hash == commit.Hash
}

func (n collection) Len() int {
	return len(n)
}

func (n collection) Less(i, j int) bool {
	return n[i].Tag.LessThan(n[j].Tag)
}

func (n collection) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n collection) Latest() taggedCommit {
	return n[len(n)-1]
}
