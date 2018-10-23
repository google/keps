package rendering_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRendered(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rendering Suite")
}
