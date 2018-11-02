package skeleton_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSkeleton(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Skeleton Suite")
}
