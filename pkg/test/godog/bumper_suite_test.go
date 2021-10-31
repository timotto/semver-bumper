package godog_test

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	. "github.com/timotto/semver-bumper/pkg/test/cli_testbed"
	. "github.com/timotto/semver-bumper/pkg/test/git_testbed"
	"os"
	"strings"
)

type bumperFeature struct {
	res CliResult
	git *TestbedRepo
}

func (f *bumperFeature) iSeeTheVersion(expectedVersion string) error {
	if f.res.Output != expectedVersion+"\n" {
		return fmt.Errorf("expected to see verion %v but the actual output is:\n%v", expectedVersion, f.res.Output)
	}

	return nil
}

func (f *bumperFeature) thereIsACommit(givenCommitMessage string) error {
	f.git.AddCommit(givenCommitMessage)

	return nil
}

func (f *bumperFeature) thereIsACommitWithTheTag(givenCommitMessage, givenTag string) error {
	f.git.AddCommit(givenCommitMessage).AddLightweightTag(givenTag)

	return nil
}

func (f *bumperFeature) iTagTheGitWith(tag string) error {
	f.git.AddLightweightTag(tag)

	return nil
}

func (f *bumperFeature) thereIsADirectoryWithAGitRepository() error {
	f.git = NewTestbedRepo(os.TempDir())

	return nil
}

func (f *bumperFeature) iRunSemverBumperInThatDirectory() error {
	return f.runSemverBumper()
}

func (f *bumperFeature) iRunSemverBumper(args string) error {
	splitArgs := strings.Split(args, " ")

	return f.runSemverBumper(splitArgs...)
}

func (f *bumperFeature) runSemverBumper(splitArgs ...string) error {
	goArgs := []string{"run", "../../../cmd/bumper"}

	if f.git != nil {
		goArgs = append(goArgs, f.git.Path())
	}

	f.res = RunCommand(""+
		"go",
		append(goArgs, splitArgs...)...,
	)

	return nil
}

func (f *bumperFeature) iSeeTheHelpPage() error {
	but := fmt.Sprintf("but is actually:\\n%v", f.res.Output)

	if !strings.HasPrefix(f.res.Output, "Usage:") {
		return fmt.Errorf(`help output should start with "Usage:"` + but)
	}

	if !strings.Contains(f.res.Output, "Application Options:") {
		return fmt.Errorf(`help output should start with "Usage:"` + but)
	}

	return nil
}

func (f *bumperFeature) theExitCodeIs(expectedExitCode int) error {
	if f.res.ExitCode != expectedExitCode {
		return fmt.Errorf("expected exit code %v but actual result is %v", expectedExitCode, f.res.ExitCode)
	}

	return nil
}

func (f *bumperFeature) before(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
	f.clear()
	return ctx, nil
}

func (f *bumperFeature) after(ctx context.Context, _ *godog.Scenario, err error) (context.Context, error) {
	f.clear()

	return ctx, err
}

func (f *bumperFeature) clear() {
	f.res = CliResult{}
	if f.git != nil {
		f.git.Teardown()
		f.git = nil
	}
}
