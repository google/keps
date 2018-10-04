package settings_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/calebamiles/keps/pkg/settings"
	"github.com/calebamiles/keps/pkg/sigs"
)

var _ = Describe("finding KEP content", func() {
	Describe("FindContentRoot()", func() {
		It("looks for the content root from the environment", func() {
			contentDir, err := createSIGDirs()
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(contentDir)

			err = os.Setenv(settings.ContentRootEnv, contentDir)
			Expect(err).ToNot(HaveOccurred())
			defer os.Unsetenv(settings.ContentRootEnv)

			foundRoot, err := settings.FindContentRoot()
			Expect(err).ToNot(HaveOccurred())

			Expect(foundRoot).To(Equal(contentDir))
		})

		It("looks for the content root in user settings", func() {
			contentDir, err := createSIGDirs()
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(contentDir)

			// make sure nobody has set the env variable
			Expect(os.Getenv(settings.ContentRootEnv)).To(BeEmpty())

			// set user cache dir to someplace unused
			err = os.Setenv(cacheEnv, contentDir)
			Expect(err).ToNot(HaveOccurred())

			err = settings.SaveContentRoot(contentDir)
			Expect(err).ToNot(HaveOccurred())

			foundRoot, err := settings.FindContentRoot()
			Expect(err).ToNot(HaveOccurred())

			Expect(foundRoot).To(Equal(contentDir))

			err = os.Unsetenv(cacheEnv)
			Expect(err).ToNot(HaveOccurred())
		})

		It("looks for the content root from the intersection of $PWD and $HOME", func() {
			u, err := user.Current()
			Expect(err).ToNot(HaveOccurred())

			homeDir := u.HomeDir
			workspace := filepath.Join(homeDir, "workspace")

			contentDir, err := ioutil.TempDir(workspace, "kep-content")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(contentDir)

			err = createSIGDirsAt(contentDir)
			Expect(err).ToNot(HaveOccurred())

			// make sure nobody has set the env variable
			Expect(os.Getenv(settings.ContentRootEnv)).To(BeEmpty())

			// set user cache dir to someplace unused
			err = os.Setenv(cacheEnv, contentDir)
			Expect(err).ToNot(HaveOccurred())
			defer os.Unsetenv(cacheEnv)

			foundRoot, err := settings.FindContentRoot()
			Expect(err).ToNot(HaveOccurred())
			Expect(foundRoot).To(Equal(contentDir))
		})

		It("returns an error when attempting to search for content not under $HOME", func() {
			tempDir, err := ioutil.TempDir("", "kep-content")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tempDir)

			// make sure nobody has set the env variable
			Expect(os.Getenv(settings.ContentRootEnv)).To(BeEmpty())

			// set user cache dir to someplace unused
			err = os.Setenv(cacheEnv, tempDir)
			Expect(err).ToNot(HaveOccurred())
			defer os.Unsetenv(cacheEnv)

			pwd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			err = os.Chdir(tempDir)
			Expect(err).ToNot(HaveOccurred())
			defer os.Chdir(pwd)

			_, err = settings.FindContentRoot()
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

			err = os.Setenv(cacheEnv, tempDir)
			Expect(err).ToNot(HaveOccurred())
			defer os.Unsetenv(cacheEnv)

			pwd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())

			err = os.Chdir(tempDir)
			Expect(err).ToNot(HaveOccurred())
			defer os.Chdir(pwd)

			_, err = settings.FindContentRoot()
			Expect(err.Error()).To(ContainSubstring("could not find KEP content"))
		})
	})

	Describe("SaveContentRoot()", func() {
		It("writes the content root location to disk", func() {
			tempDir, err := ioutil.TempDir("", "kep-content")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tempDir)

			// make sure nobody has set the env variable
			Expect(os.Getenv(settings.ContentRootEnv)).To(BeEmpty())

			// set user cache dir to someplace unused
			err = os.Setenv(cacheEnv, tempDir)
			Expect(err).ToNot(HaveOccurred())

			err = createSIGDirsAt(tempDir)
			Expect(err).ToNot(HaveOccurred())

			err = settings.SaveContentRoot(tempDir)
			Expect(err).ToNot(HaveOccurred())

			savedPath, err := settings.FindContentRoot()
			Expect(err).ToNot(HaveOccurred())
			Expect(savedPath).To(Equal(tempDir))
		})
	})
})

const (
	cacheEnv = "XDG_CACHE_HOME"
)

func createSIGDirs() (string, error) {
	tempDir, err := ioutil.TempDir("", "kep-content")
	if err != nil {
		return "", err
	}

	err = createSIGDirsAt(tempDir)
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func createSIGDirsAt(p string) error {
	for _, sig := range sigs.All() {
		err := os.MkdirAll(filepath.Join(p, sig), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
