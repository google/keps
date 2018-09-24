package metadata_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routing Metadata", func() {
	Describe("extracting routing information from a path string", func() {
		Context("with a valid SIG and subproject", func() {
			It("extracts SIG information from a path string", func() {
				Expect(sigs.ExtractFromPath("content/sig-node/kubelet")).To(Equal("sig-node"))
			})

			It("extracts subproject information from a path string", func() {
				Expect(sigs.ExtractSubprojectFromPath("content/sig-node/kubelet")).To(Equal("sig-node"))
			})
		})

		Context("with an invalid SIG and/or subproject", func() {
			It("returns an empty string for the invalid component", func() {
				By("parsing the SIG")
				Expect(sigs.ExtractFromPath("content/sig-node/not-real-subproject")).To(Equal("sig-node"))
				Expect(sigs.ExtractFromPath("content/sig-not-real/kubelet")).To(BeEmpty())

				By("parsing the subproject")
				Expect(sigs.ExtractSubprojectFromPath("content/sig-not-real/kubelet")).To(Equal("kubelet"))
				Expect(sigs.ExtractSubprojectFromPath("content/sig-not-real/not-real-subproject")).To(BeEmpty())
			})
		})

		Context("with a missing SIG and/or subproject", func() {
			It("returns an empty string for the missing component", func() {

			})
		})
	})
})
