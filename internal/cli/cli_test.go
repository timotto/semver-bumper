package cli_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/timotto/semver-bumper/internal/cli"
	. "github.com/timotto/semver-bumper/internal/cli/clifakes"
	. "github.com/timotto/semver-bumper/pkg/config"
	. "github.com/timotto/semver-bumper/pkg/test/git_testbed"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strings"
)

var _ = Describe("Run", func() {
	var (
		bed    *TestbedRepo
		fakeOs *FakeOs
		rec    *outputRecorder

		emptyTempDir string
	)
	BeforeEach(CreateBeforeEach(os.TempDir(), &bed))
	AfterEach(TeardownAfterEach(&bed))
	BeforeEach(func() {
		emptyTempDir = createAnEmptyTempDir()
	})
	AfterEach(func() {
		cleanupEmptyTempDir(emptyTempDir)
	})

	var runWithArgs = func(args ...string) error {
		fakeOs, rec = newRecordingFakeOs(args...)
		return Run(fakeOs)
	}

	It("prints the result to stdout", func() {
		err := runWithArgs(bed.Path(), "-0", "3.14.159")

		Expect(err).ToNot(HaveOccurred())
		Expect(rec.Stdout.String()).To(Equal("3.14.159\n"))
	})

	When("there is an error", func() {
		It("prints the error to stderr and returns it", func() {
			err := runWithArgs(bed.Path(), "-0", "bad-semver")

			Expect(err).To(HaveOccurred())
			var expectErrorDetails = func(errorMessage string) {
				Expect(errorMessage).To(ContainSubstring("bad-semver"))
				Expect(errorMessage).To(ContainSubstring("invalid initial version"))
			}
			Expect(rec.Stdout.String()).ToNot(ContainSubstring("bad-semver"))
			expectErrorDetails(rec.Stderr.String())
			expectErrorDetails(err.Error())
		})
	})

	Describe("--print-keywords", func() {
		It("prints the configured version bump level keywords", func() {
			err := runWithArgs(bed.Path(), "--print-keywords", "-1", "keyword-a", "-2", "keyword-b", "-3", "keyword-c")

			Expect(err).To(HaveOccurred())

			stderr := err.Error()
			Expect(stderr).To(ContainSubstring("keywords"))
			Expect(stderr).To(ContainSubstring("major"))
			Expect(stderr).To(ContainSubstring("minor"))
			Expect(stderr).To(ContainSubstring("patch"))
			Expect(stderr).To(ContainSubstring("keyword-a"))
			Expect(stderr).To(ContainSubstring("keyword-b"))
			Expect(stderr).To(ContainSubstring("keyword-c"))
		})
	})

	Describe("--output", func() {
		It("writes the result into a file instead of stdout", func() {
			filename := path.Join(emptyTempDir, "expected-file")

			err := runWithArgs(bed.Path(), "-0", "3.14.159", "--output", filename)

			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Stdout.String()).ToNot(ContainSubstring("3.14.159"))
			expectFileToHaveContent(filename, "3.14.159\n")
		})
	})

	Describe("--commits", func() {
		var filename string
		BeforeEach(func() {
			filename = path.Join(emptyTempDir, "expected-file")
		})
		It(`writes the commits into the given file formatted like "git log --format=oneline"`, func() {
			// given there are commits
			bed.AddCommits("commit message 1", "commit message 2", "commit with trailing LF\n")

			// when I want them to be stored in a file
			Expect(runWithArgs(bed.Path(), "--commits", filename)).ToNot(HaveOccurred())

			// then
			// the file content
			actualContent := fileContent(filename)
			// and the result of "git log --format=oneline > filename"
			expectedContent := formattedLikeLogFormatOneline(bed.Commits())
			// are the same
			Expect(actualContent).To(Equal(expectedContent))
			// and it's not an accident
			Expect(actualContent).ToNot(BeEmpty())
		})

		It(`includes only the commits relevant to the bump`, func() {
			bed.
				AddCommits("unexpected-1", "unexpected-2").
				AddLightweightTag("1.2.3").
				AddCommits("expected-1", "expected-2")

			err := runWithArgs(bed.Path(), "--commits", filename)

			Expect(err).ToNot(HaveOccurred())
			expectFileToContain(filename, "expected-1", "expected-2")
		})
	})

	Describe("config file", func() {
		const (
			expectedPrereleasePrefix = "expectedprereleaseprefix"
			expectedTagPrefix        = "expected_tag_prefix"
			expectedPathInclude1     = "expected_path_include_1"
			expectedPathInclude2     = "expected_path_include_2"
			expectedPathExclude1     = "expected_path_include_2/exclude_1"
			expectedPathExclude2     = "exclude_2/*/exclude_me"
			expectedInitialVersion   = "3.141.59265"
		)
		var filename string
		var runWithConfigFile = func(args ...string) error {
			allArgs := []string{"--config-file", filename, bed.Path()}
			return runWithArgs(append(allArgs, args...)...)
		}
		var writeAConfigFile = func(filename string, enc encoderFn, content *Options) {
			file, err := os.Create(filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(enc(file).Encode(content)).ToNot(HaveOccurred())
			Expect(file.Close()).ToNot(HaveOccurred())
		}
		var writeANewConfigFile = func(enc encoderFn, content *Options) func() {
			return func() {
				writeAConfigFile(filename, enc, content)
			}
		}

		BeforeEach(func() {
			filename = path.Join(emptyTempDir, "expected-file")
		})
		Describe("--write-config", func() {
			var expectedResult Options
			BeforeEach(func() {
				expectedResult = Options{
					InitialVersion: expectedInitialVersion,

					Prerelease: expectedPrereleasePrefix,
					TagPrefix:  expectedTagPrefix,

					PathInclude: []string{expectedPathInclude1, expectedPathInclude2},
					PathExclude: []string{expectedPathExclude1, expectedPathExclude2},
				}
				Expect(expectedResult.Valid()).ToNot(HaveOccurred())
			})
			var runWithWriteConfig = func() {
				err := runWithArgs(
					bed.Path(),
					"--write-config", filename,
					"--pre", expectedPrereleasePrefix,
					"--tag-prefix", expectedTagPrefix,
					"--initial-version", expectedInitialVersion,
					"--path-include", expectedPathInclude1,
					"--path-include", expectedPathInclude2,
					"--path-exclude", expectedPathExclude1,
					"--path-exclude", expectedPathExclude2,
				)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("written"))
				Expect(err.Error()).To(ContainSubstring(filename))
			}
			var expectFileToUnmarshalAndEqual = func(unmarshal func([]byte, interface{}) error) {
				data, err := os.ReadFile(filename)
				Expect(err).ToNot(HaveOccurred())

				actualResult := Options{}
				Expect(unmarshal(data, &actualResult)).ToNot(HaveOccurred())
				Expect(actualResult.Valid()).ToNot(HaveOccurred())

				Expect(actualResult).To(Equal(expectedResult))
			}
			When("the filename ends with .json", func() {
				BeforeEach(func() {
					filename = filename + ".json"
				})
				It("writes the supplied flags as JSON config file", func() {
					runWithWriteConfig()
					expectFileToUnmarshalAndEqual(json.Unmarshal)
				})
			})
			When("the filename ends with .yaml", func() {
				BeforeEach(func() {
					filename = filename + ".yaml"
				})
				It("writes the supplied flags as YAML config file", func() {
					runWithWriteConfig()
					expectFileToUnmarshalAndEqual(yaml.Unmarshal)
				})
			})
			When("the filename ends with neither .json nor .yaml", func() {
				BeforeEach(func() {
					filename = filename + ".different"
				})
				It("writes the supplied flags as YAML config file", func() {
					runWithWriteConfig()
					expectFileToUnmarshalAndEqual(yaml.Unmarshal)
				})
			})
		})
		Describe("--config-file", func() {
			Describe("file, suffix, and content", func() {
				givenConfigFileContent := Options{
					InitialVersion: expectedInitialVersion,
				}
				var expectToRunWithFlagsFromConfigFile = func() func() {
					return func() {
						Expect(runWithConfigFile()).ToNot(HaveOccurred())
						Expect(rec.Stdout.String()).To(Equal(expectedInitialVersion + "\n"))
					}
				}
				var itRunsWithTheFlagsFromTheGivenFile = func() {
					It("runs with the flags read from the given file", expectToRunWithFlagsFromConfigFile())
				}
				var itReturnsAnError = func() {
					It("returns an error", func() {
						Expect(runWithConfigFile()).To(HaveOccurred())
					})
				}
				var common = func(filenameSuffix string) {
					BeforeEach(func() {
						filename = filename + filenameSuffix
					})
					When("the content is JSON", func() {
						BeforeEach(writeANewConfigFile(jsonEncoder, &givenConfigFileContent))
						itRunsWithTheFlagsFromTheGivenFile()
					})
					When("the content is yaml", func() {
						BeforeEach(writeANewConfigFile(yamlEncoder, &givenConfigFileContent))
						itRunsWithTheFlagsFromTheGivenFile()
					})
					When("the content is neither JSON nor yaml", func() {
						BeforeEach(writeToFile(&filename, []byte("not yaml at all, so also no json")))
						itReturnsAnError()
					})
					When("the file does not exist", func() {
						itReturnsAnError()
					})
				}
				When("the filename ends with .json", func() {
					common(".json")
				})
				When("the filename ends with .yaml", func() {
					common(".yaml")
				})
				When("the filename ends with neither .json nor .yaml", func() {
					common(".something-else")
				})
			})
			When("there is a config file with flags", func() {
				givenConfigFileContent := Options{
					InitialVersion: expectedInitialVersion,
					TagPrefix:      expectedTagPrefix,
				}
				BeforeEach(writeANewConfigFile(yamlEncoder, &givenConfigFileContent))
				BeforeEach(func() {
					bed.
						AddCommits("commit-1").
						AddLightweightTag(expectedTagPrefix + "1.2.3").
						AddCommits("commit-2").
						AddLightweightTag(expectedTagPrefix + "1.2.4-" + expectedPrereleasePrefix + ".10").
						AddCommits("commit-3")
				})
				It("runs with the flags read from the config file", func() {
					Expect(runWithConfigFile()).ToNot(HaveOccurred())
					Expect(rec.Stdout.String()).To(Equal("1.2.3\n"))
				})
				When("there are also command line arguments", func() {
					It("uses both flags from the config file and the command line", func() {
						Expect(runWithConfigFile("--pre", expectedPrereleasePrefix)).ToNot(HaveOccurred())
						Expect(rec.Stdout.String()).To(Equal("1.2.4-" + expectedPrereleasePrefix + ".11\n"))
					})
					When("there is a value for a flag in both the config file and command line arguments", func() {
						It("prefers the value from the command line arguments", func() {
							// config file has
							// - initial-version
							// - tag-prefix
							// parameters have
							// - tag-prefix
							Expect(runWithConfigFile("--tag-prefix", "v")).ToNot(HaveOccurred())
							// there are no tags with the tag-prefix from the command line
							// it returns the initial version from the config file
							Expect(rec.Stdout.String()).To(Equal(expectedInitialVersion + "\n"))
						})
					})
				})
			})
		})
		Describe("project config file", func() {
			Describe("filename and content", func() {
				DescribeTable(
					"possible filenames",
					func(suffix string) {
						filename := ".semver-bumper.conf" + suffix
						writeAConfigFile(
							path.Join(bed.Path(), filename),
							yamlEncoder,
							&Options{InitialVersion: "9.8.7"},
						)
						Expect(runWithArgs(bed.Path())).ToNot(HaveOccurred())
						Expect(rec.Stdout.String()).To(Equal("9.8.7\n"))
					},
					Entry("no suffix", ""),
					Entry(".json", ".json"),
					Entry(".yaml", ".yaml"),
					Entry(".yml", ".yml"),
				)
				When("there are multiple project config files", func() {
					BeforeEach(func() {
						writeAConfigFile(
							path.Join(bed.Path(), ".semver-bumper.conf"),
							yamlEncoder,
							&Options{TagPrefix: "a"},
						)
						writeAConfigFile(
							path.Join(bed.Path(), ".semver-bumper.conf.json"),
							yamlEncoder,
							&Options{TagPrefix: "b"},
						)
						writeAConfigFile(
							path.Join(bed.Path(), ".semver-bumper.conf.yaml"),
							yamlEncoder,
							&Options{TagPrefix: "c"},
						)
					})
					It("returns an error", func() {
						Expect(runWithArgs(bed.Path())).To(HaveOccurred())
					})
				})
			})
			When("there is a project config file", func() {
				var projectConfig Options
				BeforeEach(func() {
					projectConfig = Options{
						InitialVersion: "987.65.4",
					}
					projectConfigFilename := path.Join(bed.Path(), ".semver-bumper.conf")
					writeAConfigFile(projectConfigFilename, yamlEncoder, &projectConfig)
				})
				It("runs with the flags from that config file", func() {
					Expect(runWithArgs(bed.Path())).ToNot(HaveOccurred())
					Expect(rec.Stdout.String()).To(Equal(projectConfig.InitialVersion + "\n"))
				})
				When("there also is a --config-file command line argument", func() {
					BeforeEach(writeANewConfigFile(yamlEncoder, &Options{TagPrefix: "v"}))
					It("ignores the flags from that config file", func() {
						Expect(runWithArgs(bed.Path(), "--config-file", filename)).ToNot(HaveOccurred())
						// given the initial version is set in project file
						// when it is not in the user supplied file
						// then the default initial version is returned
						Expect(rec.Stdout.String()).To(Equal("1.0.0\n"))
					})
				})
			})
		})
	})
})

