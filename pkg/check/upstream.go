package check

import (
	"github.com/hashicorp/go-multierror"

	"github.com/calebamiles/keps/pkg/keps/metadata"
)

func ThatKEPExistsUpstream(meta metadata.KEP) error {
	var errs *multierror.Error

	panic("not implemented")

	return errs
}

func ThatKEPHasBeenAcceptedUpstream(meta metadata.KEP) error {
	var errs *multierror.Error

	panic("not implemented")

	return errs
}

func ThatKEPNumberIsNotUsedUpstream(meta metadata.KEP) error {
	var errs *multierror.Error

	panic("not implemented")

	return errs
}
