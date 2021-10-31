package gitrepo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestGitrepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gitrepo Suite")
}