func expectFileToContain(filename string, expectedContents ...string) {
	actual := fileContent(filename)
	for _, expected := range expectedContents {
		Expect(actual).To(ContainSubstring(expected))
	}
}

func expectFileToHaveContent(filename, expectedContent string) {
	Expect(fileContent(filename)).To(Equal(expectedContent))
}

func fileContent(filename string) string {
	actualResult, err := os.ReadFile(filename)
	Expect(err).ToNot(HaveOccurred())

	return string(actualResult)
}

func createAnEmptyTempDir() string {
	result, err := os.MkdirTemp(os.TempDir(), "semver-cli-test-*")
	Expect(err).ToNot(HaveOccurred())

	return result
}

func cleanupEmptyTempDir(dir string) {
	Expect(os.RemoveAll(dir)).ToNot(HaveOccurred())
}

func formattedLikeLogFormatOneline(commits []*object.Commit) string {
	var lines []string
	for _, commit := range commits {
		line := fmt.Sprintf("%s %s\n", commit.Hash.String(), commit.Message)
		lines = append(lines, line)
	}

	return strings.Join(lines, "")
}

type encoderFn func(file *os.File) encoder

type encoder interface {
	Encode(v interface{}) error
}

var jsonEncoder encoderFn = func(file *os.File) encoder {
	return json.NewEncoder(file)
}

var yamlEncoder encoderFn = func(file *os.File) encoder {
	return yaml.NewEncoder(file)
}

func writeToFile(filename *string, data []byte) func() {
	return func() {
		Expect(os.WriteFile(*filename, data, 0644)).ToNot(HaveOccurred())
	}
}
