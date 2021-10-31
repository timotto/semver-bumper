package cli_testbed_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/timotto/semver-bumper/pkg/test/cli_testbed"
)

var _ = Describe("CliTestbed", func() {
	Describe("RunCommand", func() {
		It("understands the result code", func() {
			result := RunCommand("true")
			Expect(result.ExitCode).To(Equal(0))
			result.ExpectSuccess()

			result = RunCommand("false")
			Expect(result.ExitCode).To(Equal(1))
			result.ExpectError()
		})

		It("understands arguments and output", func() {
			command := "echo"
			argument1 := "argument-1"
			argument2 := "argument-2"
			argument3 := "argument-n"
			expectedOutput := "argument-1 argument-2 argument-n\n"

			result := RunCommand(command, argument1, argument2, argument3)

			Expect(result.Output).To(Equal(expectedOutput))
			result.ExpectOutput(expectedOutput)
			result.ExpectInOutput(argument1)
			result.ExpectInOutput(argument2)
			result.ExpectInOutput(argument3)
		})

		It("mixes stdout and stderr", func() {
			givenScript := `
echo on-stdout-1
echo on-stderr-1 1>&2
echo on-stdout-2
echo on-stderr-2 1>&2
`
			expectedOutput := `on-stdout-1
on-stderr-1
on-stdout-2
on-stderr-2
`

			result := RunCommand("sh", "-c", givenScript)

			result.ExpectOutput(expectedOutput)
		})
	})
})
