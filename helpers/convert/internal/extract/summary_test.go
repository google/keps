package extract_test

import (
	"bytes"
	"io/ioutil"

	"github.com/calebamiles/keps/helpers/convert/internal/extract"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extracting the Summary", func() {
	It("extracts the Summary of a KEP", func() {
		kepContent, err := ioutil.ReadFile(kep0000File)
		Expect(err).ToNot(HaveOccurred())

		sections, err := extract.Sections(kepContent)
		Expect(err).ToNot(HaveOccurred())

		summary := sections[extract.SummaryHeading]

		expectedSummaryBytes, err := ioutil.ReadFile(kep0000ExpectedSummaryFile)
		Expect(err).ToNot(HaveOccurred())

		expectedSummary := bytes.TrimSpace(expectedSummaryBytes)

		Expect(string(summary)).To(Equal(string(expectedSummary)))
	})
})
