package sigs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"io/ioutil"

	"github.com/go-yaml/yaml"

	"github.com/calebamiles/keps/pkg/sigs"
)

var _ = Describe("the SIGs helper package", func() {
	Describe("All", func() {
		It("contains each SIG specified in sigs.yaml of kubernetes/community@master", func() {
			resp, err := http.Get(upstreamSIGListURL)
			defer resp.Body.Close()

			Expect(err).ToNot(HaveOccurred(), "downloading canonical SIG list")

			respBytes, err := ioutil.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred(), "reading HTTP response")

			sl := &upstreamSIGList{}
			err = yaml.Unmarshal(respBytes, sl)
			Expect(err).ToNot(HaveOccurred(), "unmarshalling SIG list")
			Expect(sl.SIGs).ToNot(BeEmpty(), "unmarshalled SIG list is empty")

			Expect(len(sl.SIGs)).To(Equal(len(sigs.All)), "committed list of SIGs has different length than upstream")

			for _, s := range sl.SIGs {
				Expect(sigs.All).To(ContainElement(s.Name))
			}

		})
	})
})


type upstreamSIGList struct {
	SIGs []upstreamSIGEntry `yaml:"sigs"`
}

type upstreamSIGEntry struct {
	Name sigs.SIG `yaml:"dir"`// we actually want to look at what the SIG is called on disk
}

const upstreamSIGListURL = "https://raw.githubusercontent.com/kubernetes/community/master/sigs.yaml"
