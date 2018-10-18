package rendering_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"
)

var _ = Describe("The Motivation Section", func() {
	Describe("NewMotivation()", func() {
		It("renders a new Motivation", func() {
			info := newBasicRenderingInfo()
			content, err := rendering.NewMotivation(info)
			Expect(err).ToNot(HaveOccurred())

			Expect(content).To(ContainSubstring(basicInfoTitle))
			Expect(content).To(ContainSubstring("## Motivation"))
		})
	})
})
