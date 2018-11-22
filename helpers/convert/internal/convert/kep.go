package convert

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/calebamiles/keps/helpers/convert/internal/extract"
	"github.com/calebamiles/keps/helpers/convert/internal/metadata"
	"github.com/calebamiles/keps/pkg/keps"
)

func ToCurrent(kepLocation string) (string, error) {
	kepBytes, err := ioutil.ReadFile(kepLocation)
	if err != nil {
		log.Fatalf("could not open KEP location: %s", kepLocation)
	}

	kepMetadataBytes, withoutMetadataBytes, err := extract.Metadata(kepBytes)
	if err != nil {
		// TODO (re)add log warning here
		return "", err
	}

	oldMetadata := &metadata.Old{}
	newMetadata := &metadata.New{}

	err = yaml.Unmarshal(kepMetadataBytes, oldMetadata)
	if err != nil {
		// TODO (re)add log warning here
		return "", err
	}

	newMetadata.TitleField = oldMetadata.Title         // assume this will be set
	newMetadata.AuthorsField = oldMetadata.Authors     // assume this will be set
	newMetadata.StateField = oldMetadata.Status        // assume this will be set
	newMetadata.OwningSIGField = oldMetadata.OwningSIG // assume this will be set

	newMetadata.ParticipatingSIGsField = oldMetadata.ParticipatingSIGs
	newMetadata.ReviewersField = clearAllEmpty(oldMetadata.Reviewers)
	newMetadata.ApproversField = clearAllEmpty(oldMetadata.Approvers)
	newMetadata.ReplacesField = clearAllEmpty(oldMetadata.Replaces)
	newMetadata.SupersededByField = clearAllEmpty(oldMetadata.SupersededBy)

	newMetadata.CreatedField = oldMetadata.CreationDate
	newMetadata.LastUpdatedField = oldMetadata.LastUpdated

	newMetadata.EditorsField = append(newMetadata.EditorsField, oldMetadata.Editor)
	newMetadata.EditorsField = clearAllEmpty(newMetadata.EditorsField) // hack as a slice with an empty string is itself nonempty

	newMetadata.UniqueIDField = uuid.New().String() // will panic on error
	newMetadata.SIGWideField = true                 // only KEPs {0000, 0001, 0001a} exist today as Kubernetes wide
	// TODO add `converted` event to metadata

	kepSections, err := extract.Sections(withoutMetadataBytes)
	if err != nil {
		// TODO (re)add log warning here
		return "", err
	}

	convertedLocation, err := ioutil.TempDir("", filename("kep-conversion-helper-"+newMetadata.TitleField))
	if err != nil {
		// TODO (re)add log warning here
		return "", err
	}

	var errs *multierror.Error
	for sectionName, sectionContent := range kepSections {
		if sectionName == extract.TableOfContentsHeading {
			continue
		}

		newMetadata.SectionsField = append(newMetadata.SectionsField, filename(sectionName)+".md")
		multierror.Append(errs, ioutil.WriteFile(filepath.Join(convertedLocation, filename(sectionName)+".md"), sectionContent, os.ModePerm))
	}

	sort.Sort(BySectionOrder(newMetadata.SectionsField)) // we want earlier sections to appear at the top of the metadata

	if errs.ErrorOrNil() != nil {
		// TODO (re)add log warning here
		return "", err
	}

	newMetadataBytes, err := yaml.Marshal(newMetadata)
	if err != nil {
		// TODO (re)add log warning here
		return "", err
	}

	err = ioutil.WriteFile(filepath.Join(convertedLocation, "metadata.yaml"), newMetadataBytes, os.ModePerm)
	if err != nil {
		// TODO (re)add log warning here
		return "", err
	}

	k, err := keps.Open(convertedLocation)
	err = k.Check()
	if err != nil {
		os.RemoveAll(convertedLocation)
		return "", err
	}

	return convertedLocation, nil
}

func filename(s string) string {
	return strings.ToLower(strings.Replace(strings.TrimPrefix(s, "## "), " ", "_", -1))
}

func clearEmpty(s string) string {
	return strings.Replace(strings.Replace(s, "TBD", "", -1), "n/a", "", -1) // a brief survey suggests that we don't use "N/A" or "tbd"
}

func clearAllEmpty(ss []string) []string {
	allCleaned := []string{}

	for _, s := range ss {
		if cleaned := clearEmpty(s); cleaned != "" {
			allCleaned = append(allCleaned, cleaned)
		}
	}

	return allCleaned
}

var sectionOrder = map[string]int{
	"summary.md":                1,
	"motivation.md":             2,
	"proposal.md":               3,
	"graduation_criteria.md":    4,
	"implementation_history.md": 5,
	"drawbacks.md":              6,
	"alternatives.md":           7,
	"infrastructure_needed.md":  8,
}

type BySectionOrder []string

func (o BySectionOrder) Len() int           { return len(o) }
func (o BySectionOrder) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o BySectionOrder) Less(i, j int) bool { return sectionOrder[o[i]] < sectionOrder[o[j]] }
