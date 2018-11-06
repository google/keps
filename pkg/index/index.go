package index

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	//bolt  "go.etcd.io/bbolt" // TODO use to implement filtering
	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/states"
	"github.com/calebamiles/keps/pkg/settings"
)

type Index interface {
	ClaimNextShortID() int // implementations of Index MUST tolerate HasShortID being called from within Update() to prevent deadlocks
	HasShortID(int) bool   // implementations of Index MUST tolerate HasShortID being called from within Update() to prevent deadlocks

	Fetch(string) (keps.Instance, error)

	// TODO add Filter(metadata.KEP) metadata.KEP
	// TODO add Remove(keps.Instance)
	Update(keps.Instance) error
	Persist() error
}

func New(contentRoot string) (Index, error) {
	i := &index{
		contentRoot: contentRoot,
		locker:      &sync.RWMutex{},
		kepShortIDs: &sync.Map{},
		kepsSet:     make(map[string]keps.Instance),
	}

	return i, nil
}

// Open returns an existing Index persisted at contentRoot. The index will not be rebuilt so
// callers should not add a *new* KEP to the index. Open will fail if any persisted KEP entry
// cannot be loaded. Open is intended to be used in the context of operations on a *single*
// KEP
func Open(contentRoot string) (Index, error) {
	indexBytes, err := ioutil.ReadFile(filepath.Join(contentRoot, indexFilename))
	if err != nil {
		// TODO add log
		return nil, err
	}

	idx := &index{
		locker:      &sync.RWMutex{},
		kepShortIDs: &sync.Map{},
		kepsSet:     make(map[string]keps.Instance),
	}

	err = yaml.Unmarshal(indexBytes, idx)
	if err != nil {
		println(contentRoot)
		println(string(indexBytes))
		return nil, err
	}

	var allErrs *multierror.Error
	for _, entry := range idx.KEPs {
		k, err := keps.Open(entry.ContentLocationField)
		if err != nil {
			// TODO add log
			allErrs = multierror.Append(allErrs, err)
			continue // collect all possible errors
		}

		err = idx.Update(k)
		if err != nil {
			// TODO add log
			allErrs = multierror.Append(allErrs, err)
			continue // collect all possible errors
		}
	}

	if allErrs.ErrorOrNil() != nil {
		return nil, allErrs
	}

	return idx, nil
}

// Rebuild
// - walks directories starting at runtime.ContentRoot() looking for KEP metadata.yaml
//   files in the tree. For each metadata.yaml that is found
//	- open KEP
//	- add global index consistency checks to KEP
//	- add KEP to index
// - keeps track of the highest short KEP ID encountered for possible future allocation
func Rebuild(runtime settings.Runtime) (Index, error) {
	kepIndex, err := New(runtime.ContentRoot())
	if err != nil {
		return nil, err
	}

	var allErrors *multierror.Error
	err = filepath.Walk(runtime.ContentRoot(), func(path string, info os.FileInfo, incomingErr error) error {
		if incomingErr != nil {
			return incomingErr // we're probably not interested in this error
		}

		if info.Name() != metadataFilename {
			return nil // we're done here
		}

		containingDir := filepath.Dir(path)

		kep, err := keps.Open(containingDir)
		if err != nil {
			log.Errorf("error opening KEP at path: %s, with error: %s", containingDir, err)
			allErrors = multierror.Append(allErrors, err)
			return filepath.SkipDir // keep processing going
		}

		switch kep.State() {
		case states.Draft:
			return filepath.SkipDir
		case states.Provisional:
			// DISCUSS: should we index provisional KEPs
			return filepath.SkipDir
		}

		err = kepIndex.Update(kep)
		if err != nil {
			log.Errorf("error adding KEP at path: %s, to index. Error occurred: %s", containingDir, err)
			allErrors = multierror.Append(allErrors, err)
			return filepath.SkipDir // keep processing going
		}

		// skip rest of directory entries because we already found the metadata
		return filepath.SkipDir
	})

	allErrors = multierror.Append(allErrors, err)
	return kepIndex, allErrors.ErrorOrNil()
}

type index struct {
	KEPs            []*kepEntry              `yaml:"keps"`
	NextNumberField int64                    `yaml:"NEXT_KEP_NUMBER"`
	contentRoot     string                   `yaml:"-"` // do not persist this
	kepShortIDs     *sync.Map                `yaml:"-"` // do not persist this
	locker          *sync.RWMutex            `yaml:"-"` // do not persist this
	kepsSet         map[string]keps.Instance `yaml:"-"` // do not persist this
}

func (i *index) Update(k keps.Instance) error {
	i.locker.Lock()
	defer i.locker.Unlock()

	k.AddChecks(newThatIdentifiersAreUnique(i), newThatHasIndexableState(i))
	err := k.Check()
	if err != nil {
		log.Warnf("could not add kep at path: %s. Failed self-check with error: %s", k.ContentDir(), err)
		return err
	}

	i.kepsSet[k.UniqueID()] = k
	i.kepShortIDs.Store(k.ShortID(), k) // we store k here in anticipation of wanting to print the claiming KEP in cases of conflict

	return nil
}

func (i *index) Fetch(id string) (keps.Instance, error) {
	i.locker.Lock()
	defer i.locker.Unlock()

	k := i.kepsSet[id]
	if k == nil {
		return nil, fmt.Errorf("no KEP with unique ID: %s found", id)
	}

	return k, nil
}

func (i *index) HasShortID(given int) bool {
	// no additional locking should be needed here
	_, found := i.kepShortIDs.Load(given)
	return found
}

func (i *index) ClaimNextShortID() int {
	// no additional locking should be needed here
	return int(atomic.AddInt64(&i.NextNumberField, 1))
}

func (i *index) Persist() error {
	i.locker.Lock()
	defer i.locker.Unlock()

	// we build this here in order to make removing KEPs easy
	entryList := []*kepEntry{}
	for _, k := range i.kepsSet {
		entryList = append(entryList, &kepEntry{
			ShortIDField:         k.ShortID(),
			UUIDField:            k.UniqueID(),
			TitleField:           k.Title(),
			AuthorsField:         k.Authors(),
			ContentLocationField: k.ContentDir(),
			CreatedField:         k.Created(),
			LastUpdatedField:     k.LastUpdated(),
			StateField:           k.State(),
		})
	}

	// TODO add test for this behavior
	i.KEPs = entryList
	sort.Sort(ByIncreasingAge(i.KEPs))

	entriesBytes, err := yaml.Marshal(i)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(i.contentRoot, indexFilename), entriesBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

type kepEntry struct {
	ShortIDField         int         `yaml:"short_id"`
	UUIDField            string      `yaml:"uuid"`
	TitleField           string      `yaml:"title"`
	OwningSIGField       string      `yaml:"owning_sig"`
	AuthorsField         []string    `yaml:"authors"`
	ContentLocationField string      `yaml:"content_location"`
	CreatedField         time.Time   `yaml:"created"`
	LastUpdatedField     time.Time   `yaml:"last_updated"`
	StateField           states.Name `yaml:"state"`
}

type ByIncreasingAge []*kepEntry

func (a ByIncreasingAge) Len() int           { return len(a) }
func (a ByIncreasingAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIncreasingAge) Less(i, j int) bool { return a[i].CreatedField.Before(a[j].CreatedField) }

const (
	indexFilename    = "keps.yaml"
	metadataFilename = "metadata.yaml"
)
