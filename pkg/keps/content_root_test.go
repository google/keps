package keps_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/settings"
	"github.com/calebamiles/keps/pkg/sigs"
)

var _ = Describe("finding KEP content", func() {
	Describe("FindContentRoot()", func() {
		It("looks for the content root from the environment", func() {
			tempDir, err := ioutil.TempDir("", "kep-content")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tempDir)

			contentDir := filepath.Join(tempDir, "content")
			err = createSIGDirs(contentDir)
			Expect(err).ToNot(HaveOccurred())

			err = os.Setenv(keps.ContentRootEnv, contentDir)
			Expect(err).ToNot(HaveOccurred())
			defer os.Unsetenv(keps.ContentRootEnv)

			oldContentRoot, err := settings.ContentRoot()
			Expect(err).ToNot(HaveOccurred())

			err = settings.SaveContentRoot("")
			Expect(err).ToNot(HaveOccurred())
			if oldContentRoot != contentDir {
				defer settings.SaveContentRoot(oldContentRoot)
			}

			foundRoot, err := keps.FindContentRoot()
			Expect(err).ToNot(HaveOccurred())

			Expect(foundRoot).To(Equal(contentDir))
		})

		It("looks for the content root in user settings", func() {
			tempDir, err := ioutil.TempDir("", "kep-content")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tempDir)

			contentDir := filepath.Join(tempDir, "content")
			err = createSIGDirs(contentDir)
			Expect(err).ToNot(HaveOccurred())

			err = os.Unsetenv(keps.ContentRootEnv)
			Expect(err).ToNot(HaveOccurred())

			oldContentRoot, err := settings.ContentRoot()
			Expect(err).ToNot(HaveOccurred())

			err = settings.SaveContentRoot(contentDir)
			Expect(err).ToNot(HaveOccurred())
			if oldContentRoot != contentDir {
				defer settings.SaveContentRoot(oldContentRoot)
			}

			foundRoot, err := keps.FindContentRoot()
			Expect(err).ToNot(HaveOccurred())

			Expect(foundRoot).To(Equal(contentDir))
		})

		It("looks for the content root from the intersection of $PWD and $HOME", func() {
			u, err := user.Current()
			Expect(err).ToNot(HaveOccurred())
			homeDir := u.HomeDir
			workspace := filepath.Join(homeDir, "workspace")

			tempDir, err := ioutil.TempDir(workspace, "kep-content")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tempDir)

			contentDir := filepath.Join(tempDir, "content")
			err = createSIGDirs(contentDir)
			Expect(err).ToNot(HaveOccurred())

			err = os.Unsetenv(keps.ContentRootEnv)
			Expect(err).ToNot(HaveOccurred())

			oldContentRoot, err := settings.ContentRoot()
			Expect(err).ToNot(HaveOccurred())

			err = settings.SaveContentRoot("")
			Expect(err).ToNot(HaveOccurred())
			if oldContentRoot != "" {
				defer settings.SaveContentRoot(oldContentRoot)
			}

			foundRoot, err := keps.FindContentRoot()
			Expect(err).ToNot(HaveOccurred())
			Expect(foundRoot).To(Equal(contentDir))
		})

		It("returns an error when attempting to search for content not under $HOME", func() {
			tempDir, err := ioutil.TempDir("", "kep-content")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tempDir)

			err = os.Unsetenv(keps.ContentRootEnv)
			Expect(err).ToNot(HaveOccurred())

			oldContentRoot, err := settings.ContentRoot()
			Expect(err).ToNot(HaveOccurred())

			err = settings.SaveContentRoot("")
			Expect(err).ToNot(HaveOccurred())
			if oldContentRoot != "" {
				defer settings.SaveContentRoot(oldContentRoot)
			}

			err = os.Chdir(tempDir)
			Expect(err).ToNot(HaveOccurred())
			_, err = keps.FindContentRoot()
			Expect(err.Error()).To(ContainSubstring("file search must start at location under $HOME"))
		})

		It("returns an error if no content was found", func() {
			u, err := user.Current()
			Expect(err).ToNot(HaveOccurred())
			homeDir := u.HomeDir
			workspace := filepath.Join(homeDir, "workspace")

			tempDir, err := ioutil.TempDir(workspace, "kep-content")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tempDir)

			err = os.Unsetenv(keps.ContentRootEnv)
			Expect(err).ToNot(HaveOccurred())

			oldContentRoot, err := settings.ContentRoot()
			Expect(err).ToNot(HaveOccurred())

			err = settings.SaveContentRoot("")
			Expect(err).ToNot(HaveOccurred())
			if oldContentRoot != "" {
				defer settings.SaveContentRoot(oldContentRoot)
			}

			err = os.Chdir(tempDir)
			Expect(err).ToNot(HaveOccurred())

			_, err = keps.FindContentRoot()
			Expect(err.Error()).To(ContainSubstring("could not find KEP content"))

		})

	})
})

func createSIGDirs(loc string) error {
	for _, sig := range sigs.All() {
		err := os.MkdirAll(filepath.Join(loc, sig), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
