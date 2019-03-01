package inplace_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInplace(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inplace Suite")
}
