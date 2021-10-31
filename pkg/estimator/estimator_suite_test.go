package estimator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEstimator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Estimator Suite")
}
