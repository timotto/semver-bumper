package cli

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"strings"
)

func (rt runtime) beforeResult() error {
	if rt.opts.PrintKeywords {
		return rt.printKeywords()
	}

	if rt.opts.WriteConfigFile() {
		return rt.writeConfigToFile()
	}

	return nil
}

func (rt runtime) onResult(version *semver.Version, commits []*object.Commit) error {
	if err := rt.outputVersion(version); err != nil {
		return err
	}

	if err := rt.outputCommits(commits); err != nil {
		return err
	}

	return nil
}

func (rt runtime) outputVersion(version *semver.Version) error {
	if rt.opts.Output == "" {
		Outln(rt.os, version.String())
		return nil
	}

	if err := os.WriteFile(rt.opts.Output, []byte(version.String()+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to write to %v: %w", rt.opts.Output, err)
	}

	return nil
}

func (rt runtime) outputCommits(commits []*object.Commit) error {
	if !rt.opts.OutputCommits() {
		return nil
	}

	data := []byte(format(commits))
	if err := os.WriteFile(rt.opts.Commits, data, 0644); err != nil {
		return fmt.Errorf("failed to write to %v: %w", rt.opts.Commits, err)
	}

	return nil
}

func format(commits []*object.Commit) string {
	var lines []string
	for _, commit := range commits {
		msg := commit.Message
		for strings.HasSuffix(msg, "\n") {
			msg = strings.TrimSuffix(msg, "\n")
		}
		line := fmt.Sprintf("%s %s\n", commit.Hash.String(), msg)
		lines = append(lines, line)
	}

	return strings.Join(lines, "")
}

func (rt runtime) printKeywords() error {
	buf := &bytes.Buffer{}
	_, _ = fmt.Fprintln(buf, "the keywords for the different version bump level are:")
	var fn = func(label string, keywords []string) {
		_, _ = fmt.Fprintf(buf, "%s:\n", label)
		for _, key := range keywords {
			_, _ = fmt.Fprintf(buf, "\t%s\n", key)
		}
	}

	fn("major", rt.opts.KeywordsMajor)
	fn("minor", rt.opts.KeywordsMinor)
	fn("patch", rt.opts.KeywordsPatch)

	return fmt.Errorf(buf.String())
}

func (rt runtime) writeConfigToFile() error {
	if err := rt.opts.WriteToFile(rt.opts.WriteConfig); err != nil {
		return err
	} else {
		return fmt.Errorf("configuration written to %v", rt.opts.WriteConfig)
	}
}
