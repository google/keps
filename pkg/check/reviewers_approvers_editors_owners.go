package check

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/sigs"
)

func ThatThereAreEditors(meta metadata.KEP) error {
	var errs *multierror.Error

	switch {
	case len(meta.Editors()) == 0:
		errs = multierror.Append(errs, errors.New("no editors"))
	default:
		for _, editor := range meta.Editors() {
			if editor == "" {
				errs = multierror.Append(errs, errors.New("invalid editor. empty string given for editor"))
			}
		}
	}

	return errs.ErrorOrNil()
}

func ThatThereAreReviewers(meta metadata.KEP) error {
	var errs *multierror.Error

	switch {
	case len(meta.Reviewers()) == 0:
		errs = multierror.Append(errs, errors.New("no reviewers"))
	default:
		for _, reviewer := range meta.Reviewers() {
			if reviewer == "" {
				errs = multierror.Append(errs, errors.New("invalid reviewer. empty string given for reviewer"))
			}
		}
	}

	return errs.ErrorOrNil()
}

func ThatThereAreApprovers(meta metadata.KEP) error {
	var errs *multierror.Error

	switch {
	case len(meta.Approvers()) == 0:
		errs = multierror.Append(errs, errors.New("no approvers"))
	default:
		for _, approver := range meta.Approvers() {
			if approver == "" {
				errs = multierror.Append(errs, errors.New("invalid approver. empty string given for approver"))
			}
		}
	}

	return errs.ErrorOrNil()
}

func ThatHasOwningSIG(meta metadata.KEP) error {
	var errs *multierror.Error

	switch {
	case meta.OwningSIG() == "":
		errs = multierror.Append(errs, errors.New("Invalid owning SIG. Empty SIG information"))
	case !sigs.Exists(meta.OwningSIG()):
		errs = multierror.Append(errs, fmt.Errorf("Invalid owning SIG %s. No SIG information compiled in. Try updating?", meta.OwningSIG()))
	}

	return errs.ErrorOrNil()
}

func ThatThereAreOwners(meta metadata.KEP) error {
	var errs *multierror.Error
	var err error

	err = ThatHasOwningSIG(meta)
	errs = multierror.Append(errs, err)

	err = ThatThereAreReviewers(meta)
	errs = multierror.Append(errs, err)

	err = ThatThereAreApprovers(meta)
	errs = multierror.Append(errs, err)

	return errs.ErrorOrNil()
}
