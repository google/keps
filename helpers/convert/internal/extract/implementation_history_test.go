package extract_test

import (
	"io/ioutil"
	"strings"

	"github.com/calebamiles/keps/helpers/convert/internal/extract"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extracting the Implementation History section", func() {
	It("extracts the Implementation History section", func() {
		kepBytes, err := ioutil.ReadFile(kep0000File)
		Expect(err).ToNot(HaveOccurred())

		sections, err := extract.Sections(kepBytes)
		Expect(err).ToNot(HaveOccurred())

		implementationHistory := sections[extract.ImplementationHistoryHeading]

		expectedImplementationHistoryBytes, err := ioutil.ReadFile(kep0000ExpectedImplementationHistoryFile)
		Expect(err).ToNot(HaveOccurred())

		expectedImplementationHistory := strings.TrimSpace(string(expectedImplementationHistoryBytes))
		Expect(string(implementationHistory)).To(Equal(expectedImplementationHistory))
	})
})
