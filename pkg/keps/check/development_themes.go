package check

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/keps/themes"
)

func ThatKEPHasDevelopmentThemes(meta metadata.KEP) error {
	var errs *multierror.Error

	if len(meta.DevelopmentThemes()) == 0 {
		errs = multierror.Append(errs, errors.New("no development themes set"))
	}

	for _, theme := range meta.DevelopmentThemes() {
		if theme == "" {
			errs = multierror.Append(errs, errors.New("Invalid development theme: empty string given as development theme"))
		}
	}

	return errs.ErrorOrNil()
}

func ThatKEPHasStabilityTheme(meta metadata.KEP) error {
	var errs *multierror.Error

	if len(meta.DevelopmentThemes()) == 0 {
		errs = multierror.Append(errs, errors.New("no development themes set"))
	}

	for _, theme := range meta.DevelopmentThemes() {
		if theme == themes.Stability {
			return errs.ErrorOrNil() // should be nil
		}
	}

	errs = multierror.Append(errs, fmt.Errorf("stability development theme: %s not set", themes.Stability))
	return errs
}
