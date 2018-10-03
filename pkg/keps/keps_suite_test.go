package keps_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKeps(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keps Suite")
}
