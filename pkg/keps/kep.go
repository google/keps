package keps

import (
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"

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

	SetState(states.Name) error

	AddChecks(...check.That)
	AddApprovers(...string)
	AddReviewers(...string)
	AddSections(sections.Collection)

	Check() error
	Persist() error
}

// New creates a new Instance from a sections.Collection and a metadata.KEP
func New(meta metadata.KEP, content sections.Collection) (Instance, error) {
	k := &kep{
		meta:    meta,
		content: []sections.Collection{content},
	}

	checks := []check.That{check.ThatAllBasicInvariantsAreSatisfied}
	switch meta.State() {
	case states.Provisional:
		checks = append(checks, check.ThatIsValidForProvisionalState)
	case states.Implementable:
		checks = append(checks, check.ThatIsValidForImplementableState)
	}

	k.checks = checks

	return k, nil
}

// Open returns an Instance based on information stored on disk at path
func Open(path string) (Instance, error) {
	meta, err := metadata.Open(path)
	if err != nil {
		return nil, err
	}

	content, err := sections.OpenCollection(meta)
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

	k := &kep{
		meta:    meta,
		content: []sections.Collection{content},
	}

	return k, nil
}

type kep struct {
	meta    metadata.KEP
	content sectionContent
	checks  []check.That
	locker  sync.RWMutex
}

func (k *kep) Persist() error {
	k.locker.Lock()
	defer k.locker.Unlock()

	// TODO call k.check() before attempting to persist

	// TODO check that sections won't clobber each other before persisting

	var errs *multierror.Error

	errs = multierror.Append(errs, k.content.Persist())
	if errs.ErrorOrNil() != nil {
		return multierror.Append(errs, k.content.Erase())
	}

	k.meta.AddSections(k.content.Sections())
	errs = multierror.Append(errs, k.meta.Persist())
	if errs.ErrorOrNil() != nil {
		return multierror.Append(errs, k.content.Erase())
	}

	return nil
}

func (k *kep) AddSections(newContent sections.Collection) {
	k.locker.Lock()
	defer k.locker.Unlock()

	k.content = append(k.content, newContent)
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

	// we run checks directly here to avoid possibly creating an invalid KEP state by promoting a KEP which is not valid for the desired state
	switch state {
	case states.Provisional:
		err = check.ThatIsValidForProvisionalState(k.meta)
		if err != nil {
			return err
		}

		k.meta.SetState(states.Provisional)
		k.checks = append(k.checks, check.ThatIsValidForProvisionalState) // allow anyone to check that the KEP remains valid
	case states.Implementable:
		err = check.ThatIsValidForImplementableState(k.meta)
		if err != nil {
			return err
		}

		k.meta.SetState(states.Implementable)
		k.checks = append(k.checks, check.ThatIsValidForImplementableState) // allow anyone to check that the KEP remains valid
	}

	return nil
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

func (k *kep) check() error {
	var errs *multierror.Error

	for _, c := range k.checks {
		errs = multierror.Append(errs, c(k.meta))
	}

	return errs.ErrorOrNil()
}

type sectionContent []sections.Collection

func (content sectionContent) Persist() error {
	var errs *multierror.Error

	for _, c := range content {
		errs = multierror.Append(errs, c.Persist())
	}

	return errs.ErrorOrNil()
}

func (content sectionContent) Erase() error {
	var errs *multierror.Error

	for _, c := range content {
		errs = multierror.Append(errs, c.Erase())
	}

	return errs.ErrorOrNil()
}

func (content sectionContent) Sections() []string {
	allSections := []string{}

	for _, c := range content {
		allSections = append(allSections, c.Sections()...)
	}

	return allSections
}
