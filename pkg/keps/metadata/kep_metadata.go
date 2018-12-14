package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

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
	DevelopmentThemes() []string
	Sections() []string // really are section paths

	// should be pointers to other KEPs
	Replaces() []string
	SupersededBy() []string

	Created() time.Time
	LastUpdated() time.Time

	// Flattened routing info
	OwningSIG() string
	AffectedSubprojects() []string
	ParticipatingSIGs() []string
	KubernetesWide() bool
	SIGWide() bool
	ContentDir() string

	// Mutators
	SetState(states.Name)
	AddSections([]string)
	AddApprovers([]string)
	AddReviewers([]string)
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
		StateField:               states.Draft,
		contentDir:               routingInfo.ContentDir(),
		hasSectionPath:           make(map[string]bool),
		inApproversSet:           make(map[string]bool),
		inReviewersSet:           make(map[string]bool),
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
	k := &kep{
		hasSectionPath: make(map[string]bool),
		inApproversSet: make(map[string]bool),
		inReviewersSet: make(map[string]bool),
	}

	err := yaml.Unmarshal(b, k)
	if err != nil {
		return nil, err
	}

	k.hasSectionPath = make(map[string]bool)

	for _, s := range k.SectionsField {
		if !k.hasSectionPath[sectionKey(s)] {
			k.hasSectionPath[sectionKey(s)] = true
			k.uniqueSections = append(k.uniqueSections, s)
		}
	}

	for _, a := range k.ApproversField {
		k.inApproversSet[a] = true
	}

	for _, r := range k.ReviewersField {
		k.inReviewersSet[r] = true
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
	AuthorsField           []string    `yaml:"authors,omitempty"`
	TitleField             string      `yaml:"title,omitempty"`
	ShortIDField           *int        `yaml:"kep_number,omitempty"`
	ReviewersField         []string    `yaml:"reviewers,omitempty"`
	ApproversField         []string    `yaml:"approvers,omitempty"`
	EditorsField           []string    `yaml:"editors,omitempty"`
	StateField             states.Name `yaml:"state,omitempty"`
	ReplacesField          []string    `yaml:"replaces,omitempty"`
	SupersededByField      []string    `yaml:"superseded_by,omitempty"`
	DevelopmentThemesField []string    `yaml:"development_themes,omitempty"`
	LastUpdatedField       time.Time   `yaml:"last_updated,omitempty"`
	CreatedField           time.Time   `yaml:"created,omitempty"`
	UniqueIDField          string      `yaml:"uuid,omitempty"`
	SectionsField          []string    `yaml:"sections,omitempty"`

	OwningSIGField           string   `yaml:"owning_sig,omitempty"`
	AffectedSubprojectsField []string `yaml:"affected_subprojects,omitempty"`
	ParticipatingSIGsField   []string `yaml:"participating_sigs,omitempty"`
	KubernetesWideField      bool     `yaml:"kubernetes_wide,omitempty"`
	SIGWideField             bool     `yaml:"sig_wide,omitempty"`

	hasSectionPath map[string]bool `yaml:"-"` // do not persist this
	uniqueSections []string        `yaml:"-"` // do not persist this
	inApproversSet map[string]bool `yaml:"-"` // do not persist this
	inReviewersSet map[string]bool `yaml:"-"` // do not persist this
	contentDir     string          `yaml:"-"` // do not persist this
	lock           sync.RWMutex    `yaml:"-"` // do not persist this
}

func (k *kep) Persist() error {
	k.lock.Lock()
	defer k.lock.Unlock()

	k.LastUpdatedField = time.Now().UTC()
	k.SectionsField = k.uniqueSections

	k.ApproversField = []string{}
	for approver := range k.inApproversSet {
		k.ApproversField = append(k.ApproversField, approver)
	}

	k.ReviewersField = []string{}
	for reviewer := range k.inReviewersSet {
		k.ReviewersField = append(k.ReviewersField, reviewer)
	}

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
		if !k.hasSectionPath[s] {
			k.hasSectionPath[sectionKey(s)] = true
			// we maintain a separate slice in order to maintain section order
			k.uniqueSections = append(k.uniqueSections, s)
		}
	}
}

func (k *kep) AddApprovers(approvers []string) {
	k.lock.Lock()
	defer k.lock.Unlock()

	// we don't currently care about the order of approvers when persisting this back to YAML
	for _, approver := range approvers {
		k.inApproversSet[approver] = true
	}
}

func (k *kep) AddReviewers(reviewers []string) {
	k.lock.Lock()
	defer k.lock.Unlock()

	// we don't currently care about the order of reviewers when persisting this back to YAML
	for _, reviewer := range reviewers {
		k.inReviewersSet[reviewer] = true
	}
}

func (k *kep) SetState(state states.Name) {
	k.lock.Lock()
	defer k.lock.Unlock()

	k.StateField = state
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

	reviewers := []string{}
	for reviewer := range k.inReviewersSet {
		reviewers = append(reviewers, reviewer)
	}

	return reviewers
}

func (k *kep) Approvers() []string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	approvers := []string{}
	for approver := range k.inApproversSet {
		approvers = append(approvers, approver)
	}

	return approvers
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

func (k *kep) Replaces() []string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.ReplacesField
}

func (k *kep) SupersededBy() []string {
	k.lock.RLock()
	defer k.lock.RUnlock()

	return k.SupersededByField
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
	UnsetShortID     = -1
)

func sectionKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
