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

var _ = Describe("Propose", func() {
	const (
		authorOne         = "handleOne"
		title             = "A Great but Complicated Idea"
		kubernetesWideDir = "kubernetes-wide"
		metadataFilename  = "metadata.yaml"
	)

	It("prepares the KEP for acceptance|deferment|rejection", func() {
		tmpDir, err := ioutil.TempDir("", "kep-propose")
		Expect(err).ToNot(HaveOccurred())
		defer os.RemoveAll(tmpDir)

		kepDirName := "a-good-but-complicated-idea"
		targetDir := filepath.Join(tmpDir, kepDirName)

		runtimeSettings := &settingsfakes.FakeRuntime{}
		runtimeSettings.PrincipalReturns(authorOne)
		runtimeSettings.TargetDirReturns(targetDir)
		runtimeSettings.ContentRootReturns(tmpDir)

		err = workflow.Init(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		// simulate targeting the newly created KEP
		targetDir = filepath.Join(tmpDir, kubernetesWideDir, kepDirName)
		runtimeSettings.TargetDirReturns(targetDir)

		err = workflow.Propose(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		By("updating the KEP state and persisting the KEP")
		kepMetaBytes, err := ioutil.ReadFile(filepath.Join(targetDir, metadataFilename))
		Expect(err).ToNot(HaveOccurred())

		kepMeta, err := metadata.FromBytes(kepMetaBytes)
		Expect(err).ToNot(HaveOccurred())

		Expect(kepMeta.State()).To(Equal(states.Draft))
	})
})
