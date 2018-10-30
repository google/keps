package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/google/uuid"

	"github.com/calebamiles/keps/pkg/keps/states"
)

type KEP interface {
	UniqueID() string
	Authors() []string
	Title() string
	ShortID() int
	Reviewers() []string
	Approvers() []string
	Editors() []string
	State() states.Name
	Replaces() string
	DevelopmentThemes() []string
	Created() time.Time
	LastUpdated() time.Time

	// Flattened routing info
	OwningSIG() string
	AffectedSubprojects() []string
	ParticipatingSIGs() []string
	KubernetesWide() bool
	SIGWide() bool
	ContentDir() string

	AddSections([]string)
	Sections() []string

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
		CreatedField:             time.Now().UTC(),
		LastUpdatedField:         time.Now().UTC(),
		UniqueIDField:            uuid.New().String(), // note: will panic on error
		StateField:               states.Provisional,
		contentDir:               routingInfo.ContentDir(),
		sectionNames:             make(map[string]bool),
	}

	return k, nil
}

func Open(p string) (KEP, error) {
	kepBytes, err := ioutil.ReadFile(filepath.Join(p, metadataFilename))
	if err != nil {
		return nil, err
	}

	k, err := fromBytes(kepBytes)
	if err != nil {
		return nil, err
	}

	k.contentDir = p

	return k, nil
}

func FromBytes(b []byte) (KEP, error) {
	return fromBytes(b)
}

func fromBytes(b []byte) (*kep, error) {
	k := &kep{}
	err := yaml.Unmarshal(b, k)
	if err != nil {
		return nil, err
	}

	k.sectionNames = make(map[string]bool)

	for _, s := range k.SectionsField {
		switch k.sectionNames[sectionKey(s)] {
		case false:
			k.sectionNames[sectionKey(s)] = true

			k.uniqueSections = append(k.uniqueSections, s)
		default:
			// do nothing if name exists
		}
	}

	return k, nil
}

type kepSection struct {
	FilenameField string `yaml:"filename"`
	NameField     string `yaml:"name"`
}

func (s *kepSection) Name() string     { return s.NameField }
func (s *kepSection) Filename() string { return s.FilenameField }

type kep struct {
	AuthorsField           []string    `yaml:"authors"`
	TitleField             string      `yaml:"title"`
	ShortIDField            *int        `yaml:"kep_number",omitempty`
	ReviewersField         []string    `yaml:"reviewers"`
	ApproversField         []string    `yaml:"approvers"`
	EditorsField           []string    `yaml:"editors"`
	StateField             states.Name `yaml:"state"`
	ReplacesField          string      `yaml:"replaces"`
	DevelopmentThemesField []string    `yaml:"development_themes"`
	LastUpdatedField       time.Time   `yaml:"last_updated"`
	CreatedField           time.Time   `yaml:"created"`
	UniqueIDField          string      `yaml:"uuid"`
	SectionsField          []string    `yaml:"sections"`

	OwningSIGField           string   `yaml:"owning_sig"`
	AffectedSubprojectsField []string `yaml:"affected_subprojects"`
	ParticipatingSIGsField   []string `yaml:"participating_sigs"`
	KubernetesWideField      bool     `yaml:"kubernetes_wide"`
	SIGWideField             bool     `yaml:"sig_wide"`

	sectionNames   map[string]bool `yaml:"-"` // do not persist this
	uniqueSections []string        `yaml:"-"` // do not persist this
	contentDir     string          `yaml:"-"` // do not persist this
	lock           sync.RWMutex    `yaml:"-"` // do not persist this
}

func (k *kep) Persist() error {
	k.lock.Lock()
	defer k.lock.Unlock()

	k.LastUpdatedField = time.Now().UTC()
	k.SectionsField = k.uniqueSections

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

func (k *kep) AddSections(paths []string) {
	k.lock.Lock()
	defer k.lock.Unlock()

	for _, s := range paths {
		switch k.sectionNames[s] {
		case false:
			k.sectionNames[sectionKey(s)] = true

			k.uniqueSections = append(k.uniqueSections, s)
		default:
			// silently do nothing if name exists
		}
	}
}

func (k *kep) Sections() []string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.uniqueSections
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

func (k *kep) ShortID() int {
	k.lock.RLock()
	defer k.lock.RUnlock()

	if k.ShortIDField != nil {
		return *k.ShortIDField
	}

	return UnsetShortID
}

func (k *kep) UniqueID() string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.UniqueIDField
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

func (k *kep) Created() time.Time {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.CreatedField
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

func (k *kep) ParticipatingSIGs() []string {
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
	metadataFilename = "metadata.yaml"
	UnsetShortID   = -1
)

func sectionKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
