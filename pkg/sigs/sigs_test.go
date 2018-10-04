package sigs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"net/http"

	"github.com/go-yaml/yaml"

	"github.com/calebamiles/keps/pkg/sigs"
)

var _ = Describe("the SIGs helper package", func() {
	Describe("All()", func() {
		It("contains each SIG specified in sigs.yaml of kubernetes/community@master", func() {
			upstreamList := fetchUpstreamSIGNames()
			compiledList := sigs.All()

			Expect(len(upstreamList)).To(Equal(len(compiledList)), "compiled list of SIGs has different length than upstream")

			for _, s := range upstreamList {
				Expect(compiledList).To(ContainElement(s))
			}
		})
	})

	Describe("Exists()", func() {
		It("returns whether a specific SIG exists", func() {
			upstreamList := fetchUpstreamSIGNames()

			for _, s := range upstreamList {
				Expect(sigs.Exists(s)).To(BeTrue())
			}
		})
	})

	Describe("BuildRoutingFromPath", func() {
		Context("when the path is the content root", func() {
			It("returns kubernetes wide information", func() {
				contentRoot := "/home/user/workspace/keps/content/"
				givenPath := "/home/user/workspace/keps/content"

				info, err := sigs.BuildRoutingFromPath(contentRoot, givenPath)
				Expect(err).ToNot(HaveOccurred())

				Expect(info.OwningSIG).To(Equal("sig-architecture"))
				Expect(info.KubernetesWide).To(BeTrue())
			})
		})

		Context("when the path is at a SIG root", func() {
			It("returns SIG wide information", func() {
				contentRoot := "/home/user/workspace/keps/content/"
				givenPath := "/home/user/workspace/keps/content/sig-node"

				info, err := sigs.BuildRoutingFromPath(contentRoot, givenPath)
				Expect(err).ToNot(HaveOccurred())

				Expect(info.OwningSIG).To(Equal("sig-node"))
				Expect(info.SIGWide).To(BeTrue())
			})
		})

		Context("when the path does not include a SIG", func() {
			It("returns an error", func() {
				contentRoot := "/home/user/workspace/keps/content/"
				givenPath := "/home/user/workspace/keps/content/sig-not-real"

				_, err := sigs.BuildRoutingFromPath(contentRoot, givenPath)
				Expect(err.Error()).To(ContainSubstring("unable to determine SIG information for given path"))
			})
		})

		Context("when the path includes a SIG and subproject", func() {
			It("returns SIG and subproject information", func() {
				contentRoot := "/home/user/workspace/keps/content/"
				givenPath := "/home/user/workspace/keps/content/sig-node/kubelet"

				info, err := sigs.BuildRoutingFromPath(contentRoot, givenPath)
				Expect(err).ToNot(HaveOccurred())

				Expect(info.OwningSIG).To(Equal("sig-node"))
				Expect(info.SIGWide).To(BeFalse())
				Expect(info.AffectedSubprojects).To(ContainElement("kubelet"))
			})
		})
	})

})

func fetchUpstreamSIGNames() []string {
	resp, err := http.Get(upstreamSIGListURL)
	defer resp.Body.Close()

	Expect(err).ToNot(HaveOccurred(), "downloading canonical SIG list")

	respBytes, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred(), "reading HTTP response")

	sl := &upstreamSIGList{}
	err = yaml.Unmarshal(respBytes, sl)
	Expect(err).ToNot(HaveOccurred(), "unmarshalling SIG list")
	Expect(sl.SIGs).ToNot(BeEmpty(), "unmarshalled SIG list is empty")

	names := []string{}
	for i := range sl.SIGs {
		names = append(names, sl.SIGs[i].Name)
	}

	return names
}

type upstreamSIGList struct {
	SIGs []upstreamSIGEntry `yaml:"sigs"`
}

type upstreamSIGEntry struct {
	Name string `yaml:"dir"` // we actually want to look at what the SIG is called on disk
}

const upstreamSIGListURL = "https://raw.githubusercontent.com/kubernetes/community/master/sigs.yaml"
