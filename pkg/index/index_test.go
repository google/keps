package index_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"

	"github.com/calebamiles/keps/pkg/index"
	"github.com/calebamiles/keps/pkg/keps/states"

	"github.com/calebamiles/keps/pkg/keps/kepsfakes"
	"github.com/calebamiles/keps/pkg/settings/settingsfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("working with an index of KEPs", func() {
	Describe("Rebuild()", func() {
		It("builds an index of KEP entries", func() {
			tmpDir, err := ioutil.TempDir("", "kep-index")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			kepDirOne := filepath.Join(tmpDir, "sig-architecture", "kep1")
			kepDirTwo := filepath.Join(tmpDir, "sig-architecture", "kep2")
			kepDirThree := filepath.Join(tmpDir, "sig-architecture", "kep3")

			err = os.MkdirAll(kepDirOne, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			err = os.MkdirAll(kepDirTwo, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			err = os.MkdirAll(kepDirThree, os.ModePerm)
			Expect(err).ToNot(HaveOccurred())

			writeTestMetadata(kepDirOne)
			writeTestMetadata(kepDirTwo)
			writeTestMetadata(kepDirThree)

			fakeSettings := &settingsfakes.FakeRuntime{}
			fakeSettings.ContentRootReturns(tmpDir)

			_, err = index.Rebuild(fakeSettings)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when attempting to add a KEP returns an error", func() {
			It("finishes looking for KEPs before returning an error", func() {
				tmpDir, err := ioutil.TempDir("", "kep-index")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				kepDirOne := filepath.Join(tmpDir, "sig-architecture", "kep1")
				kepDirTwo := filepath.Join(tmpDir, "sig-architecture", "kep2")
				kepDirThree := filepath.Join(tmpDir, "sig-architecture", "kep3")

				err = os.MkdirAll(kepDirOne, os.ModePerm)
				Expect(err).ToNot(HaveOccurred())

				err = os.MkdirAll(kepDirTwo, os.ModePerm)
				Expect(err).ToNot(HaveOccurred())

				err = os.MkdirAll(kepDirThree, os.ModePerm)
				Expect(err).ToNot(HaveOccurred())

				writeTestMetadata(kepDirOne)

				err = ioutil.WriteFile(filepath.Join(kepDirTwo, "metadata.yaml"), []byte("invalid yaml"), os.ModePerm)
				Expect(err).ToNot(HaveOccurred())

				err = ioutil.WriteFile(filepath.Join(kepDirThree, "metadata.yaml"), []byte("invalid yaml"), os.ModePerm)
				Expect(err).ToNot(HaveOccurred())

				fakeSettings := &settingsfakes.FakeRuntime{}
				fakeSettings.ContentRootReturns(tmpDir)

				_, err = index.Rebuild(fakeSettings)
				Expect(err).To(HaveOccurred())

				merr, ok := err.(*multierror.Error)
				Expect(ok).To(BeTrue())

				Expect(merr.Errors).To(HaveLen(2))
			})
		})
	})

	Describe("operations on an Index", func() {
		Describe("#Persist()", func() {
			It("persists the KEP entries to disk", func() {
				tmpDir, err := ioutil.TempDir("", "kep-index")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				k := &kepsfakes.FakeInstance{}
				k.ShortIDReturns(42)
				k.TitleReturns("The Kubernetes Enhancement Proposal Process")
				k.UniqueIDReturns("a-valid-uuid")
				k.AuthorsReturns([]string{"jbeda", "calebamiles"})
				k.ContentDirReturns("a-valid-directory-location")
				k.CreatedReturns(time.Now().Add(-time.Hour))
				k.LastUpdatedReturns(time.Now())
				k.StateReturns(states.Implementable)
				k.CheckReturns(nil)

				kepIndex, err := index.New(tmpDir)
				Expect(err).ToNot(HaveOccurred())

				err = kepIndex.Update(k)
				Expect(err).ToNot(HaveOccurred())

				err = kepIndex.Persist()
				Expect(err).ToNot(HaveOccurred())

				Expect(filepath.Join(tmpDir, "keps.yaml")).To(BeARegularFile())
			})
		})

		Describe("#Update()", func() {
			It("attempts to add a KEP instance to the index", func() {
				tmpDir, err := ioutil.TempDir("", "kep-index")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				k := &kepsfakes.FakeInstance{}
				k.ShortIDReturns(42)
				k.TitleReturns("The Kubernetes Enhancement Proposal Process")
				k.UniqueIDReturns("a-valid-uuid")
				k.AuthorsReturns([]string{"jbeda", "calebamiles"})
				k.ContentDirReturns("a-valid-directory-location")
				k.CreatedReturns(time.Now().Add(-time.Hour))
				k.LastUpdatedReturns(time.Now())
				k.StateReturns(states.Implementable)

				kepIndex, err := index.New(tmpDir)
				Expect(err).ToNot(HaveOccurred())

				By("adding KEPs that self report as valid")
				k.CheckReturns(nil)
				err = kepIndex.Update(k)
				Expect(err).ToNot(HaveOccurred())

				By("adding global index consistency checks to KEPs")
				Expect(k.AddChecksCallCount()).To(Equal(1))
				addedChecks := k.AddChecksArgsForCall(0)
				Expect(addedChecks).To(HaveLen(2))
			})
		})

		Describe("#Fetch", func() {
			It("returns a KEP instance by looking up its unique ID", func() {
				tmpDir, err := ioutil.TempDir("", "kep-index")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				k := &kepsfakes.FakeInstance{}
				k.ShortIDReturns(42)
				k.TitleReturns("The Kubernetes Enhancement Proposal Process")
				k.UniqueIDReturns("a-valid-uuid")
				k.AuthorsReturns([]string{"jbeda", "calebamiles"})
				k.ContentDirReturns("a-valid-directory-location")
				k.CreatedReturns(time.Now().Add(-time.Hour))
				k.LastUpdatedReturns(time.Now())
				k.StateReturns(states.Implementable)

				kepIndex, err := index.New(tmpDir)
				Expect(err).ToNot(HaveOccurred())

				k.CheckReturns(nil)
				err = kepIndex.Update(k)
				Expect(err).ToNot(HaveOccurred())

				By("returning the KEP with the given unique ID")
				_, err = kepIndex.Fetch(k.UniqueID())
				Expect(err).ToNot(HaveOccurred())
			})

			Context("when no KEP with given unique ID exists", func() {
				It("returns an error", func() {
					tmpDir, err := ioutil.TempDir("", "kep-index")
					Expect(err).ToNot(HaveOccurred())
					defer os.RemoveAll(tmpDir)

					kepIndex, err := index.New(tmpDir)
					Expect(err).ToNot(HaveOccurred())

					_, err = kepIndex.Fetch("not-an-id-in-the-index")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("no KEP with unique ID"))
				})
			})
		})

		Describe("Open()", func() {
			It("reads a kep.yaml from disk", func() {
				tmpDir, err := ioutil.TempDir("", "kep-index")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(tmpDir)

				kepIndex, err := index.New(tmpDir)
				Expect(err).ToNot(HaveOccurred())

				err = kepIndex.Persist()
				Expect(err).ToNot(HaveOccurred())

				_, err = index.Open(tmpDir)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})

func writeTestMetadata(dir string) {
	tm := &testMetadata{
		AuthorsField:     []string{"Test Author " + uuid.New().String(), "Test Author " + uuid.New().String()},
		TitleField:       "Test Title " + uuid.New().String(),
		StateField:       states.Implemented,
		LastUpdatedField: time.Now(),
		CreatedField:     time.Now().Add(-time.Hour),
		UniqueIDField:    uuid.New().String(),
		OwningSIGField:   "sig-architecture",
	}

	tmBytes, err := yaml.Marshal(tm)
	Expect(err).ToNot(HaveOccurred())

	err = ioutil.WriteFile(filepath.Join(dir, "metadata.yaml"), tmBytes, os.ModePerm)
	Expect(err).ToNot(HaveOccurred())
}

type testMetadata struct {
	AuthorsField     []string    `yaml:"authors"`
	TitleField       string      `yaml:"title"`
	ShortIDField     *int        `yaml:"kep_number",omitempty`
	StateField       states.Name `yaml:"state"`
	LastUpdatedField time.Time   `yaml:"last_updated"`
	CreatedField     time.Time   `yaml:"created"`
	UniqueIDField    string      `yaml:"uuid"`
	SectionsField    []string    `yaml:"sections"`
	OwningSIGField   string      `yaml:"owning_sig"`
}
