package cli_testbed

import (
	"bytes"
	"fmt"
	. "github.com/onsi/gomega"
	"os/exec"
)

type CliResult struct {
	Command  string
	Args     []string
	ExitCode int
	Error    error
	Output   string
}

func RunCommand(command string, args ...string) CliResult {
	out := &bytes.Buffer{}
	cmd := exec.Command(command, args...)
	cmd.Stdout = out
	cmd.Stderr = out

	err := cmd.Run()

	return CliResult{
		Command:  command,
		Args:     args,
		ExitCode: cmd.ProcessState.ExitCode(),
		Error:    err,
		Output:   out.String(),
	}
}

func (r CliResult) AsString() string {
	return fmt.Sprintf("%v", &r)
}

func (r CliResult) ExpectSuccess() CliResult {
	Expect(r.ExitCode).To(Equal(0), r.AsString())
	return r
}

func (r CliResult) ExpectError() CliResult {
	Expect(r.ExitCode).ToNot(Equal(0), r.AsString())
	return r
}

func (r CliResult) ExpectOutput(expectedOutput string) CliResult {
	Expect(r.Output).To(Equal(expectedOutput))

	return r
}

func (r CliResult) ExpectInOutput(substrs ...string) CliResult {
	for _, substr := range substrs {
		Expect(r.Output).To(ContainSubstring(substr), r.AsString())
	}
	return r
}
