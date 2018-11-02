package skeleton_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/keps/skeleton"
)

var _ = Describe("A KEP Skeleton (an individual KEP directory)", func() {
	Describe("initializing a skeleton", func() {
		It("creates the sub directories of a KEP", func() {
			tmpDir, err := ioutil.TempDir("", "kep-skeleton")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			dirProvider := &testDirProvider{contentDir: tmpDir}

			err = skeleton.Init(dirProvider)
			Expect(err).ToNot(HaveOccurred())

			Expect(filepath.Join(tmpDir, "guides")).To(BeADirectory())
			Expect(filepath.Join(tmpDir, "guides", ".gitkeep")).To(BeARegularFile())

			Expect(filepath.Join(tmpDir, "experience_reports")).To(BeADirectory())
			Expect(filepath.Join(tmpDir, "experience_reports", ".gitkeep")).To(BeARegularFile())

			Expect(filepath.Join(tmpDir, "assets")).To(BeADirectory())
			Expect(filepath.Join(tmpDir, "assets", ".gitkeep")).To(BeARegularFile())
		})
	})
})

type testDirProvider struct {
	contentDir string
}

func (p *testDirProvider) ContentDir() string { return p.contentDir }
