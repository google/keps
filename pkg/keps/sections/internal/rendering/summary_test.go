package rendering_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"
)

var _ = Describe("The Summary Section", func() {
	Describe("NewSummary()", func() {
		It("renders a new Summary", func() {
			info := newBasicRenderingInfo()

			content, err := rendering.NewSummary(info)
			Expect(err).ToNot(HaveOccurred())

			Expect(content).To(ContainSubstring(basicInfoTitle))
			Expect(content).To(ContainSubstring("## Summary"), "expected `Summary` section heading to exist")
		})
	})
})
