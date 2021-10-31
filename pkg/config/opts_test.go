package config_test

import (
	"github.com/Masterminds/semver/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	. "github.com/timotto/semver-bumper/pkg/config"
)

var _ = Describe("Options", func() {
	Describe("Valid", func() {

		Describe("Default values", func() {
			var uut *Options
			BeforeEach(func() {
				uut = &Options{}
				Expect(uut.Valid()).ToNot(HaveOccurred())
			})
			Describe("InitialVersion", func() {
				It(`is "1.0.0" by default`, func() {
					Expect(uut.InitialVersion).To(Equal("1.0.0"))
				})
			})
			Describe("NoMatchBump", func() {
				It(`is "none" by default`, func() {
					Expect(uut.NoMatchBump).To(Equal("none"))
				})
			})
			Describe("Bump level Keyword pattern", func() {
				It(`is the semantic commit pattern`, func() {
					Expect(uut.KeywordsMajor).To(ConsistOf("^BREAKING CHANGE:"))
					Expect(uut.KeywordsMinor).To(ConsistOf("^feat:"))
					Expect(uut.KeywordsPatch).To(ConsistOf("^fix:", "^chore:"))
				})
			})
		})
		Describe("Validation", func() {
			DescribeTable(
				"InitialVersion",
				func(val string, expect types.GomegaMatcher) {
					uut := &Options{InitialVersion: val}
					Expect(uut.Valid()).To(expect)
				},
				Entry("valid semantic version", "1.2.3", BeNil()),
				Entry("another valid semantic version", "0.0.1", BeNil()),
				Entry("unexpected prefix", "v1.0.0", HaveOccurred()),
				Entry("random string", "random string", HaveOccurred()),
			)
			Describe("InitialVersion", func() {
				It("must be a semantic version", func() {
					uut := &Options{
						InitialVersion: "not a semver version",
					}
					Expect(uut.Valid()).To(HaveOccurred())
				})
			})
			DescribeTable(
				"NoMatchBump",
				func(val string, expect types.GomegaMatcher) {
					uut := &Options{NoMatchBump: val}
					Expect(uut.Valid()).To(expect)
				},
				Entry("valid value: patch", "patch", BeNil()),
				Entry("valid value: none", "none", BeNil()),
				Entry("invalid value", "other", HaveOccurred()),
			)
		})

		Describe("Value objects", func() {
			Describe("InitialVersion", func() {
				It("makes it available as value object", func() {
					uut := &Options{InitialVersion: "3.2.1"}
					Expect(uut.Valid()).ToNot(HaveOccurred())

					Expect(uut.InitialVersionValue()).To(Equal(semver.MustParse("3.2.1")))
				})
			})
			Describe("ConfigFile", func() {
				It("provides a bool to check if it is set", func() {
					uut := &Options{ConfigFile: ""}
					Expect(uut.Valid()).ToNot(HaveOccurred())
					Expect(uut.ReadConfigFile()).To(BeFalse())

					uut = &Options{ConfigFile: "some-file"}
					Expect(uut.Valid()).ToNot(HaveOccurred())
					Expect(uut.ReadConfigFile()).To(BeTrue())
				})
			})
			Describe("WriteConfig", func() {
				It("provides a bool to check if it is set", func() {
					uut := &Options{WriteConfig: ""}
					Expect(uut.Valid()).ToNot(HaveOccurred())
					Expect(uut.WriteConfigFile()).To(BeFalse())

					uut = &Options{WriteConfig: "some-file"}
					Expect(uut.Valid()).ToNot(HaveOccurred())
					Expect(uut.WriteConfigFile()).To(BeTrue())
				})
			})
			Describe("Prerelease", func() {
				It("provides a bool to check if it is set", func() {
					uut := &Options{Prerelease: ""}
					Expect(uut.Valid()).ToNot(HaveOccurred())
					Expect(uut.BumpPrerelease()).To(BeFalse())

					uut = &Options{Prerelease: "someprefix"}
					Expect(uut.Valid()).ToNot(HaveOccurred())
					Expect(uut.BumpPrerelease()).To(BeTrue())
				})
			})
			Describe("NoMatchBump", func() {
				It("makes it available as value object", func() {
					uut := &Options{NoMatchBump: "none"}
					Expect(uut.Valid()).ToNot(HaveOccurred())
					Expect(uut.NoMatchBumpValue()).To(Equal(FallbackStrategyNone))

					uut = &Options{NoMatchBump: "patch"}
					Expect(uut.Valid()).ToNot(HaveOccurred())
					Expect(uut.NoMatchBumpValue()).To(Equal(FallbackStrategyPatch))

				})
			})
			Describe("Commits", func() {
				It("provides a bool to check if it is set", func() {
					uut := &Options{Commits: ""}
					Expect(uut.Valid()).ToNot(HaveOccurred())
					Expect(uut.OutputCommits()).To(BeFalse())

					uut = &Options{Commits: "some-filename"}
					Expect(uut.Valid()).ToNot(HaveOccurred())
					Expect(uut.OutputCommits()).To(BeTrue())
				})
			})
		})
	})
})
