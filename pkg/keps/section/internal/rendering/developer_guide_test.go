package rendering_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/section/internal/rendering"
)

var _ = Describe("The Developer Guide", func() {
	Describe("NewDeveloperGuide()", func() {
		It("renders a new Developer Guide", func() {
			info := newBasicRenderingInfo()
			content, err := rendering.NewDeveloperGuide(info)
			Expect(err).ToNot(HaveOccurred())

			Expect(err).ToNot(HaveOccurred())
			Expect(content).To(ContainSubstring(basicInfoTitle))
			Expect(content).To(ContainSubstring("## Developer Guide"))
		})
	})
})
