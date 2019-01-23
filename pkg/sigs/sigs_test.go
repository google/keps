package sigs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/calebamiles/keps/pkg/sigs"
)

var _ = Describe("the SIGs helper package", func() {
	Describe("All()", func() {
		It("contains each SIG specified in sigs.yaml of kubernetes/community@master", func() {
			upstreamList := fetchUpstreamSIGNames()
			compiledList := sigs.All()

			Expect(len(upstreamList)).To(Equal(len(compiledList)), "compiled list of SIGs should have the length as upstream")

			for _, s := range upstreamList {
				Expect(compiledList).To(ContainElement(s), "expected compiled SIG list to include: "+s)
			}
		})
	})

	Describe("Exists()", func() {
		It("returns whether a specific SIG exists", func() {
			upstreamList := fetchUpstreamSIGNames()

			for _, s := range upstreamList {
				Expect(sigs.Exists(s)).To(BeTrue(), "expected "+s+" to exist as listed SIG")
			}
		})
	})

	Describe("SubprojectExists()", func() {
		It("returns wether a specific subproject exists", func() {
			upstreamList := fetchUpstreamSIGList()

			for _, s := range upstreamList.SIGs {
				for _, sp := range s.Subprojects {
					Expect(sigs.SubprojectExists(sp.Name)).To(BeTrue(), "expected "+sp.Name+" to exist as subproject")
				}
			}
		})
	})

	Describe("BuildRoutingFromPath", func() {
		Context("when the path is at a SIG root", func() {
			It("returns SIG wide information", func() {
				contentRoot := "/home/user/workspace/keps/content/"
				givenPath := "sig-node/device-plugins"

				info, err := sigs.BuildRoutingFromPath(contentRoot, givenPath)
				Expect(err).ToNot(HaveOccurred(), "building routing info from a path relative to the content root should not return an error")

				Expect(info.OwningSIG()).To(Equal("sig-node"), "SIG node should be owner of KEPs in the `sig-node` directory")
				Expect(info.SIGWide()).To(BeTrue(), "a KEP at the root of a SIG dir should have `sig-wide` scope")
			})
		})

		Context("when the path includes a SIG and subproject", func() {
			It("returns SIG and subproject information", func() {
				contentRoot := "/home/user/workspace/keps/content/"
				givenPath := "sig-node/kubelet/dynamic-kubelet-configuration"

				info, err := sigs.BuildRoutingFromPath(contentRoot, givenPath)
				Expect(err).ToNot(HaveOccurred(), "building routing info from a path relative to the content root and including a subproject should not return an error")

				Expect(info.OwningSIG()).To(Equal("sig-node"), "SIG node should be owner of KEPs created in its subproject directories")
				Expect(info.SIGWide()).To(BeFalse(), "KEPs within a subproject directory should not be given `sig-wide` scope")
				Expect(info.AffectedSubprojects()).To(ContainElement("kubelet"), "the affected subprojects for a KEP created inside a subproject directory should contain the subproject")
			})
		})
	})

})

func fetchUpstreamSIGNames() []string {
	names := []string{}
	upstreamList := fetchUpstreamSIGList()
	for _, sig := range upstreamList.SIGs {
		names = append(names, sig.Name)
	}

	return names
}

func fetchUpstreamSIGList() *upstreamSIGList {
	resp, err := http.Get(upstreamSIGListURL)
	defer resp.Body.Close()

	Expect(err).ToNot(HaveOccurred(), "downloading canonical SIG list")

	respBytes, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred(), "reading HTTP response")

	sl := &upstreamSIGList{}
	err = yaml.Unmarshal(respBytes, sl)
	Expect(err).ToNot(HaveOccurred(), "unmarshalling SIG list")
	Expect(sl.SIGs).ToNot(BeEmpty(), "unmarshalled SIG list is empty")

	return sl
}

type upstreamSIGList struct {
	SIGs []upstreamSIGEntry `yaml:"sigs"`
}

type upstreamSIGEntry struct {
	Name        string                    `yaml:"dir"` // we actually want to look at what the SIG is called on disk
	Subprojects []upstreamSubprojectEntry `yaml:"subprojects"`
}

type upstreamSubprojectEntry struct {
	Name string `yaml:"name"`
}

const upstreamSIGListURL = "https://raw.githubusercontent.com/kubernetes/community/master/sigs.yaml"
