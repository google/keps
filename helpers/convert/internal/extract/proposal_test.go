package extract_test

import (
	"bytes"
	"io/ioutil"

	"github.com/calebamiles/keps/helpers/convert/internal/extract"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extracting the Proposal", func() {
	It("extracts the Proposal from a KEP", func() {
		kepBytes, err := ioutil.ReadFile(kep0000File)
		Expect(err).ToNot(HaveOccurred())

		sections, err := extract.Sections(kepBytes)
		Expect(err).ToNot(HaveOccurred())

		proposal := sections[extract.ProposalHeading]

		expectedProposalBytes, err := ioutil.ReadFile(kep0000ExpectedProposalFile)
		Expect(err).ToNot(HaveOccurred())

		expectedProposal := bytes.TrimSpace(expectedProposalBytes)

		Expect(string(proposal)).To(Equal(string(expectedProposal)))
	})
})
