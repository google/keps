package index

import (
	"fmt"

	"github.com/calebamiles/keps/pkg/keps/check"
	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/keps/states"
)

func newThatIdentifiersAreUnique(idx Index) check.That {
	return func(meta metadata.KEP) error {
		if meta.ShortID() == metadata.UnsetShortID {
			return nil
		}

		if idx.HasShortID(meta.ShortID()) {
			return fmt.Errorf("short ID: %d, already exists in index", meta.ShortID())
		}

		return nil
	}
}

func newThatHasIndexableState(_ Index) check.That {
	return func(meta metadata.KEP) error {
		switch meta.State() {
		case states.Implementable:
			return nil
		case states.Implemented:
			return nil
		}

		return fmt.Errorf("cannot add KEP with state: %s to KEP index", meta.State())
	}
}
