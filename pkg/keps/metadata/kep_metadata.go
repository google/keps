package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-yaml/yaml"

	"github.com/calebamiles/keps/pkg/keps/states"
)

type KEP interface {
	Authors() []string
	Title() string
	Number() int
	Reviewers() []string
	Approvers() []string
	Editors() []string
	State() states.Name
	Replaces() string
	DevelopmentThemes() []string
	LastUpdated() time.Time

	// Flattened routing info
	OwningSIG() string
	AffectedSubprojects() []string
	ParticpiatingSIGs() []string
	KubernetesWide() bool
	SIGWide() bool
	ContentDir() string

	Persist() error
}

type routingInfoProvider interface {
	OwningSIG() string
	AffectedSubprojects() []string
	ParticipatingSIGs() []string
	KubernetesWide() bool
	SIGWide() bool
	ContentDir() string
}

func New(authors []string, title string, routingInfo routingInfoProvider) (KEP, error) {
	k := &kep{
		AuthorsField:             authors,
		TitleField:               title,
		OwningSIGField:           routingInfo.OwningSIG(),
		AffectedSubprojectsField: routingInfo.AffectedSubprojects(),
		ParticipatingSIGsField:   routingInfo.ParticipatingSIGs(),
		KubernetesWideField:      routingInfo.KubernetesWide(),
		SIGWideField:             routingInfo.SIGWide(),
		LastUpdatedField:         time.Now(),
		StateField:               states.Provisional,
		contentDir:               routingInfo.ContentDir(),
	}

	return k, nil
}

func Open(p string) (KEP, error) {
	kepBytes, err := ioutil.ReadFile(filepath.Join(p, metadataFilename))
	if err != nil {
		return nil, err
	}

	k := &kep{}
	err = yaml.Unmarshal(kepBytes, k)
	if err != nil {
		return nil, err
	}

	k.contentDir = p

	return k, nil
}

type kep struct {
	AuthorsField           []string    `yaml:"authors"`
	TitleField             string      `yaml"title"`
	NumberField            int         `yaml:"kep_number"`
	ReviewersField         []string    `yaml:"reviewers"`
	ApproversField         []string    `yaml:"approvers"`
	EditorsField           []string    `yaml:"editors"`
	StateField             states.Name `yaml:"state"`
	ReplacesField          string      `yaml:"replaces"`
	DevelopmentThemesField []string    `yaml:"development_themes"`
	LastUpdatedField       time.Time   `yaml:"last_updated"`

	OwningSIGField           string   `yaml:"owning_sig"`
	AffectedSubprojectsField []string `yaml:"affected_subprojects"`
	ParticipatingSIGsField   []string `yaml:"participating_sigs"`
	KubernetesWideField      bool     `yaml:"kubernetes_wide"`
	SIGWideField             bool     `yaml:"sig_wide"`

	contentDir string       `yaml:"-"` // do not persist this
	lock       sync.RWMutex `yaml:"-"` // do not persist this
}

func (k *kep) Persist() error {
	k.lock.Lock()
	defer k.lock.Unlock()

	loc := k.contentDir
	filename := filepath.Join(loc, metadataFilename)

	metaBytes, err := yaml.Marshal(k)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, metaBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (k *kep) Authors() []string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.AuthorsField
}

func (k *kep) Title() string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.TitleField
}

func (k *kep) Number() int {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.NumberField
}

func (k *kep) Reviewers() []string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.ReviewersField
}

func (k *kep) Approvers() []string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.ApproversField
}

func (k *kep) Editors() []string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.EditorsField
}

func (k *kep) State() states.Name {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.StateField
}

func (k *kep) Replaces() string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.ReplacesField
}

func (k *kep) DevelopmentThemes() []string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.DevelopmentThemesField
}

func (k *kep) LastUpdated() time.Time {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.LastUpdatedField
}

func (k *kep) OwningSIG() string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.OwningSIGField
}

func (k *kep) AffectedSubprojects() []string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.AffectedSubprojectsField
}

func (k *kep) ParticpiatingSIGs() []string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.ParticipatingSIGsField
}

func (k *kep) KubernetesWide() bool {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.KubernetesWideField
}

func (k *kep) SIGWide() bool {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.SIGWideField
}

func (k *kep) ContentDir() string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.contentDir
}

const (
	metadataFilename = "metdata.yaml"
)
