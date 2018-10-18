package rendering_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"
)

var _ = Describe("The Operator Guide", func() {
	Describe("NewOperatorGuide()", func() {
		It("renders a new Operator Guide", func() {
			info := newBasicRenderingInfo()
			content, err := rendering.NewOperatorGuide(info)
			Expect(err).ToNot(HaveOccurred())

			Expect(content).To(ContainSubstring(basicInfoTitle))
			Expect(content).To(ContainSubstring("## Operator Guide"))
		})
	})
})
