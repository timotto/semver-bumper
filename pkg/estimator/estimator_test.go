package estimator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/timotto/semver-bumper/pkg/config"
	. "github.com/timotto/semver-bumper/pkg/estimator"
	. "github.com/timotto/semver-bumper/pkg/model"
)

const (
	testMajorKeyword = "major"
	testMinorKeyword = "minor"
	testPatchKeyword = "patch"
	testNoKeyword    = "no-keyword"

	testPrereleasePrefix = "prelive"
)

var _ = Describe("Estimator", func() {
	Describe("BumpLevelFrom", func() {
		var config *config.Options
		var commits []string
		var expectBump = func(expected BumpLevel) func() {
			return func() {
				Expect(NewEstimator(config).BumpLevelFrom(commits)).To(Equal(expected))
			}
		}
		var common = func(expectedDefaultLevel BumpLevel) {
			DescribeTable(
				"bumps based on commits",
				func(expected BumpLevel, c ...string) {
					commits = c
					expectBump(expected)()
				},
				Entry("no commits -> default", expectedDefaultLevel),
				Entry("no matching commits -> default", expectedDefaultLevel, testNoKeyword+"1", testNoKeyword+"2"),

				Entry("commits matching patch -> patch", BumpLevelPatch, testPatchKeyword, testNoKeyword),
				Entry("commits matching minor -> minor", BumpLevelMinor, testMinorKeyword, testNoKeyword),
				Entry("commits matching major -> minor", BumpLevelMajor, testMajorKeyword, testNoKeyword),

				Entry("commits matching patch & minor -> minor", BumpLevelMinor, testPatchKeyword, testMinorKeyword),
				Entry("commits matching patch & major -> major", BumpLevelMajor, testPatchKeyword, testMajorKeyword),
				Entry("commits matching minor & major -> major", BumpLevelMajor, testMinorKeyword, testMajorKeyword),
				Entry("commits matching patch & minor & major -> major", BumpLevelMajor, testPatchKeyword, testMinorKeyword, testMajorKeyword),
			)
		}

		When(`the NoMatchBump configuration is "none"`, func() {
			BeforeEach(func() {
				config = aConfiguration()
				config.NoMatchBump = "none"
				Expect(config.Valid()).NotTo(HaveOccurred())
			})

			common(BumpLevelNone)
		})

		When(`the NoMatchBump configuration is "patch"`, func() {
			BeforeEach(func() {
				config = aConfiguration()
				config.NoMatchBump = "patch"
				Expect(config.Valid()).NotTo(HaveOccurred())
			})

			common(BumpLevelPatch)
		})
	})

	Describe("NextPrerelease", func() {
		DescribeTable(
			"behavior",
			func(givenInput, expectedOutput string) {
				Expect(
					NewEstimator(aConfiguration()).
						NextPrerelease(givenInput)).
					To(
						Equal(expectedOutput))
			},
			Entry("empty input returns pre#1", "", testPrereleasePrefix+".1"),
			Entry("input pre#1 returns pre#2", testPrereleasePrefix+".1", testPrereleasePrefix+".2"),
			Entry("input pre#1000 returns pre#1001", testPrereleasePrefix+".1000", testPrereleasePrefix+".1001"),
		)

		When("the input does not match the prerelease prefix configuration", func() {
			It("returns an error", func() {
				_, err := NewEstimator(aConfiguration()).NextPrerelease("badpre.1")
				Expect(err).To(HaveOccurred())

				_, err = NewEstimator(aConfiguration()).NextPrerelease(testPrereleasePrefix)
				Expect(err).To(HaveOccurred())

				_, err = NewEstimator(aConfiguration()).NextPrerelease(testPrereleasePrefix + ".")
				Expect(err).To(HaveOccurred())

				_, err = NewEstimator(aConfiguration()).NextPrerelease(testPrereleasePrefix + ".1")
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})

func aConfiguration() *config.Options {
	cfg := &config.Options{
		KeywordsMajor: []string{testMajorKeyword},
		KeywordsMinor: []string{testMinorKeyword},
		KeywordsPatch: []string{testPatchKeyword},
		Prerelease:    testPrereleasePrefix,
	}

	Expect(cfg.Valid()).ToNot(HaveOccurred())
	return cfg
}
