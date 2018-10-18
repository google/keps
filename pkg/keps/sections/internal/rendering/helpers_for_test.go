package rendering_test

const (
	basicInfoTitle     = "The Kubernetes Enhancement Proposal Process"
	basicInfoOwningSIG = "sig-architecture"
)

var basicInfoAuthors = []string{"calebmiles", "jbeda"}

func newBasicRenderingInfo() *mockInfoProvider {
	info := newMockInfoProvider()
	info.TitleOutput.Ret0 <- basicInfoTitle
	info.AuthorsOutput.Ret0 <- basicInfoAuthors
	info.OwningSIGOutput.Ret0 <- basicInfoOwningSIG

	return info
}
