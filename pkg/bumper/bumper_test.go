package bumper_test

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing/object"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/timotto/semver-bumper/pkg/bumper"
	. "github.com/timotto/semver-bumper/pkg/config"
	"github.com/timotto/semver-bumper/pkg/estimator"
	"github.com/timotto/semver-bumper/pkg/gitrepo"
	. "github.com/timotto/semver-bumper/pkg/test/git_testbed"
	"os"
)

const (
	testInitialVersion    = "9.8.7"
	beforeInitialVersion  = "0.1.2"
	afterInitialVersion   = "11.12.13"
	testPrereleasePrefix  = "almost"
	testInitialPrerelease = testInitialVersion + "-" + testPrereleasePrefix + ".1"

	patchLevelCommitMessage = "fix: bug"
	minorLevelCommitMessage = "feat: feature"
	majorLevelCommitMessage = "BREAKING CHANGE: all new"
)

var _ = Describe("Bumper", func() {
	var (
		cfg  *Options
		esti Estimator
		repo GitRepo
		bed  *TestbedRepo
	)

	BeforeEach(CreateBeforeEach(os.TempDir(), &bed))
	AfterEach(TeardownAfterEach(&bed))

	BeforeEach(func() {
		cfg = aConfiguration()
		esti = anEstimator()
		repo = aGitRepo(bed)
	})

	var bedWith = func(fns ...func(testbedRepo *TestbedRepo)) func() {
		return func() {
			for _, fn := range fns {
				fn(bed)
			}
		}
	}
	var expectVersion = func(expectedVersion string, expectedCommitMessages ...interface{}) func() {
		return func() {
			actualResult, actualCommits, err := Bump(cfg, repo, esti)
			Expect(err).ToNot(HaveOccurred())
			Expect(actualResult).ToNot(BeNil())
			Expect(actualResult.String()).To(Equal(expectedVersion))
			if len(expectedCommitMessages) > 0 {
				Expect(messagesFrom(actualCommits...)).To(ConsistOf(expectedCommitMessages...))
			}
		}
	}

	When("BumpPrerelease==false", func() {
		BeforeEach(func() {
			cfg.Prerelease = ""
		})

		const expectedInitialVersion = testInitialVersion
		var itReturnsTheInitialVersion = func(commits ...interface{}) {
			It("returns the initial version",
				expectVersion(expectedInitialVersion, commits...))
		}

		When("there are no commits", func() {
			itReturnsTheInitialVersion()
		})

		When("there are commits", func() {
			BeforeEach(bedWith(commits("one", "two")))

			When("there already is a tag with a release version", func() {
				const existingReleaseVersion = beforeInitialVersion
				var expectBump = func(expectedBump func(version *semver.Version) *semver.Version, expectedCommitMessages ...interface{}) func() {
					return expectVersion(bumped(existingReleaseVersion, expectedBump), expectedCommitMessages...)
				}

				BeforeEach(bedWith(lightweightTags(existingReleaseVersion)))

				When("there is a major level commit", func() {
					BeforeEach(bedWith(commits(majorLevelCommitMessage)))
					It("bumps the major level", expectBump(majorBump, majorLevelCommitMessage))
				})

				When("there is a minor level commit", func() {
					BeforeEach(bedWith(commits(minorLevelCommitMessage)))
					It("bumps the minor level", expectBump(minorBump, minorLevelCommitMessage))

					When("there is also a major level commit", func() {
						BeforeEach(bedWith(commits(majorLevelCommitMessage)))
						It("bumps the major level", expectBump(majorBump, majorLevelCommitMessage, minorLevelCommitMessage))
					})
				})

				When("there is a patch level commit", func() {
					BeforeEach(bedWith(commits(patchLevelCommitMessage)))
					It("bumps the patch level", expectBump(patchBump, patchLevelCommitMessage))

					When("there is also a minor level commit", func() {
						BeforeEach(bedWith(commits(minorLevelCommitMessage)))
						It("bumps the minor level", expectBump(minorBump, minorLevelCommitMessage, patchLevelCommitMessage))

						When("there is also a major level commit", func() {
							BeforeEach(bedWith(commits(majorLevelCommitMessage)))
							It("bumps the major level", expectBump(majorBump, majorLevelCommitMessage, minorLevelCommitMessage, patchLevelCommitMessage))
						})
					})

					When("there is also a major level commit", func() {
						BeforeEach(bedWith(commits(majorLevelCommitMessage)))
						It("bumps the major level", expectBump(majorBump, majorLevelCommitMessage, patchLevelCommitMessage))
					})
				})

				When("there is a minor level commit", func() {
					BeforeEach(bedWith(commits(minorLevelCommitMessage)))
					It("bumps the minor level", expectBump(minorBump, minorLevelCommitMessage))

					When("there is also a major level commit", func() {
						BeforeEach(bedWith(commits(majorLevelCommitMessage)))
						It("bumps the major level", expectBump(majorBump, majorLevelCommitMessage, minorLevelCommitMessage))
					})
				})
			})

			When("there are no tags", func() {
				itReturnsTheInitialVersion("one", "two")
			})

			When("there is a prerelease below the initial version", func() {
				BeforeEach(func() {
					bed.AddLightweightTag(asPrerelease1(beforeInitialVersion))
				})
				itReturnsTheInitialVersion("one", "two")
			})

			When("there is a prerelease above the initial version", func() {
				BeforeEach(func() {
					bed.AddLightweightTag(asPrerelease1(afterInitialVersion))
				})
				itReturnsTheInitialVersion("one", "two")
			})
		})
	})

	When("BumpPrerelease==true", func() {
		BeforeEach(func() {
			cfg.Prerelease = testPrereleasePrefix
		})

		const expectedInitialVersion = testInitialPrerelease
		var itReturnsTheInitialVersion = func(commits ...interface{}) {
			It("returns the initial version as prerelease # 1",
				expectVersion(expectedInitialVersion, commits...))
		}

		When("there are no commits", func() {
			itReturnsTheInitialVersion()
		})

		When("there are commits but no tags", func() {
			BeforeEach(func() {
				bed.AddCommits("one", "two")
			})
			itReturnsTheInitialVersion("one", "two")
		})

		When("there are commits with a prerelease tag before the initial version", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag(asPrerelease(beforeInitialVersion, 999)).
					AddCommit("three")
			})
			It("returns the existing prerelease + 1",
				expectVersion(asPrerelease(beforeInitialVersion, 1000)))
		})
		When("there are commits with prerelease tags after the initial version", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one", "two").
					AddLightweightTag(asPrerelease(afterInitialVersion, 5)).
					AddCommit("three")
			})
			It("returns the existing prerelease + 1",
				expectVersion(asPrerelease(afterInitialVersion, 6)))
		})
		When("there is a release, a prerelease, and commits", func() {
			BeforeEach(func() {
				bed.
					AddCommits("one").
					AddLightweightTag("1.1.0").
					AddCommits("fix: bug").
					AddLightweightTag("1.1.1-almost.1").
					AddCommit("fix: bug")
			})
			It("returns the existing prerelease + 1",
				expectVersion("1.1.1-almost.2"))
		})
	})
})

