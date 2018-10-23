package check

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"

	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/sigs"
)

func ThatAllSectionsExistWithContent(meta metadata.KEP) error {
	var errs *multierror.Error

	for _, sectionFilename := range meta.Sections() {
		sectionBytes, err := ioutil.ReadFile(filepath.Join(meta.ContentDir(), sectionFilename))
		switch {
		case os.IsNotExist(err):
			errs = multierror.Append(errs, fmt.Errorf("invalid section: %s. Section does not exist on disk", sectionFilename))
		case err != nil:
			errs = multierror.Append(errs, err)
		case len(sectionBytes) == 0:
			errs = multierror.Append(errs, fmt.Errorf("invalid section: %s. Section contains no content", sectionFilename))
		}
	}

	return errs.ErrorOrNil()
}

func ThatAllSIGsExist(meta metadata.KEP) error {
	var errs *multierror.Error

	allSIGs := []string{meta.OwningSIG()}
	allSIGs = append(allSIGs, meta.ParticipatingSIGs()...)
	for _, sig := range allSIGs {
		if !sigs.Exists(sig) {
			errs = multierror.Append(errs, fmt.Errorf("invalid SIG: %s. No SIG information compiled into binary. Try updating", sig))
		}
	}

	return errs.ErrorOrNil()
}

func ThatTitleIsSet(meta metadata.KEP) error {
	var errs *multierror.Error

	if meta.Title() == "" {
		errs = multierror.Append(errs, errors.New("no title set"))
	}

	return errs.ErrorOrNil()
}

func ThatAuthorsExist(meta metadata.KEP) error {
	var errs *multierror.Error

	if len(meta.Authors()) == 0 {
		errs = multierror.Append(errs, errors.New("no authors listed"))
		return errs
	}

	for _, author := range meta.Authors() {
		if author == "" {
			errs = multierror.Append(errs, errors.New("empty string given for author"))
		}
	}

	return errs.ErrorOrNil()
}

func ThatSubprojectsExist(meta metadata.KEP) error {
	var errs *multierror.Error

	for _, subproject := range meta.AffectedSubprojects() {
		if !sigs.SubprojectExists(subproject) {
			errs = multierror.Append(errs, fmt.Errorf("invalid subproject: %s. No SIG information compiled into binary. Try updating.", subproject))
		}
	}

	return errs.ErrorOrNil()
}

func ThatKEPHasUUID(meta metadata.KEP) error {
	var errs *multierror.Error

	if meta.UniqueID() == "" {
		errs = multierror.Append(errs, errors.New("empty string given as UUID"))
	}

	_, err := uuid.Parse(meta.UniqueID())
	errs = multierror.Append(errs, err)

	return errs.ErrorOrNil()
}

func ThatCreatedTimeExists(meta metadata.KEP) error {
	var errs *multierror.Error

	emptyTime := time.Time{}

	if !meta.Created().After(emptyTime) {
		errs = multierror.Append(errs, errors.New("created at time is invalid: not before empty time"))
	}

	return errs.ErrorOrNil()
}

func ThatLastUpdatedAfterCreated(meta metadata.KEP) error {
	var errs *multierror.Error

	emptyTime := time.Time{}

	if !meta.Created().After(emptyTime) {
		errs = multierror.Append(errs, errors.New("created at time is invalid. Created at time is not before empty time"))
	}

	if !meta.LastUpdated().After(meta.Created()) {
		errs = multierror.Append(errs, errors.New("created at or last updated time is invalid. Created at time is after last updated"))
	}

	return errs.ErrorOrNil()
}
