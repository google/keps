package workflow_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/keps/states"
	"github.com/calebamiles/keps/pkg/settings/settingsfakes"

	"github.com/calebamiles/keps/pkg/workflow"
)

var _ = Describe("Accept", func() {
	const (
		approverOne       = "handleOne"
		title             = "A Great but Complicated Idea"
		kubernetesWideDir = "kubernetes-wide"
		metadataFilename  = "metadata.yaml"
	)

	It("ensures the the KEP is ready for planning", func() {
		tmpDir, err := ioutil.TempDir("", "kep-accept")
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

		kepMetaBytes, err := ioutil.ReadFile(filepath.Join(targetDir, metadataFilename))
		Expect(err).ToNot(HaveOccurred())

		kepMeta, err := metadata.FromBytes(kepMetaBytes)
		Expect(err).ToNot(HaveOccurred())

		By("adding the principal as an approver")
		Expect(kepMeta.Approvers()).To(ContainElement(approverOne))

		By("marking the KEP as provisional")
		Expect(kepMeta.State()).To(Equal(states.Provisional))
	})
})
