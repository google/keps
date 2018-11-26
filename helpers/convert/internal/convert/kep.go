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
	"github.com/calebamiles/keps/helpers/convert/internal/render"
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/skeleton"
)

func ToCurrent(outputDir string, kepLocation string) (string, error) {
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

	err = os.MkdirAll(filepath.Join(outputDir, strings.Replace(newMetadata.OwningSIGField, " ", "-", -1)), os.ModePerm)
	if err != nil {
		// TODO (re)add log warning here
		return "", err
	}

	convertedLocation, err := ioutil.TempDir(filepath.Join(outputDir, strings.Replace(newMetadata.OwningSIGField, " ", "-", -1)), strings.Replace(newMetadata.TitleField, " ", "-", -1))
	if err != nil {
		// TODO (re)add log warning here
		return "", err
	}

	err = skeleton.Init(dirProvider(convertedLocation))
	if err != nil {
		return "", err
	}

	sectionLocations := map[string]string{}

	var errs *multierror.Error
	for sectionName, sectionContent := range kepSections {
		switch sectionName {
		case extract.TableOfContentsHeading:
			continue // we want to auto generate this in a KEPs README.md
		case extract.ProposalHeading:
			loc := filepath.Join(convertedLocation, developerGuideFilename)
			sectionLocations[developerGuideName] = developerGuideFilename

			newMetadata.SectionsField = append(newMetadata.SectionsField, developerGuideFilename)
			developerGuide, renderErr := render.DeveloperGuide(kepSections)
			if renderErr != nil {
				errs = multierror.Append(errs, renderErr)
				continue // try writing the rest of the sections
			}

			errs = multierror.Append(errs, ioutil.WriteFile(loc, developerGuide, os.ModePerm))
		case extract.DrawbacksHeading:
			// covered by template rendering in ProposalHeading case
		case extract.AlternativesHeading:
			// covered by template rendering in ProposalHeading case
		case extract.ImplementationHistoryHeading:
			loc := filepath.Join(convertedLocation, changelogFilename)
			sectionLocations[changelogName] = changelogFilename

			newMetadata.SectionsField = append(newMetadata.SectionsField, changelogFilename)
			errs = multierror.Append(errs, ioutil.WriteFile(loc, sectionContent, os.ModePerm))
		default:
			loc := filepath.Join(convertedLocation, markdownFilename(sectionName))
			sectionLocations[sectionName] = markdownFilename(sectionName)

			newMetadata.SectionsField = append(newMetadata.SectionsField, markdownFilename(sectionName))
			errs = multierror.Append(errs, ioutil.WriteFile(loc, sectionContent, os.ModePerm))
		}

	}

	sort.Sort(BySectionOrder(newMetadata.SectionsField)) // we want earlier sections to appear at the top of the metadata

	if errs.ErrorOrNil() != nil {
		// TODO (re)add log warning here
		return "", errs
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

	readme, err := render.NewReadme(newMetadata, sectionLocations)
	if err != nil {
		// TODO (re)add log warning here
		return "", err
	}

	err = ioutil.WriteFile(filepath.Join(convertedLocation, "README.md"), readme, os.ModePerm)
	if err != nil {
		// TODO (re)add log warning here
		return "", err
	}

	k, err := keps.Open(convertedLocation)
	if err != nil {
		return "", err
	}

	err = k.Check()
	if err != nil {
		// TODO log location for debugging rather than remove
		os.RemoveAll(convertedLocation)
		return "", err
	}

	return convertedLocation, nil
}

func markdownFilename(s string) string {
	return strings.ToLower(strings.Replace(strings.TrimPrefix(s, "## "), " ", "_", -1)) + ".md"
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

const (
	developerGuideName     = "Developer Guide"
	changelogName          = "CHANGELOG"
	developerGuideFilename = "guides/developer.md"
	changelogFilename      = "CHANGELOG.md"
)

var sectionOrder = map[string]int{
	"summary.md":                1,
	"motivation.md":             2,
	"proposal.md":               3,
	"guides/developer.md":       3, // everything under guides/ is essentially equal in order
	"graduation_criteria.md":    4,
	"implementation_history.md": 5,
	"CHANGELOG.md":              5, // CHANGELOG is the same as implementation_history
	"drawbacks.md":              6,
	"alternatives.md":           7,
	"infrastructure_needed.md":  8,
}

type BySectionOrder []string

func (o BySectionOrder) Len() int           { return len(o) }
func (o BySectionOrder) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o BySectionOrder) Less(i, j int) bool { return sectionOrder[o[i]] < sectionOrder[o[j]] }

type dirProvider string

func (c dirProvider) ContentDir() string { return string(c) }
