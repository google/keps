package sigs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSigs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sigs Suite")
}
