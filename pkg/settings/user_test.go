package settings_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"

	"github.com/calebamiles/keps/pkg/settings"
)

var _ = Describe("Reading Cached User Settings", func() {
	Describe("ContentRoot()", func() {
		It("returns the saved KEP content root location", func() {
			prevHome := os.Getenv("HOME")

			tmpDir, err := ioutil.TempDir("", "kep-settings")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			settingsDir := filepath.Join(tmpDir, settings.Dirname)
			settingsLocation := filepath.Join(settingsDir, settings.Filename)
			contentLocation := filepath.Join(tmpDir, "content")

			err = os.MkdirAll(settingsDir, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			us := &settings.User{
				ContentRoot: contentLocation,
			}

			settingsBytes, err := yaml.Marshal(us)
			Expect(err).ToNot(HaveOccurred())

			err = ioutil.WriteFile(settingsLocation, settingsBytes, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			err = os.Setenv("HOME", tmpDir)
			Expect(err).ToNot(HaveOccurred())

			err = os.Setenv("XDG_CACHE_HOME", tmpDir)
			Expect(err).ToNot(HaveOccurred())

			defer os.Setenv("HOME", prevHome)
			defer os.Unsetenv("XDG_CACHE_HOME")

			location, err := settings.ContentRoot()
			Expect(err).ToNot(HaveOccurred())
			Expect(location).To(Equal(contentLocation))
		})
	})

	Describe("SaveContentRoot()", func() {
		It("saves the KEP content root location", func() {
			prevHome := os.Getenv("HOME")

			tmpDir, err := ioutil.TempDir("", "kep-settings")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			settingsDir := filepath.Join(tmpDir, settings.Dirname)
			contentLocation := filepath.Join(tmpDir, "content")

			err = os.MkdirAll(settingsDir, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			err = os.Setenv("HOME", tmpDir)
			Expect(err).ToNot(HaveOccurred())

			err = os.Setenv("XDG_CACHE_HOME", tmpDir)
			Expect(err).ToNot(HaveOccurred())

			defer os.Setenv("HOME", prevHome)
			defer os.Unsetenv("XDG_CACHE_HOME")

			err = settings.SaveContentRoot(contentLocation)
			Expect(err).ToNot(HaveOccurred())

			savedRoot, err := settings.ContentRoot()
			Expect(err).ToNot(HaveOccurred())

			Expect(savedRoot).To(Equal(contentLocation))
		})
	})
})
