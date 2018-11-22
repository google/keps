package extract_test

import (
	"bytes"
	"io/ioutil"

	"github.com/calebamiles/keps/helpers/convert/internal/extract"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extracting the Motivation", func() {
	It("extracts the Motivation from a KEP", func() {
		kepBytes, err := ioutil.ReadFile(kep0000File)
		Expect(err).ToNot(HaveOccurred())

		sections, err := extract.Sections(kepBytes)
		Expect(err).ToNot(HaveOccurred())

		motivation := sections[extract.MotivationHeading]

		expectedMotivationBytes, err := ioutil.ReadFile(kep0000ExpectedMotivationFile)
		Expect(err).ToNot(HaveOccurred())

		expectedMotivation := bytes.TrimSpace(expectedMotivationBytes)

		Expect(string(motivation)).To(Equal(string(expectedMotivation)))
	})
})
