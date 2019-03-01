package cmd

import (
	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/orgs"
	"github.com/calebamiles/keps/pkg/orgs/planctae"
	"github.com/calebamiles/keps/pkg/settings"
)

func Setup(targetPath string) (settings.Runtime, keps.Instance, orgs.Instance, error) {
	contentRoot, err := settings.FindContentRoot()
	if err != nil {
		return nil, nil, nil, err
	}

	// save it now to avoid the expensive look everywhere under $HOME next time
	err = settings.SaveContentRoot(contentRoot)
	if err != nil {
		return nil, nil, nil, err
	}

	principal, err := settings.FindPrincipal()
	if err != nil {
		return nil, nil, nil, err
	}

	runtimeSettings := settings.NewRuntime(contentRoot, targetPath, principal)

	p, err := keps.Path(runtime.ContentRoot(), runtime.TargetDir())
	if err != nil {
		return nil, nil, nil, err
	}

	kep, err := keps.Open(p)
	if err != nil {
		return nil, nil, nil, err
	}

	return runtimeSettings, kep, planctae.Organization(), nil
}
