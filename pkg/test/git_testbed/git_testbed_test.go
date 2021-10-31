package git_testbed_test

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/timotto/semver-bumper/pkg/test/cli_testbed"
	. "github.com/timotto/semver-bumper/pkg/test/git_testbed"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var _ = Describe("GitTestbed", func() {
	var uut *TestbedRepo
	var dir string
	BeforeEach(createTempDirectory(&dir))
	AfterEach(cleanupTempDirectory(&dir))

	var runGit = func(args ...string) CliResult {
		allArgs := append([]string{"-C", uut.Path()}, args...)
		return RunCommand("git", allArgs...)
	}

	BeforeEach(func() {
		uut = NewTestbedRepo(dir)
	})

	Describe("NewTestbedRepo", func() {
		It("creates a new directory in the given baseDirectory with a defined prefix", func() {
			item := theItemIn(dir)
			Expect(item.IsDir()).To(BeTrue())
			Expect(item.Name()).To(HavePrefix("repo-"))
		})

		It("remembers the directory name", func() {
			actualPath := uut.Path()
			expectedPath := path.Join(dir, theItemIn(dir).Name())
			Expect(actualPath).To(Equal(expectedPath))
		})

		It("creates an empty working git repo in that directory", func() {
			statusResult := runGit("status")
			statusResult.ExpectSuccess()
			statusResult.ExpectInOutput("On branch", "No commits")

			runGit("log").
				ExpectError()

			runGit("tag").
				ExpectSuccess().
				ExpectOutput("")

			Expect(os.WriteFile(path.Join(uut.Path(), "file"), []byte("content"), 0640)).ToNot(HaveOccurred())
			runGit("add", "file").
				ExpectSuccess()
			runGit("commit", "-m", "test").
				ExpectSuccess()

			runGit("log").
				ExpectSuccess()
		})
	})

	Describe("Teardown", func() {
		It("deletes the directory", func() {
			uut.Teardown()
			Expect(itemsIn(dir)).To(BeEmpty())
		})
	})

	Describe("AddCommits", func() {
		It("adds commits", func() {
			firstMessage := "first commit message"
			secondMessage := "a different message for another commit"
			uut.AddCommits(firstMessage, secondMessage)
			result := runGit("log", "--format=oneline")
			result.ExpectSuccess()
			result.ExpectInOutput(firstMessage, secondMessage)
		})
		It("sets the committer time apart by at least one second", func() {
			uut.AddCommits("1", "2", "3", "4", "5")

			result := runGit("log", "--format=%ct")
			result.ExpectSuccess()
			timestamps := parseUnixTimestampFromEachLine(result.Output)
			Expect(timestamps).To(HaveLen(5))
			intervals := secondsBetweenTimestamps(timestamps)
			Expect(intervals).To(HaveLen(4))
			for _, interval := range intervals {
				Expect(interval.Seconds()).To(BeNumerically(">=", 1))
			}
		})
	})

	Describe("AddCommitAt", func() {
		It("adds a commit with a specific filename", func() {
			expectedFilename := "some-directory/some-filename"
			firstMessage := "first commit message"
			uut.AddCommitAt(expectedFilename, firstMessage)
			result := runGit("log", "--format=oneline")
			result.ExpectSuccess()
			result.ExpectInOutput(firstMessage)

			Expect(filenamesIn(uut.Path())).To(ContainElement("some-directory"))
			Expect(filenamesIn(path.Join(uut.Path(), "some-directory"))).To(ContainElement("some-filename"))
		})
	})

	Describe("AddLightweightTag", func() {
		It("adds a lightweight tag to the head", func() {
			firstMessage := "first commit message"
			secondMessage := "commit-with-tag"
			tagName := "a-tag-name"
			uut.AddCommits(firstMessage, secondMessage)

			uut.AddLightweightTag(tagName)

			result := runGit("for-each-ref", "refs/tags")
			result.ExpectSuccess()
			lines := strings.Split(result.Output, "\n")
			Expect(lines).To(HaveLen(2)) // last is blank

			columns := strings.Split(lines[0], " ")
			Expect(columns).To(HaveLen(2))
			columns = strings.Split(columns[1], "\t")
			Expect(columns).To(Equal([]string{"commit", "refs/tags/" + tagName}))
		})
	})

	Describe("AddAnnotatedTag", func() {
		It("adds an annotated tag to the head", func() {
			firstMessage := "first commit message"
			secondMessage := "commit-with-tag"
			tagName := "a-tag-name"
			uut.AddCommits(firstMessage, secondMessage)

			uut.AddAnnotatedTag(tagName)

			result := runGit("for-each-ref", "refs/tags")
			result.ExpectSuccess()
			lines := strings.Split(result.Output, "\n")
			Expect(lines).To(HaveLen(2)) // last is blank

			columns := strings.Split(lines[0], " ")
			Expect(columns).To(HaveLen(2))
			columns = strings.Split(columns[1], "\t")
			Expect(columns).To(Equal([]string{"tag", "refs/tags/" + tagName}))
		})
	})

	Describe("Commits", func() {
		It(`returns all the commits like "git log --format=oneline"`, func() {
			// given
			firstMessage := "first commit message"
			secondMessage := "a different message for another commit"
			uut.AddCommits(firstMessage, secondMessage)
			result := runGit("log", "--format=oneline")
			result.ExpectSuccess()

			// when
			commits := uut.Commits()

			// then
			Expect(formattedLikeLogFormatOneline(commits)).
				To(Equal(result.Output))
		})
	})

	Describe("Test Helpers", func() {
		var bed *TestbedRepo
		BeforeEach(CreateBeforeEach(dir, &bed))
		AfterEach(TeardownAfterEach(&bed))

		Describe("CreateBeforeEach", func() {
			var lastInstance string
			It("creates a new instance", func() {
				anotherInstance := fmt.Sprintf("%v", bed)
				Expect(anotherInstance).ToNot(Equal(lastInstance))
				lastInstance = anotherInstance
			})
			It("creates a new instance", func() {
				anotherInstance := fmt.Sprintf("%v", bed)
				Expect(anotherInstance).ToNot(Equal(lastInstance))
				lastInstance = anotherInstance
			})
		})
	})
})

func theItemIn(dir string) fs.FileInfo {
	items := itemsIn(dir)
	Expect(items).To(HaveLen(1))
	return items[0]
}

func itemsIn(dir string) []fs.FileInfo {
	items, err := ioutil.ReadDir(dir)
	Expect(err).ToNot(HaveOccurred())

	return items
}

func filenamesIn(dir string) []string {
	items := itemsIn(dir)

	var result []string
	for _, item := range items {
		result = append(result, item.Name())
	}

	return result
}

func formattedLikeLogFormatOneline(commits []*object.Commit) string {
	var lines []string
	for _, commit := range commits {
		line := fmt.Sprintf("%s %s\n", commit.Hash.String(), commit.Message)
		lines = append(lines, line)
	}

	return strings.Join(lines, "")
}

func parseUnixTimestampFromEachLine(output string) []time.Time {
	var result []time.Time

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		val, err := strconv.ParseInt(line, 10, 32)
		Expect(err).ToNot(HaveOccurred())

		t := time.Unix(val, 0)
		result = append(result, t)
	}

	return result
}

func secondsBetweenTimestamps(timestamps []time.Time) []time.Duration {
	var result []time.Duration
	var moreRecentTime time.Time
	for i, timestamp := range timestamps {
		if i == 0 {
			moreRecentTime = timestamp
			continue
		}

		delta := moreRecentTime.Sub(timestamp)
		moreRecentTime = timestamp

		result = append(result, delta)
	}

	return result
}
