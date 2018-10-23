package check

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/calebamiles/keps/pkg/keps/metadata"
)

func ThatAuthorIsNotApprover(meta metadata.KEP) error {
	var errs *multierror.Error

	inAuthorsSet := map[string]bool{}

	if len(meta.Authors()) == 0 {
		errs = multierror.Append(errs, errors.New("no authors set"))
	}

	for _, author := range meta.Authors() {
		inAuthorsSet[author] = true
	}

	//TODO downcase comparisons
	for _, approver := range meta.Approvers() {
		if inAuthorsSet[approver] {
			errs = multierror.Append(errs, fmt.Errorf("%s is listed as both an author and approver", approver))
		}
	}

	return errs.ErrorOrNil()
}

func ThatAuthorIsNotReviewer(meta metadata.KEP) error {
	var errs *multierror.Error

	inAuthorsSet := map[string]bool{}

	if len(meta.Authors()) == 0 {
		errs = multierror.Append(errs, errors.New("no authors set"))
	}

	for _, author := range meta.Authors() {
		inAuthorsSet[author] = true
	}

	//TODO downcase comparisons
	for _, reviewer := range meta.Reviewers() {
		if inAuthorsSet[reviewer] {
			errs = multierror.Append(errs, fmt.Errorf("%s is listed as both an author and reviewer", reviewer))
		}
	}

	return errs.ErrorOrNil()
}
