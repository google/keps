package extract_test

import (
	"bytes"
	"io/ioutil"

	"github.com/calebamiles/keps/helpers/convert/internal/extract"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extracting the Drawbacks section", func() {
	It("extracts the Drawbacks section", func() {
		kepBytes, err := ioutil.ReadFile(kep0000File)
		Expect(err).ToNot(HaveOccurred())

		sections, err := extract.Sections(kepBytes)
		Expect(err).ToNot(HaveOccurred())

		drawbacks := sections[extract.DrawbacksHeading]

		expectedDrawbacksBytes, err := ioutil.ReadFile(kep0000ExpectedDrawbacksFile)
		Expect(err).ToNot(HaveOccurred())

		expectedDrawbacks := bytes.TrimSpace(expectedDrawbacksBytes)

		Expect(string(drawbacks)).To(Equal(string(expectedDrawbacks)))
	})
})
