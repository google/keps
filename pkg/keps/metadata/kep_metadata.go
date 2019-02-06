package metadata

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

	"github.com/calebamiles/keps/pkg/keps/sections"
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
	SectionLocations() []string // really are section paths

	// should be (string) references to other KEPs
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

	// Mutators (locking)
	SetState(states.Name)
	AddSectionLocations([]string)
	AddApprovers([]string)
	AddReviewers([]string)
	Persist() error

	// External locking support
	sync.Locker
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
		inApproversSet:           make(map[string]bool),
		inReviewersSet:           make(map[string]bool),
		inSectionLocationsSet:    make(map[string]bool),
		RWMutex:                  new(sync.RWMutex),
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
	SectionLocationsField  []string    `yaml:"sections,omitempty"`

	OwningSIGField           string   `yaml:"owning_sig,omitempty"`
	AffectedSubprojectsField []string `yaml:"affected_subprojects,omitempty"`
	ParticipatingSIGsField   []string `yaml:"participating_sigs,omitempty"`
	KubernetesWideField      bool     `yaml:"kubernetes_wide,omitempty"`
	SIGWideField             bool     `yaml:"sig_wide,omitempty"`

	inApproversSet        map[string]bool `yaml:"-"` // do not persist this
	inReviewersSet        map[string]bool `yaml:"-"` // do not persist this
	inSectionLocationsSet map[string]bool `yaml:"-"` // do not persist this
	contentDir            string          `yaml:"-"` // do not persist this

	*sync.RWMutex `yaml:"-"` // do not persist this
}

func (k *kep) Persist() error {
	k.Lock()
	defer k.Unlock()

	// TODO make this a no op if no KEP content has changed

	k.LastUpdatedField = time.Now().UTC()

	k.SectionLocationsField = []string{}
	for p := range k.inSectionLocationsSet {
		k.SectionLocationsField = append(k.SectionLocationsField, p)

	}

	k.ApproversField = []string{}
	for approver := range k.inApproversSet {
		k.ApproversField = append(k.ApproversField, approver)

	}

	k.ReviewersField = []string{}
	for reviewer := range k.inReviewersSet {
		k.ReviewersField = append(k.ReviewersField, reviewer)

	}

	// TODO write tests for this
	sort.Sort(sections.ByOrder(k.SectionLocationsField))
	sort.Strings(k.ApproversField)
	sort.Strings(k.ReviewersField)

	// TODO ensure all section locations exist

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

// sections

func (k *kep) AddSectionLocations(locs []string) {
	k.Lock()
	defer k.Unlock()

	for _, loc := range locs {
		k.inSectionLocationsSet[loc] = true
	}
}

func (k *kep) SectionLocations() []string {
	k.RLock()
	defer k.RUnlock()

	locs := []string{}
	for loc := range k.inSectionLocationsSet {
		locs = append(locs, loc)
	}

	return locs
}

// state

func (k *kep) SetState(state states.Name) {
	k.Lock()
	defer k.Unlock()

	k.StateField = state
}

func (k *kep) State() states.Name {
	k.RLock()
	defer k.RUnlock()

	return k.StateField
}

// owners

func (k *kep) AddApprovers(approvers []string) {
	k.Lock()
	defer k.Unlock()

	for _, approver := range approvers {
		k.inApproversSet[approver] = true
	}
}

func (k *kep) AddReviewers(reviewers []string) {
	k.Lock()
	defer k.Unlock()

	for _, reviewer := range reviewers {
		k.inReviewersSet[reviewer] = true
	}
}

func (k *kep) Authors() []string {
	k.RLock()
	defer k.RUnlock()

	return k.AuthorsField
}

func (k *kep) Reviewers() []string {
	k.RLock()
	defer k.RUnlock()

	reviewers := []string{}
	for reviewer := range k.inReviewersSet {
		reviewers = append(reviewers, reviewer)
	}

	return reviewers
}

func (k *kep) Approvers() []string {
	k.RLock()
	defer k.RUnlock()

	approvers := []string{}
	for approver := range k.inApproversSet {
		approvers = append(approvers, approver)
	}

	return approvers
}

func (k *kep) Editors() []string {
	k.RLock()
	defer k.RUnlock()

	return k.EditorsField
}

// basic metadata

func (k *kep) Title() string {
	k.RLock()
	defer k.RUnlock()

	return k.TitleField
}

func (k *kep) ShortID() int {
	k.RLock()
	defer k.RUnlock()

	if k.ShortIDField != nil {
		return *k.ShortIDField
	}

	return UnsetShortID
}

func (k *kep) UniqueID() string {
	k.RLock()
	defer k.RUnlock()

	return k.UniqueIDField
}

func (k *kep) ContentDir() string {
	k.RLock()
	defer k.RUnlock()

	return k.contentDir
}

func (k *kep) LastUpdated() time.Time {
	k.RLock()
	defer k.RUnlock()

	return k.LastUpdatedField
}

func (k *kep) Created() time.Time {
	k.RLock()
	defer k.RUnlock()

	return k.CreatedField
}

// other KEP references

func (k *kep) Replaces() []string {
	k.RLock()
	defer k.RUnlock()

	return k.ReplacesField
}

func (k *kep) SupersededBy() []string {
	k.RLock()
	defer k.RUnlock()

	return k.SupersededByField
}

// development themes (SIG PM)

func (k *kep) DevelopmentThemes() []string {
	k.RLock()
	defer k.RUnlock()

	return k.DevelopmentThemesField
}

// SIG info

func (k *kep) OwningSIG() string {
	k.RLock()
	defer k.RUnlock()

	return k.OwningSIGField
}

func (k *kep) AffectedSubprojects() []string {
	k.RLock()
	defer k.RUnlock()

	return k.AffectedSubprojectsField
}

func (k *kep) ParticipatingSIGs() []string {
	k.RLock()
	defer k.RUnlock()

	return k.ParticipatingSIGsField
}

func (k *kep) KubernetesWide() bool {
	k.RLock()
	defer k.RUnlock()

	return k.KubernetesWideField
}

func (k *kep) SIGWide() bool {
	k.RLock()
	defer k.RUnlock()

	return k.SIGWideField
}

const (
	metadataFilename = "metadata.yaml"
	UnsetShortID     = -1
)

func sectionKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func fromBytes(b []byte) (*kep, error) {
	k := &kep{
		inSectionLocationsSet: make(map[string]bool),
		inApproversSet:        make(map[string]bool),
		inReviewersSet:        make(map[string]bool),
		RWMutex:               new(sync.RWMutex),
	}

	err := yaml.Unmarshal(b, k)
	if err != nil {
		return nil, err
	}

	for _, p := range k.SectionLocationsField {
		k.inSectionLocationsSet[p] = true
	}

	for _, a := range k.ApproversField {
		k.inApproversSet[a] = true
	}

	for _, r := range k.ReviewersField {
		k.inReviewersSet[r] = true
	}

	return k, nil
}
