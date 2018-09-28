package metadata

import (
	"fmt"

	"github.com/calebamiles/keps/pkg/sigs"
)

type Routing struct {
	OwningSIG           string   `yaml:"owning_sig"`
	ParticipatingSIGs   []string `yaml:"participating_sigs"`
	AffectedSubprojects []string `yaml:"affected_subprojects"`
}

func NewRoutingFromPath(p string) (*Routing, error) {
	owningSIG := sigs.ExtractNameFromPath(p)
	if owningSIG == "" {
		return nil, fmt.Errorf("no SIG information found in: %s", p)
	}

	r := &Routing{
		OwningSIG: owningSIG,
	}

	// subproject information is optional for now
	targetSubproject := sigs.ExtractSubprojectNameFromPath(p)
	if targetSubproject != "" {
		r.AffectedSubprojects = append(r.AffectedSubprojects, targetSubproject)
	}

	return r, nil
}
