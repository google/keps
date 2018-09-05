package propose_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPropose(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Propose Suite")
}
