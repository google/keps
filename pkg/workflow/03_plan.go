package workflow

import (
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/sections"
	"github.com/calebamiles/keps/pkg/settings"
)

func Plan(runtime settings.Runtime) error {
	kep, err := keps.Open(runtime.TargetDir())
	if err != nil {
		return err
	}

	sectionContent, err := sections.ForImplementableState(kep)
	if err != nil {
		return err
	}

	kep.AddSections(sectionContent)
	err = kep.Persist()
	if err != nil {
		return err
	}

	err = kep.Persist()
	if err != nil {
		return err
	}

	// TODO add mechanics for creating PR

	return nil
}
