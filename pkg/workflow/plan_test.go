package workflow_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/settings/settingsfakes"

	"github.com/calebamiles/keps/pkg/workflow"
)

var _ = Describe("Plan()", func() {
	const (
		approverOne            = "handleOne"
		kubernetesWideDir      = "kubernetes-wide"
		metadataFilename       = "metadata.yaml"
		teacherGuideFilename   = "guides/teacher.md"
		developerGuideFilename = "guides/developer.md"
		operatorGuideFilename  = "guides/operator.md"
	)

	It("ensures the KEP is ready for approval", func() {
		tmpDir, err := ioutil.TempDir("", "kep-plan")
		Expect(err).ToNot(HaveOccurred())
		defer os.RemoveAll(tmpDir)

		kepDirName := "a-good-but-complicated-idea"
		targetDir := filepath.Join(tmpDir, kepDirName)

		runtimeSettings := &settingsfakes.FakeRuntime{}
		runtimeSettings.PrincipalReturns(approverOne)
		runtimeSettings.TargetDirReturns(targetDir)
		runtimeSettings.ContentRootReturns(tmpDir)

		targetDir, err = workflow.Init(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		// simulate targeting the newly created KEP
		runtimeSettings.TargetDirReturns(targetDir)

		err = workflow.Propose(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		err = workflow.Accept(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		err = workflow.Plan(runtimeSettings)
		Expect(err).ToNot(HaveOccurred())

		By("creating templates for content under guides/")
		Expect(filepath.Join(targetDir, teacherGuideFilename)).To(BeARegularFile())
		Expect(filepath.Join(targetDir, developerGuideFilename)).To(BeARegularFile())
		Expect(filepath.Join(targetDir, operatorGuideFilename)).To(BeARegularFile())
	})
})
