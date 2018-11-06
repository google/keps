package workflow_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/index"
	"github.com/calebamiles/keps/pkg/settings/settingsfakes"
	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/keps/states"

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
		targetDir := filepath.Join(tmpDir, kepDirName)

		runtimeSettings := &settingsfakes.FakeRuntime{}
		runtimeSettings.PrincipalReturns(approverOne)
		runtimeSettings.TargetDirReturns(targetDir)
		runtimeSettings.ContentRootReturns(tmpDir)

		err = workflow.Init(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		// simulate targeting the newly created KEP
		targetDir = filepath.Join(tmpDir, kubernetesWideDir, kepDirName)
		runtimeSettings.TargetDirReturns(targetDir)

		kepMetaBytes, err := ioutil.ReadFile(filepath.Join(targetDir, metadataFilename))
		Expect(err).ToNot(HaveOccurred())

		kepMeta, err := metadata.FromBytes(kepMetaBytes)
		Expect(err).ToNot(HaveOccurred())

		// read the created unique ID
		createdKEPId := kepMeta.UniqueID()

		err = workflow.Propose(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		err = workflow.Accept(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		err = workflow.Plan(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		err = workflow.Approve(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		kepMetaBytes, err = ioutil.ReadFile(filepath.Join(targetDir, metadataFilename))
		Expect(err).ToNot(HaveOccurred())

		kepMeta, err = metadata.FromBytes(kepMetaBytes)
		Expect(err).ToNot(HaveOccurred())

		By("marking the KEP as implementable")
		Expect(kepMeta.State()).To(Equal(states.Implementable))

		By("rebuilding and persisting the KEP index")
		kepIndex, err := index.Open(tmpDir)
		Expect(err).ToNot(HaveOccurred())

		_, err = kepIndex.Fetch(createdKEPId)
		Expect(err).ToNot(HaveOccurred())
	})
})
