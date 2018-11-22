package extract_test

import (
	"io/ioutil"
	"strings"

	"github.com/calebamiles/keps/helpers/convert/internal/extract"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("extracting the Graduation Criteria", func() {
	It("extracts the Graduation Criteria from a KEP", func() {
		kepBytes, err := ioutil.ReadFile(kep0000File)
		Expect(err).ToNot(HaveOccurred())

		sections, err := extract.Sections(kepBytes)
		Expect(err).ToNot(HaveOccurred())

		graduationCriteria := sections[extract.GraduationCriteriaHeading]

		expectedGraduationCriteriaBytes, err := ioutil.ReadFile(kep0000ExpectedGraduationCriteriaFile)
		Expect(err).ToNot(HaveOccurred())

		expectedGraduationCriteria := strings.TrimSpace(string(expectedGraduationCriteriaBytes))

		Expect(string(graduationCriteria)).To(Equal(expectedGraduationCriteria))
	})
})
