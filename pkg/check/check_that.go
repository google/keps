package check

import (
	"github.com/hashicorp/go-multierror"

	"github.com/calebamiles/keps/pkg/keps/metadata"
)

// A That provides a simple way of expressing a requirement that a KEP
// (as expressed by its metadata) must satisfy
type That func(metadata.KEP) error

// All combines multiple checks into a single check
func All(checks []That) That {
	return func(meta metadata.KEP) error {
		var errs *multierror.Error
		for _, c := range checks {
			err := c(meta)
			errs = multierror.Append(errs, err)
		}
		return errs.ErrorOrNil()
	}
}
