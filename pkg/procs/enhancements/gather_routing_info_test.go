package enhancements_test

import (

	"github.com/calebamiles/keps/pkg/changes/routing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/kepsfakes"
	"github.com/calebamiles/keps/pkg/orgs/orgsfakes"
	"github.com/calebamiles/keps/pkg/settings/settingsfakes"

	"github.com/calebamiles/keps/pkg/procs/enhancements"
)

var _ = Describe("Gathering Information Required to Work With Git and GitHub", func() {
	Describe("GatherRoutingFrom()", func() {
		It("collects information from runtime, kep, and org", func() {
			Fail("test not written")
			var routingInfo routing.Info
			var err error

			fakeOrg := &orgsfakes.FakeInstance{}
			fakeKep := &kepsfakes.FakeInstance{}
			fakeRuntime := &settingsfakes.FakeRuntime{}

			By("returning an implementation of routing.Info")

			routingInfo, err = enhancements.GatherRoutingFrom(fakeRuntime, fakeKep, fakeOrg)
			Expect(err).ToNot(HaveOccurred(), "expected no error gathering routing info from fake but valid settings.Runtime, keps.Instance, orgs.Instance")

			Expect(routingInfo.SourceRepositoryOwner()).To(Equal(fakeRuntime.PrincipalGithubHandle()), "expected source repository owner to match that given by the runtime")
			Expect(routingInfo.SourceRepository()).To(Equal(fakeOrg.EnhancementsRepository()), "expected source repository to match that given by the org")
			Expect(routingInfo.SourceBranch()).To(Equal(fakeOrg.EnhancementsRepositoryDefaultBranch()), "expected source repository to match that given by the org")
			Expect(routingInfo.TargetOwner()).To(Equal(fakeOrg.EnhancementsRepository()), "expected target repository to match that given by the org")
			Expect(routingInfo.TargetRepository()).To(Equal(fakeOrg.EnhancementsRepository()), "expected target repository to match that given by the org")
			Expect(routingInfo.TargetBranch()).To(Equal(fakeOrg.EnhancementsRepositoryDefaultBranch()), "expected target repository to match that given by the org")
		})
	})
})
