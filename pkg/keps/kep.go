package keps

import (
	"sync"
	"time"

	"github.com/calebamiles/keps/pkg/keps/check"
	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/keps/sections"
	"github.com/calebamiles/keps/pkg/keps/states"
)

type Instance interface {
	UniqueID() string
	ShortID() int
	Title() string
	OwningSIG() string
	Authors() []string
	ContentDir() string
	State() states.Name
	Created() time.Time
	LastUpdated() time.Time
	Sections() []string

	// simple pass through mutators
	AddApprovers(...string)
	AddReviewers(...string)

	// heavy lifting mutators
	SetState(states.Name) error

	// consistency
	AddChecks(...check.That)
	Check() error

	// flush to disk
	Persist() error
}

// New creates a new Instance from a sections.Collection and a metadata.KEP
func New(meta metadata.KEP, existingEntries []sections.Entry) (Instance, error) {
	k := &kep{
		meta:    meta,
		locker:  new(sync.RWMutex),
		content: make(map[sections.Entry]bool),
	}

	checks := []check.That{check.ThatAllBasicInvariantsAreSatisfied}
	switch meta.State() {
	case states.Provisional:
		checks = append(checks, check.ThatIsValidForProvisionalState)
	case states.Implementable:
		checks = append(checks, check.ThatIsValidForImplementableState)
	}

	checkAll := check.All(checks)
	err := checkAll(meta)
	if err != nil {
		return nil, err
	}

	k.checks = checks
	k.addSections(existingEntries)

	return k, nil
}

// Open returns an Instance based on information stored on disk at path
func Open(path string) (Instance, error) {
	meta, err := metadata.Open(path)
	if err != nil {
		return nil, err
	}

	sectionEntries, err := sections.Open(meta)
	if err != nil {
		return nil, err
	}

	checks := []check.That{check.ThatAllBasicInvariantsAreSatisfied}
	switch meta.State() {
	case states.Provisional:
		checks = append(checks, check.ThatIsValidForProvisionalState)
	case states.Implementable:
		checks = append(checks, check.ThatIsValidForImplementableState)
	}

	checkAll := check.All(checks)
	err = checkAll(meta)
	if err != nil {
		return nil, err
	}

	k := &kep{
		meta:    meta,
		locker:  new(sync.RWMutex),
		content: make(map[sections.Entry]bool),
		checks:  checks,
	}

	k.addSections(sectionEntries)

	return k, nil
}

type kep struct {
	meta    metadata.KEP
	content map[sections.Entry]bool
	checks  []check.That
	locker  *sync.RWMutex
}

func (k *kep) Persist() error {
	k.locker.Lock()
	defer k.locker.Unlock()

	entries := []sections.Entry{}
	for entry := range k.content {
		if sections.IsAutogenerated(entry.Name()) {
			continue
		}

		entries = append(entries, entry)
	}

	err := sections.Persist(entries)
	if err != nil {
		return err
	}

	sectionLocations := sections.Locations(entries)
	k.meta.AddSectionLocations(sectionLocations)

	autogeneratedEntries, err := sections.AutoGeneratedFrom(k.meta)
	if err != nil {
		return err
	}

	err = sections.Persist(autogeneratedEntries)
	if err != nil {
		return err
	}

	autogeneratedLocations := sections.Locations(autogeneratedEntries)
	k.meta.AddSectionLocations(autogeneratedLocations)

	err = k.meta.Persist()
	if err != nil {
		return err
	}

	return k.check() // TODO figure out how to roll back failures "safely"
}

func (k *kep) SetState(state states.Name) error {
	k.locker.Lock()
	defer k.locker.Unlock()

	var err error

	// make sure the KEP is valid before doing anything
	err = k.check()
	if err != nil {
		return err
	}

	switch state {
	case states.Provisional:
		newEntries, err := sections.RenderMissingForProvisionalState(k.meta)
		if err != nil {
			return err
		}

		k.checks = append(k.checks, check.ThatIsValidForProvisionalState) // allow anyone to check that the KEP remains valid
		k.addSections(newEntries)

		k.meta.SetState(states.Provisional)

	case states.Implementable:
		newEntries, err := sections.RenderMissingForImplementableState(k.meta)
		if err != nil {
			return err
		}

		k.checks = append(k.checks, check.ThatIsValidForImplementableState) // allow anyone to check that the KEP remains valid
		k.addSections(newEntries)

		k.meta.SetState(states.Implementable)
	}

	return nil
}

func (k *kep) addSections(entries []sections.Entry) {
	// newly added entires will be persisted in a call to k.Persist()
	for i := range entries {
		k.content[entries[i]] = true
	}

	k.meta.AddSectionLocations(sections.Locations(entries))
}

func (k *kep) Sections() []string {
	k.locker.RLock()
	defer k.locker.RUnlock()

	sectionNames := []string{}
	for entry := range k.content {
		sectionNames = append(sectionNames, entry.Name())
	}

	return sectionNames
}

// Check enforces the consistency rules for an individual KEP
// important checks
// - that all referenced sections exist on disk
// - that all the minimally required sections have been included
func (k *kep) Check() error {
	k.locker.RLock()
	defer k.locker.RUnlock()

	return k.check()
}

func (k *kep) check() error {
	checkAll := check.All(k.checks)
	return checkAll(k.meta)
}

func (k *kep) AddChecks(checks ...check.That) {
	k.locker.Lock()
	defer k.locker.Unlock()

	k.checks = append(k.checks, checks...)
}

func (k *kep) AddApprovers(approvers ...string) {
	k.locker.Lock()
	defer k.locker.Unlock()

	k.meta.AddApprovers(approvers)
}

func (k *kep) AddReviewers(reviewers ...string) {
	k.locker.Lock()
	defer k.locker.Unlock()

	k.meta.AddReviewers(reviewers)
}

func (k *kep) UniqueID() string {
	k.locker.RLock()
	defer k.locker.RUnlock()

	return k.meta.UniqueID()
}

func (k *kep) ShortID() int {
	k.locker.RLock()
	defer k.locker.RUnlock()

	return k.meta.ShortID()
}

func (k *kep) Title() string {
	k.locker.RLock()
	defer k.locker.RUnlock()

	return k.meta.Title()
}

func (k *kep) OwningSIG() string {
	k.locker.RLock()
	defer k.locker.RUnlock()

	return k.meta.OwningSIG()
}

func (k *kep) Authors() []string {
	k.locker.RLock()
	defer k.locker.RUnlock()

	return k.meta.Authors()
}

func (k *kep) State() states.Name {
	k.locker.RLock()
	defer k.locker.RUnlock()

	return k.meta.State()
}

func (k *kep) ContentDir() string {
	k.locker.RLock()
	defer k.locker.RUnlock()

	return k.meta.ContentDir()
}

func (k *kep) Created() time.Time {
	k.locker.RLock()
	defer k.locker.RUnlock()

	return k.meta.Created()
}

func (k *kep) LastUpdated() time.Time {
	k.locker.RLock()
	defer k.locker.RUnlock()

	return k.meta.LastUpdated()
}