func aConfiguration() *Options {
	cfg := &Options{InitialVersion: testInitialVersion}
	Expect(cfg.Valid()).ToNot(HaveOccurred())

	return cfg
}

func anEstimator() Estimator {
	eCfg := &Options{Prerelease: testPrereleasePrefix}
	Expect(eCfg.Valid()).ToNot(HaveOccurred())

	return estimator.NewEstimator(eCfg)
}

func aGitRepo(bed *TestbedRepo) GitRepo {
	rCfg := &Options{}
	Expect(rCfg.Valid()).ToNot(HaveOccurred())
	repo, err := gitrepo.NewGitRepo(rCfg, bed.Path())
	Expect(err).ToNot(HaveOccurred())

	return repo
}

func asPrerelease1(version string) string {
	return asPrerelease(version, 1)
}

func asPrerelease(version string, n int) string {
	return fmt.Sprintf("%s-%s.%d", version, testPrereleasePrefix, n)
}

func messagesFrom(commits ...*object.Commit) []string {
	var result []string

	for _, commit := range commits {
		result = append(result, commit.Message)
	}

	return result
}

func commits(commits ...string) func(repo *TestbedRepo) {
	return func(bed *TestbedRepo) {
		bed.AddCommits(commits...)
	}
}

func lightweightTags(tags ...string) func(repo *TestbedRepo) {
	return func(bed *TestbedRepo) {
		for _, tag := range tags {
			bed.AddLightweightTag(tag)
		}
	}
}

func bumped(version string, bump func(*semver.Version) *semver.Version) string {
	v, err := semver.StrictNewVersion(version)
	Expect(err).ToNot(HaveOccurred())
	return bump(v).String()
}

func majorBump(v *semver.Version) *semver.Version {
	b := v.IncMajor()
	return &b
}

func minorBump(v *semver.Version) *semver.Version {
	b := v.IncMinor()
	return &b
}

func patchBump(v *semver.Version) *semver.Version {
	b := v.IncPatch()
	return &b
}
