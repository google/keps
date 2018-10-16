package rendering_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/section/internal/rendering"
)

var _ = Describe("The Teacher Guide", func() {
	Describe("NewTeacherGuide()", func() {
		It("renders a new Teacher Guide", func() {
			info := newBasicRenderingInfo()
			content, err := rendering.NewTeacherGuide(info)
			Expect(err).ToNot(HaveOccurred())

			Expect(content).To(ContainSubstring(basicInfoTitle))
			Expect(content).To(ContainSubstring("## Teacher Guide"))
		})
	})
})
