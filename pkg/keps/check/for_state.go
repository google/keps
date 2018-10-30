package check

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/calebamiles/keps/pkg/keps/metadata"
)

func ThatIsValidForProvisionalState(meta metadata.KEP) error {
	var errs *multierror.Error
	var err error
	// TODO update Init() to make initial state "proposal/draft"

	err = ThatThereAreOwners(meta)
	errs = multierror.Append(errs, err)

	err = ThatHasAllSectionsForProvisionalState(meta)
	errs = multierror.Append(errs, err)

	if meta.ShortID() != metadata.UnsetShortID {
		errs = multierror.Append(errs, fmt.Errorf("has %d short ID which should be unset for provisional KEPs", meta.ShortID()))
		return errs
	}

	return errs.ErrorOrNil()
}

func ThatIsValidForImplementableState(meta metadata.KEP) error {
	var errs *multierror.Error
	var err error

	err = ThatThereAreOwners(meta)
	errs = multierror.Append(errs, err)

	err = ThatHasAllSectionsForImplementableState(meta)
	errs = multierror.Append(errs, err)

	err = ThatKEPHasBeenAcceptedUpstream(meta)
	errs = multierror.Append(errs, err)

	return errs.ErrorOrNil()
}

func ThatHasAllSectionsForProvisionalState(meta metadata.KEP) error {
	var errs *multierror.Error
	var err error

	err = thatHasIntroduction(meta)
	errs = multierror.Append(errs, err)

	err = thatHasGuides(meta)
	errs = multierror.Append(errs, err)

	return errs.ErrorOrNil()
}

func ThatHasAllSectionsForImplementableState(meta metadata.KEP) error {
	var errs *multierror.Error
	var err error

	err = thatHasIntroduction(meta)
	errs = multierror.Append(errs, err)

	err = thatHasGuides(meta)
	errs = multierror.Append(errs, err)

	err = thatHasAcceptanceCriteria(meta)
	errs = multierror.Append(errs, err)

	return errs.ErrorOrNil()
}

func thatHasIntroduction(meta metadata.KEP) error {
	var errs *multierror.Error

	hasSection := map[string]bool{}
	for _, path := range meta.Sections() {
		hasSection[path] = true
	}

	if !hasSection[motivationFilename] {
		errs = multierror.Append(errs, errors.New("missing Motivation"))

	}

	if !hasSection[summaryFilename] {
		errs = multierror.Append(errs, errors.New("missing Summary"))
	}

	return errs.ErrorOrNil()
}

func thatHasGuides(meta metadata.KEP) error {
	var errs *multierror.Error

	hasSection := map[string]bool{}
	for _, path := range meta.Sections() {
		hasSection[path] = true
	}

	if !hasSection[teachersGuideFilename] {
		errs = multierror.Append(errs, errors.New("missing Teachers Guide"))
	}

	if !hasSection[operatorsGuideFilename] {
		errs = multierror.Append(errs, errors.New("missing Operators Guide"))
	}

	if !hasSection[developersGuideFilename] {
		errs = multierror.Append(errs, errors.New("missing Developers Guide"))
	}

	return errs.ErrorOrNil()
}

func thatHasAcceptanceCriteria(meta metadata.KEP) error {
	var errs *multierror.Error

	hasSection := map[string]bool{}
	for _, path := range meta.Sections() {
		hasSection[path] = true
	}

	if !hasSection[graduationCriteriaFilename] {
		errs = multierror.Append(errs, errors.New("missing Graduation Criteria"))
	}

	return errs.ErrorOrNil()
}

// TODO decide whether to import rendering package to avoid simple errors here
const (
	readmeFilename             = "README.md"
	summaryFilename            = "summary.md"
	motivationFilename         = "motivation.md"
	teachersGuideFilename      = "guides/teacher.md"
	operatorsGuideFilename     = "guides/operator.md"
	developersGuideFilename    = "guides/developer.md"
	graduationCriteriaFilename = "graduation_criteria.md"
)
