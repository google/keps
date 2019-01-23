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

var _ = Describe("Init", func() {
	const (
		authorOne           = "handleOne"
		title               = "A Great but Complicated Idea"
		kubernetesWideDir   = "kubernetes-wide"
		sigWideDir          = "sig-wide"
		metadataFilename    = "metadata.yaml"
		readmeFilename      = "README.md"
		experienceReportDir = "experience_reports"
		assetsDir           = "assets"
		guidesDir           = "guides"
		summaryFilename     = "summary.md"
		motivationFilename  = "motivation.md"
	)

	Context("when creating a Kubernetes Wide KEP", func() {
		It("creates a KEP in the kubernetes-wide directory", func() {
			tmpDir, err := ioutil.TempDir("", "kep-init")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			kepDirName := "a-good-but-complicated-idea"

			runtimeSettings := &settingsfakes.FakeRuntime{}
			runtimeSettings.PrincipalReturns(authorOne)
			runtimeSettings.TargetDirReturns(kepDirName)
			runtimeSettings.ContentRootReturns(tmpDir)

			expectedKEPContentDir, err := workflow.Init(runtimeSettings)
			Expect(err).ToNot(HaveOccurred())

			Expect(expectedKEPContentDir).To(BeADirectory())
			Expect(filepath.Join(expectedKEPContentDir, experienceReportDir)).To(BeADirectory())
			Expect(filepath.Join(expectedKEPContentDir, assetsDir)).To(BeADirectory())
			Expect(filepath.Join(expectedKEPContentDir, metadataFilename)).To(BeARegularFile())
			Expect(filepath.Join(expectedKEPContentDir, summaryFilename)).To(BeARegularFile())
			Expect(filepath.Join(expectedKEPContentDir, motivationFilename)).To(BeARegularFile())
			Expect(filepath.Join(expectedKEPContentDir, guidesDir)).To(BeADirectory())
		})
	})

	Context("when creating a SIG Wide KEP", func() {
		It("creates a KEP in the sig-wide directory for the inferred SIG", func() {
			tmpDir, err := ioutil.TempDir("", "kep-init")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			kepDirName := "kubelet-v2-api"
			sigDir := "sig-node"

			targetDir := filepath.Join(sigDir, kepDirName)

			runtimeSettings := &settingsfakes.FakeRuntime{}
			runtimeSettings.PrincipalReturns(authorOne)
			runtimeSettings.TargetDirReturns(targetDir)
			runtimeSettings.ContentRootReturns(tmpDir)

			expectedKEPContentDir, err := workflow.Init(runtimeSettings)
			Expect(err).ToNot(HaveOccurred())

			Expect(expectedKEPContentDir).To(BeADirectory())
			Expect(filepath.Join(expectedKEPContentDir, experienceReportDir)).To(BeADirectory())
			Expect(filepath.Join(expectedKEPContentDir, assetsDir)).To(BeADirectory())
			Expect(filepath.Join(expectedKEPContentDir, guidesDir)).To(BeADirectory())
			Expect(filepath.Join(expectedKEPContentDir, metadataFilename)).To(BeARegularFile())
			Expect(filepath.Join(expectedKEPContentDir, summaryFilename)).To(BeARegularFile())
			Expect(filepath.Join(expectedKEPContentDir, motivationFilename)).To(BeARegularFile())
		})
	})

	Context("when creating a SIG specific KEP under an existing subproject", func() {
		It("creates a KEP skeleton at a given path", func() {
			tmpDir, err := ioutil.TempDir("", "kep-init")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			sigDir := "sig-node"
			subprojectDir := "kubelet"
			kepDirName := "dynamic-kubelet-configuration"

			targetDir := filepath.Join(sigDir, subprojectDir, kepDirName)

			runtimeSettings := &settingsfakes.FakeRuntime{}
			runtimeSettings.PrincipalReturns(authorOne)
			runtimeSettings.TargetDirReturns(targetDir)
			runtimeSettings.ContentRootReturns(tmpDir)

			expectedKEPContentDir, err := workflow.Init(runtimeSettings)
			Expect(err).ToNot(HaveOccurred())

			Expect(expectedKEPContentDir).To(BeADirectory())
			Expect(filepath.Join(expectedKEPContentDir, experienceReportDir)).To(BeADirectory())
			Expect(filepath.Join(expectedKEPContentDir, assetsDir)).To(BeADirectory())
			Expect(filepath.Join(expectedKEPContentDir, guidesDir)).To(BeADirectory())
			Expect(filepath.Join(expectedKEPContentDir, metadataFilename)).To(BeARegularFile())
			Expect(filepath.Join(expectedKEPContentDir, summaryFilename)).To(BeARegularFile())
			Expect(filepath.Join(expectedKEPContentDir, motivationFilename)).To(BeARegularFile())
		})
	})
})
