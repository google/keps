package settings_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/calebamiles/keps/pkg/settings"
)

var _ = Describe("determining the KEP principal user", func() {
	cacheEnv := "XDG_CACHE_HOME"


	Describe("FindPrincipal()", func() {
		It("looks for the principal's GitHub handle from the environment", func() {
			githubHandle := "@smartcoder"

			tmpDir, err := ioutil.TempDir("", "kep-settings-test")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			// make sure nobody has set the env variable
			Expect(os.Getenv(settings.PrincipalEnv)).To(BeEmpty())
			Expect(os.Getenv(cacheEnv)).To(BeEmpty())

			// make sure that we don't find a user settings file
			err = os.Setenv(cacheEnv, tmpDir)
			Expect(err).ToNot(HaveOccurred())

			By("reading the GitHub handle from the environment")
			err = os.Setenv(settings.PrincipalEnv, githubHandle)
			Expect(err).ToNot(HaveOccurred())
			defer os.Unsetenv(settings.PrincipalEnv)

			foundPrincipal, err := settings.FindPrincipal()
			Expect(err).ToNot(HaveOccurred())

			Expect(foundPrincipal).To(Equal(githubHandle))
		})

		It("looks for the principal's GitHub handle from the settings file", func() {
			githubHandle := "@smartcoder"

			userCacheDir, err := ioutil.TempDir("", "kep-settings")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(userCacheDir)

			err = os.MkdirAll(filepath.Join(userCacheDir, settings.Dirname), os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			err = os.Setenv("XDG_CACHE_HOME", userCacheDir)
			defer os.Unsetenv(cacheEnv)

			var testUserSettings struct {
				GitHubHandle string `yaml:"github_handle"`
			}

			testUserSettings.GitHubHandle = githubHandle

			userSettingsBytes, err := yaml.Marshal(testUserSettings)
			Expect(err).ToNot(HaveOccurred())

			err = ioutil.WriteFile(filepath.Join(userCacheDir, settings.Dirname, settings.Filename), userSettingsBytes, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			foundPrincipal, err := settings.FindPrincipal()
			Expect(err).ToNot(HaveOccurred())

			Expect(foundPrincipal).To(Equal(githubHandle))
		})
	})
})
