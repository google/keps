package rendering_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/section/internal/rendering"
)

var _ = Describe("The Graduation Criteria Section", func() {
	Describe("NewGraduationCriteria()", func() {
		It("renders a new Graduation Criteria", func() {
			info := newBasicRenderingInfo()
			content, err := rendering.NewGraduationCriteria(info)
			Expect(err).ToNot(HaveOccurred())

			Expect(content).To(ContainSubstring(basicInfoTitle))
			Expect(content).To(ContainSubstring("## Graduation Criteria"))
		})
	})
})
