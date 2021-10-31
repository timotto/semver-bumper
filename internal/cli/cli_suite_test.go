package cli_test

import (
	"bytes"
	. "github.com/timotto/semver-bumper/internal/cli/clifakes"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cli Suite")
}

type outputRecorder struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

func newRecordingFakeOs(osArgs ...string) (*FakeOs, *outputRecorder) {
	f := &FakeOs{}
	f.ArgsReturns(append([]string{os.Args[0]}, osArgs...))

	r := &outputRecorder{}
	f.StdoutReturns(&r.Stdout)
	f.StderrReturns(&r.Stderr)

	return f, r
}
