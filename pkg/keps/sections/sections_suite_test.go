package sections_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSections(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sections Suite")
}
