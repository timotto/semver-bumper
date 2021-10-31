package git_testbed

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	. "github.com/onsi/gomega"
	"os"
	"path"
	"time"
)

type TestbedRepo struct {
	path string
	repo *git.Repository
	time time.Time
}

func CreateBeforeEach(dir string, p **TestbedRepo) func() {
	return func() {
		Expect(p).ToNot(BeNil())
		repo := NewTestbedRepo(dir)
		*p = repo
	}
}

func TeardownAfterEach(p **TestbedRepo) func() {
	return func() {
		Expect(p).ToNot(BeNil())
		Expect(*p).ToNot(BeNil())
		(*p).Teardown()
	}
}

func NewTestbedRepo(baseDirectory string) *TestbedRepo {
	dir, err := os.MkdirTemp(baseDirectory, "repo-*")
	Expect(err).ToNot(HaveOccurred())

	repo, err := git.PlainInit(dir, false)
	Expect(err).ToNot(HaveOccurred())

	cfg, err := repo.Config()
	Expect(err).ToNot(HaveOccurred())

	cfg.User.Name = "Test User"
	cfg.User.Email = "Test-User@Example.Com"

	Expect(repo.SetConfig(cfg)).ToNot(HaveOccurred())

	return &TestbedRepo{
		path: dir,
		repo: repo,
		time: time.Now().Add(-240 * time.Hour),
	}
}

func (b *TestbedRepo) Teardown() {
	Expect(b.Path()).To(HavePrefix(os.TempDir()))
	Expect(b.Path()).ToNot(Equal(os.TempDir()))
	Expect(os.RemoveAll(b.Path())).To(BeNil())
}

func (b *TestbedRepo) Path() string {
	return b.path
}

func (b *TestbedRepo) AddCommits(messages ...string) *TestbedRepo {
	for _, message := range messages {
		b.AddCommit(message)
	}
	return b
}

func (b *TestbedRepo) AddCommit(message string) *TestbedRepo {
	filename := fmt.Sprintf("some-file-%d", time.Now().UnixMilli())
	return b.AddCommitAt(filename, message)
}

func (b *TestbedRepo) AddCommitAt(filename, message string) *TestbedRepo {
	fullPath := path.Join(b.path, filename)
	fullDir := path.Dir(fullPath)
	Expect(os.MkdirAll(fullDir, 0755)).ToNot(HaveOccurred())

	err := os.WriteFile(fullPath, []byte("some content"), 0640)
	Expect(err).ToNot(HaveOccurred())

	w, err := b.repo.Worktree()
	Expect(err).ToNot(HaveOccurred())

	_, err = w.Add(filename)
	Expect(err).ToNot(HaveOccurred())

	commit, err := w.Commit(message, &git.CommitOptions{
		Author:    b.aSignature(),
		Committer: b.aSignature(),
	})
	Expect(err).ToNot(HaveOccurred())

	_, err = b.repo.CommitObject(commit)
	Expect(err).ToNot(HaveOccurred())

	return b
}

func (b *TestbedRepo) AddAnnotatedTag(tag string) *TestbedRepo {
	head, err := b.repo.Head()
	Expect(err).ToNot(HaveOccurred())

	_, err = b.repo.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Tagger:  b.aSignature(),
		Message: tag,
	})
	Expect(err).ToNot(HaveOccurred())

	return b
}

func (b *TestbedRepo) AddLightweightTag(tag string) *TestbedRepo {
	head, err := b.repo.Head()
	Expect(err).ToNot(HaveOccurred())

	_, err = b.repo.CreateTag(tag, head.Hash(), nil)
	Expect(err).ToNot(HaveOccurred())

	return b
}

func (b *TestbedRepo) Commits() []*object.Commit {
	iter, err := b.repo.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
	})
	Expect(err).ToNot(HaveOccurred())

	var result []*object.Commit
	var collect = func(commit *object.Commit) error {
		result = append(result, commit)
		return nil
	}

	Expect(iter.ForEach(collect)).ToNot(HaveOccurred())

	return result
}

func (b *TestbedRepo) aSignature() *object.Signature {
	return &object.Signature{
		Name:  "A Name",
		Email: "address@example.com",
		When:  b.nextTimeNow(),
	}
}

func (b *TestbedRepo) nextTimeNow() time.Time {
	now := b.time
	b.time = now.Add(time.Second)

	return now
}
