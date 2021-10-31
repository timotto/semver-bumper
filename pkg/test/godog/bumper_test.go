package godog_test

import (
	"github.com/cucumber/godog"
	"github.com/onsi/gomega"
	"testing"
)

func TestFeatures(t *testing.T) {
	gomega.RegisterFailHandler(func(message string, _ ...int) {
		panic(message)
	})

	suite := godog.TestSuite{
		ScenarioInitializer: initializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func initializeScenario(ctx *godog.ScenarioContext) {

	bum := &bumperFeature{}

	ctx.Before(bum.before)
	ctx.After(bum.after)

	ctx.Step(`^I see the version (.+)$`, bum.iSeeTheVersion)
	ctx.Step(`^there is a commit "([^"]*)"$`, bum.thereIsACommit)
	ctx.Step(`^there is a commit "([^"]*)" with the tag (.*)$`, bum.thereIsACommitWithTheTag)
	ctx.Step(`^there is a directory with a git repository$`, bum.thereIsADirectoryWithAGitRepository)

	ctx.Step(`^I run semver-bumper (.+)$`, bum.iRunSemverBumper)
	ctx.Step(`^I run semver-bumper$`, bum.iRunSemverBumperInThatDirectory)
	ctx.Step(`^I tag the git with (.+)$`, bum.iTagTheGitWith)

	ctx.Step(`^I see the help page$`, bum.iSeeTheHelpPage)
	ctx.Step(`^the exit code is (\d+)$`, bum.theExitCodeIs)
}
