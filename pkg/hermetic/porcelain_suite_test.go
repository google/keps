package hermetic_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPorcelain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Porcelain Suite")
}
