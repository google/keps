package rendering_test

import (
	"github.com/calebamiles/keps/pkg/keps/metadata/metadatafakes"
)

const (
	basicInfoTitle     = "The Kubernetes Enhancement Proposal Process"
	basicInfoOwningSIG = "sig-architecture"
)

var basicInfoAuthors = []string{"calebmiles", "jbeda"}

func newBasicRenderingInfo() *metadatafakes.FakeKEP {
	info := &metadatafakes.FakeKEP{}

	info.TitleReturns(basicInfoTitle)
	info.AuthorsReturns(basicInfoAuthors)
	info.OwningSIGReturns(basicInfoOwningSIG)

	return info
}
