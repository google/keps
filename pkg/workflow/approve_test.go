package workflow_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/states"
	"github.com/calebamiles/keps/pkg/settings/settingsfakes"

	"github.com/calebamiles/keps/pkg/workflow"
)

var _ = Describe("Approve()", func() {
	const (
		approverOne       = "handleOne"
		kubernetesWideDir = "kubernetes-wide"
		metadataFilename  = "metadata.yaml"
		indexFilename     = "keps.yaml"
	)

	It("ensures the KEP is ready for implementation", func() {
		tmpDir, err := ioutil.TempDir("", "kep-approve")
		Expect(err).ToNot(HaveOccurred())
		defer os.RemoveAll(tmpDir)

		kepDirName := "a-good-but-complicated-idea"

		runtimeSettings := &settingsfakes.FakeRuntime{}
		runtimeSettings.PrincipalReturns(approverOne)
		runtimeSettings.TargetDirReturns(kepDirName)
		runtimeSettings.ContentRootReturns(tmpDir)

		targetDir, err := workflow.Init(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		// simulate targeting the newly created KEP
		runtimeSettings.TargetDirReturns(targetDir)

		err = workflow.Propose(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		err = workflow.Accept(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		err = workflow.Plan(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		By("updating the KEP state and persisting the KEP")
		err = workflow.Approve(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		By("marking the KEP as implementable")
		kep, err := keps.Open(targetDir)
		Expect(err).ToNot(HaveOccurred(), "opening KEP after approve")

		Expect(kep.State()).To(Equal(states.Implementable))
	})
})
