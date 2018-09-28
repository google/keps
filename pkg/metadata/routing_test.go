package metadata_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/metadata"
)

var _ = Describe("Routing Metadata", func() {
	Describe("NewRoutingFromPath()", func() {
		Context("when the SIG and subproject exist", func() {
			It("returns routing info with SIG and subproject", func() {
				p := "content/sig-node/kubelet/device-plugins"

				routingInfo, err := metadata.NewRoutingFromPath(p)
				Expect(err).ToNot(HaveOccurred())

				Expect(routingInfo.OwningSIG).To(Equal("sig-node"))
				Expect(routingInfo.AffectedSubprojects).To(HaveLen(1))
				Expect(routingInfo.AffectedSubprojects).To(ContainElement("kubelet"))
			})
		})

		Context("when the SIG does not exist", func() {
			It("returns an error", func() {
				p := "content/sig-not-real/kublet/device-plugins"

				_, err := metadata.NewRoutingFromPath(p)
				Expect(err.Error()).To(ContainSubstring("no SIG information"))
			})
		})

		Context("when the SIG exists but the subproject does not", func() {
			It("returns routing info with SIG", func() {
				By("treating subproject information as optional")

				p := "content/sig-node/not-a-subproject/lost-idea"

				routingInfo, err := metadata.NewRoutingFromPath(p)
				Expect(err).ToNot(HaveOccurred())

				Expect(routingInfo.OwningSIG).To(Equal("sig-node"))
				Expect(routingInfo.AffectedSubprojects).To(BeEmpty())
			})
		})
	})
})
