package keps_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps"
)

var _ = Describe("resolving the abosulte path to a KEP", func() {
	Describe("Path()", func() {
		It("resolves the absolute path for a top level KEP", func() {
			contentRoot := "/tmp/keps-sandbox"
			targetDir := "large-value-delivered-incrementally"

			p, err := keps.Path(contentRoot, targetDir)
			Expect(err).ToNot(HaveOccurred(), "building absolute path to KEP directory")

			Expect(p).To(Equal(filepath.Join(contentRoot, "kubernetes-wide", targetDir)))

			By("returning an unmodified fully qualified path if given")

			p2, err := keps.Path(contentRoot, p)
			Expect(p2).To(Equal(p))
		})

		It("resolves the absolute path for a SIG level KEP", func() {
			contentRoot := "/tmp/keps-sandbox"
			targetDir := "sig-node/functional-value-delivered-incrementally"
			kepName := filepath.Base(targetDir)

			p, err := keps.Path(contentRoot, targetDir)
			Expect(err).ToNot(HaveOccurred(), "building absolute path to KEP directory")

			Expect(p).To(Equal(filepath.Join(contentRoot, "sig-node", "sig-wide", kepName)))

			By("returning an unmodified fully qualified path if given")

			p2, err := keps.Path(contentRoot, p)
			Expect(p2).To(Equal(p))
		})

		It("resolves the absolute path for a subproject level KEP", func() {
			contentRoot := "/tmp/keps-sandbox"
			targetDir := "sig-node/kubelet/subproject-specific-value-delivered-incrementally"
			kepName := filepath.Base(targetDir)

			p, err := keps.Path(contentRoot, targetDir)
			Expect(err).ToNot(HaveOccurred(), "building absolute path to KEP directory")

			Expect(p).To(Equal(filepath.Join(contentRoot, "sig-node", "kubelet", kepName)))

			By("returning an unmodified fully qualified path if given")

			p2, err := keps.Path(contentRoot, p)
			Expect(p2).To(Equal(p))
		})

		Context("when SIG information cannot be determiend", func() {
			It("returns an error", func() {
				contentRoot := "/tmp/keps-sandbox"
				targetDir := "not-a-sig/unsponsored-idea"

				_, err := keps.Path(contentRoot, targetDir)
				Expect(err.Error()).To(ContainSubstring("unable to determine SIG information for given path"))
			})
		})
	})
})
