package workflow_test

import (
	"io/ioutil"
	_ "os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps"
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

	FIt("prepares the KEP for acceptance|deferment|rejection", func() {
		tmpDir, err := ioutil.TempDir("", "kep-propose")
		Expect(err).ToNot(HaveOccurred())
		//defer os.RemoveAll(tmpDir)
		println(tmpDir)

		contentRoot := filepath.Join(tmpDir, "content")

		err = createSIGDirsAt(contentRoot)
		Expect(err).ToNot(HaveOccurred(), "creating SIG directories")

		kepDirName := "value-delivered-over-multiple-releases"
		targetDir := kepDirName

		runtimeSettings := &settingsfakes.FakeRuntime{}
		runtimeSettings.PrincipalReturns(authorOne)
		runtimeSettings.TargetDirReturns(targetDir)
		runtimeSettings.ContentRootReturns(contentRoot)

		targetDir, err = workflow.Init(runtimeSettings)
		Expect(err).ToNot(HaveOccurred(), "simulating `kep init`")

		// simulate targeting the newly created KEP
		runtimeSettings.TargetDirReturns(targetDir)

		By("updating the KEP state and persisting the KEP")
		err = workflow.Propose(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		By("marking the KEP as draft")
		kep, err := keps.Open(targetDir)
		Expect(err).ToNot(HaveOccurred(), "opening KEP after propose")

		Expect(kep.State()).To(Equal(states.Draft))
		Fail("this test should fail because README.md is duplicated in metadata.yaml")

		// check that has matching state and last updated
		Fail("no checking of README.md")
	})
})
