package enhancements_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEnhancements(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Enhancements Suite")
}
