package extract_test

import (
	"io/ioutil"

	"github.com/calebamiles/keps/helpers/convert/internal/extract"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extracting the Metadata", func() {
	It("extracts metadata from a KEP", func() {
		kepBytes, err := ioutil.ReadFile(kep0000File)
		Expect(err).ToNot(HaveOccurred())

		meta, remaining, err := extract.Metadata(kepBytes)
		Expect(err).ToNot(HaveOccurred())

		expectedMetadata, err := ioutil.ReadFile(kep0000ExpectedMetadataFile)
		Expect(err).ToNot(HaveOccurred())

		Expect(string(meta)).To(Equal(string(expectedMetadata)))
		Expect(remaining).ToNot(BeEmpty())
	})

	Context("when a frontmatter separator is missing", func() {
		XIt("returns an error", func() {
			Fail("test not written")
		})
	})
})
