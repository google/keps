package planctae_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPlanctae(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Planctae Suite")
}
