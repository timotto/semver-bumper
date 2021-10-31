package cli_testbed_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCliTestbed(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CliTestbed Suite")
}
