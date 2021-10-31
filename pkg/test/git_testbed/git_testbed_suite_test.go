package git_testbed_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGitTestbed(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GitTestbed Suite")
}

func createTempDirectory(p *string) func() {
	return func() {
		Expect(p).ToNot(BeNil())
		dir, err := os.MkdirTemp(os.TempDir(), "git-testbed-*")
		Expect(err).ToNot(HaveOccurred())

		*p = dir
	}
}

func cleanupTempDirectory(p *string) func() {
	return func() {
		Expect(p).ToNot(BeNil())
		dir := *p
		Expect(dir).To(HavePrefix(os.TempDir()))
		Expect(dir).ToNot(Equal(os.TempDir()))
		Expect(os.RemoveAll(dir)).ToNot(HaveOccurred())
	}
}
