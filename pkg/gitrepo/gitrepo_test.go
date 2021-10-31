package gitrepo_test

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing/object"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	. "github.com/timotto/semver-bumper/pkg/config"
	. "github.com/timotto/semver-bumper/pkg/gitrepo"
	. "github.com/timotto/semver-bumper/pkg/test/git_testbed"
	"os"
)

var _ = Describe("Gitrepo", func() {
	var (
		uut *Gitrepo
		bed *TestbedRepo
	)
	BeforeEach(CreateBeforeEach(os.TempDir(), &bed))
	AfterEach(TeardownAfterEach(&bed))
	var aUnitUnderTest = func(with ...func(configuration *Options)) func() {
		return func() {
			var err error
			uut, err = NewGitRepo(aConfig(with...), bed.Path())
			Expect(err).ToNot(HaveOccurred())
		}
	}
	BeforeEach(aUnitUnderTest())
	Describe("LatestTaggedRelease", func() {
		var expectVersion = func(expectedVersion string) {
			actualResult, err := uut.LatestTaggedRelease()
			Expect(err).ToNot(HaveOccurred())
			Expect(actualResult.String()).To(Equal(expectedVersion))
		}

		When("the latest commit is a release tag", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag("1.2.3")
			})
			It("returns the version", func() {
				expectVersion("1.2.3")
			})
		})

		When("there are commits after a release tag", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag("1.2.3").
					AddCommits("some", "more")
			})
			It("returns the version", func() {
				expectVersion("1.2.3")
			})
		})

		When("there are multiple commits with release tags", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag("1.2.3").
					AddCommits("some", "more").
					AddLightweightTag("0.9.2")
			})
			It("returns the latest version", func() {
				expectVersion("1.2.3")
			})
		})

		When("there is a prerelease tag after a release tag", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag("1.2.3").
					AddCommits("some", "more").
					AddLightweightTag("1.2.4-rc.5")
			})
			It("returns the release tag", func() {
				expectVersion("1.2.3")
			})
		})

		It("respects tag filters", func() {
			cfg := aConfig()
			cfg.TagPrefix = "t"
			Expect(cfg.Valid()).ToNot(HaveOccurred())
			uut, err := NewGitRepo(cfg, bed.Path())

			bed.
				AddCommits("very unexpected").
				AddLightweightTag("5.0.1").
				AddCommits("expected").
				AddLightweightTag("t1.2.3").
				AddCommits("unexpected").
				AddLightweightTag("5.0.2")

			actualResult, err := uut.LatestTaggedRelease()
			Expect(err).ToNot(HaveOccurred())
			Expect(actualResult.String()).To(Equal("1.2.3"))
		})

		It("respects annotated tags", func() {
			bed.
				AddCommits("unexpected").
				AddLightweightTag("1.0.1").
				AddCommits("expected").
				AddAnnotatedTag("1.2.3")

			actualResult, err := uut.LatestTaggedRelease()
			Expect(err).ToNot(HaveOccurred())
			Expect(actualResult.String()).To(Equal("1.2.3"))
		})

		It("respects lightweight tags", func() {
			bed.
				AddCommits("very unexpected").
				AddAnnotatedTag("1.0.1").
				AddCommits("expected").
				AddLightweightTag("1.2.3")

			actualResult, err := uut.LatestTaggedRelease()
			Expect(err).ToNot(HaveOccurred())
			Expect(actualResult.String()).To(Equal("1.2.3"))
		})

		Describe("edge cases", func() {
			var expectNil = func() {
				It("returns nil and no error 🙀", func() {
					actualResult, err := uut.LatestTaggedRelease()
					Expect(err).ToNot(HaveOccurred())
					Expect(actualResult).To(BeNil())
				})
			}
			When("there is no commit", expectNil)
			When("there are commits but no tag", func() {
				BeforeEach(func() {
					bed.AddCommits("something", "something else")
				})
				expectNil()
			})
			When("there is a prerelease tag", func() {
				BeforeEach(func() {
					bed.
						AddCommits("something", "something else").
						AddLightweightTag("1.2.3-rc.4")
				})
				expectNil()
			})
		})
	})
	Describe("LatestTaggedPrerelease", func() {
		var expectVersion = func(expectedVersion string) {
			actualResult, err := uut.LatestTaggedPrerelease()
			Expect(err).ToNot(HaveOccurred())
			Expect(actualResult.String()).To(Equal(expectedVersion))
		}

		When("the latest commit is a release tag", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag("1.2.3")
			})
			It("returns the version", func() {
				expectVersion("1.2.3")
			})
		})

		When("there are commits after a release tag", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag("1.2.3").
					AddCommits("some", "more")
			})
			It("returns the version", func() {
				expectVersion("1.2.3")
			})
		})

		When("there are multiple commits with release tags", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag("1.2.3").
					AddCommits("some", "more").
					AddLightweightTag("0.9.2")
			})
			It("returns the latest version", func() {
				expectVersion("1.2.3")
			})
		})

		When("there is a prerelease tag after a release tag", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag("1.2.3").
					AddCommits("some", "more").
					AddLightweightTag("1.2.4-rc.5")
			})
			It("returns the prerelease tag", func() {
				expectVersion("1.2.4-rc.5")
			})
		})

		When("there is a prerelease tag", func() {
			BeforeEach(func() {
				bed.
					AddCommits("something", "something else").
					AddLightweightTag("1.2.4-rc.5")
			})
			It("returns the prerelease tag", func() {
				expectVersion("1.2.4-rc.5")
			})
		})

		It("respects tag filters", func() {
			cfg := aConfig()
			cfg.TagPrefix = "t"
			Expect(cfg.Valid()).ToNot(HaveOccurred())
			uut, err := NewGitRepo(cfg, bed.Path())

			bed.
				AddCommits("very unexpected").
				AddLightweightTag("5.0.1-rc.1").
				AddCommits("expected").
				AddLightweightTag("t1.2.3-rc.5").
				AddCommits("unexpected").
				AddLightweightTag("5.0.1-rc.2")

			actualResult, err := uut.LatestTaggedPrerelease()
			Expect(err).ToNot(HaveOccurred())
			Expect(actualResult.String()).To(Equal("1.2.3-rc.5"))
		})

		Describe("edge cases", func() {
			var expectNil = func() {
				It("returns nil and no error 🙀", func() {
					actualResult, err := uut.LatestTaggedRelease()
					Expect(err).ToNot(HaveOccurred())
					Expect(actualResult).To(BeNil())
				})
			}
			When("there is no commit", expectNil)
			When("there are commits but no tag", func() {
				BeforeEach(func() {
					bed.AddCommits("something", "something else")
				})
				expectNil()
			})
		})
	})
	Describe("CommitMessagesSince", func() {
		It("returns the commit messages that happened after the commit with the given version", func() {
			// setup
			const (
				specificVersion = "1.1.0"
				expectedCommit1 = "expected1"
				expectedCommit2 = "expected2"
			)
			givenVersion := semver.MustParse(specificVersion)
			consistOfExpectedCommitMessages := ConsistOf(expectedCommit1, expectedCommit2)

			// given
			bed.
				// there are some commits with some tag
				AddCommits("unexpected1", "unexpected2").
				AddLightweightTag("1.0.0").
				// and there are some commits with the given tag
				AddCommits("unexpected3", "unexpected4").
				AddLightweightTag(specificVersion).
				// and there are some commits after the given tag
				AddCommits(expectedCommit1, expectedCommit2)

			// when
			actualCommits, err := uut.CommitMessagesSince(givenVersion)

			// then
			Expect(err).ToNot(HaveOccurred())
			Expect(messagesFrom(actualCommits...)).To(consistOfExpectedCommitMessages)
		})

		It("returns the commit messages that happened after the commit with that tag with prefix", func() {
			cfg := aConfig()
			cfg.TagPrefix = "test-v"
			uut, err := NewGitRepo(cfg, bed.Path())
			Expect(err).ToNot(HaveOccurred())
			bed.
				AddCommits("unexpected1", "unexpected2").
				AddLightweightTag("test-v1.0.0").
				AddCommits("unexpected3", "unexpected4").
				AddLightweightTag("test-v1.1.0").
				AddCommits("expected1", "expected2")

			actualCommits, err := uut.CommitMessagesSince(semver.MustParse("1.1.0"))
			Expect(err).ToNot(HaveOccurred())
			Expect(messagesFrom(actualCommits...)).To(ConsistOf("expected1", "expected2"))
		})

		Describe("path filters", func() {
			BeforeEach(func() {
				bed.
					AddCommitAt("unexpected", "unexpected").
					AddCommitAt("root-0", "commit-0").
					AddLightweightTag("1.0.0").
					AddCommitAt("root-1", "commit-1").
					AddCommitAt("root-2", "commit-2").
					AddCommitAt("first/first-1", "commit-3").
					AddCommitAt("first/first-2", "commit-4").
					AddCommitAt("first/second/second-1", "commit-5").
					AddCommitAt("third/third-1", "commit-6").
					AddCommitAt("third/first/third-2", "commit-7").
					AddCommitAt("forth/unexpected", "commit-8")
			})

			var expectCommits = func(index ...int) func() {
				return func() {
					actualCommits, err := uut.CommitMessagesSince(semver.MustParse("1.0.0"))
					Expect(err).ToNot(HaveOccurred())
					Expect(messagesFrom(actualCommits...)).To(consistOfCommits(index...))
				}
			}

			When("there is a include path filter", func() {
				var includeFilter = withIncludeFilters("first/**")
				BeforeEach(aUnitUnderTest(includeFilter))

				It("only finds commits with changes matching the path filter",
					expectCommits(3, 4, 5))

				When("there is another path filter", func() {
					var anotherIncludeFilter = withIncludeFilters("third")
					BeforeEach(aUnitUnderTest(includeFilter, anotherIncludeFilter))

					It("only finds commits with changes matching any of the path filters",
						expectCommits(3, 4, 5, 6, 7))
				})

				When("there is also a ignore path filter", func() {
					var excludeFilter = withExcludeFilters("first/second/*")
					BeforeEach(aUnitUnderTest(includeFilter, excludeFilter))

					It("ignores the commits with changes matching the ignore path filter",
						expectCommits(3, 4))
				})
			})

			When("there is an exclude path filter", func() {
				var excludeFilter = withExcludeFilters("first/*")
				BeforeEach(aUnitUnderTest(excludeFilter))

				It("ignores the commits with changes matching the ignore path filter",
					expectCommits(1, 2, 6, 7, 8))

				When("there is another exclude path filter", func() {
					var anotherExcludeFilter = withExcludeFilters("forth/unexpected")
					BeforeEach(aUnitUnderTest(excludeFilter, anotherExcludeFilter))

					It("ignores all commits with changes matching any of the ignore path filter",
						expectCommits(1, 2, 6, 7))
				})

				When(`an exclude filter matches the given "since" parameter`, func() {
					BeforeEach(aUnitUnderTest(excludeFilter, withExcludeFilters("root-0")))
					It("returns the commits as expected",
						expectCommits(1, 2, 6, 7, 8))
				})
			})
		})

		When("there are no commits", func() {
			It("returns an empty result", func() {
				actualCommits, err := uut.CommitMessagesSince(semver.MustParse("1.2.3"))
				Expect(err).ToNot(HaveOccurred())
				Expect(actualCommits).To(BeEmpty())
			})
		})

		When("there is no matching tag", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag("1.0.1").
					AddCommits("one", "two").
					AddLightweightTag("4.5.6").
					AddCommits("three")
			})

			It("returns an empty result", func() {
				actualCommits, err := uut.CommitMessagesSince(semver.MustParse("1.2.3"))
				Expect(err).ToNot(HaveOccurred())
				Expect(actualCommits).To(BeEmpty())
			})
		})

		When("the given version is nil 🙀", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag("1.2.3").
					AddCommits("three", "four")
			})
			It("returns all the commits", func() {
				actualCommits, err := uut.CommitMessagesSince(nil)
				Expect(err).ToNot(HaveOccurred())
				Expect(messagesFrom(actualCommits...)).To(ConsistOf("one", "two", "three", "four"))
			})
		})
	})

	When("the directory is not a git repo", func() {
		It("returns an error", func() {
			dir, err := os.MkdirTemp(os.TempDir(), "gitrepo-test-*")
			Expect(err).ToNot(HaveOccurred())
			defer func() {
				Expect(os.RemoveAll(dir)).ToNot(HaveOccurred())
			}()

			cfg := aConfig()
			_, err = NewGitRepo(cfg, dir)
			Expect(err).To(HaveOccurred())
		})
	})
})

func aConfig(with ...func(configuration *Options)) *Options {
	cfg := &Options{}

	for _, fn := range with {
		fn(cfg)
	}
	Expect(cfg.Valid()).To(BeNil())

	return cfg
}

func withIncludeFilters(filter ...string) func(p *Options) {
	return func(p *Options) {
		p.PathInclude = append(p.PathInclude, filter...)
	}
}

func withExcludeFilters(filter ...string) func(p *Options) {
	return func(p *Options) {
		p.PathExclude = append(p.PathExclude, filter...)
	}
}

func consistOfCommits(index ...int) types.GomegaMatcher {
	var result []interface{}
	for _, i := range index {
		result = append(result, fmt.Sprintf("commit-%v", i))
	}
	return ConsistOf(result...)
}

func messagesFrom(commits ...*object.Commit) []string {
	var result []string

	for _, commit := range commits {
		result = append(result, commit.Message)
	}

	return result
}
