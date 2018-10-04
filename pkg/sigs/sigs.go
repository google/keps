package sigs

import (
	"fmt"
	"path/filepath"

	"github.com/calebamiles/keps/pkg/sigs/internal/generated"
)

type RoutingInfo struct {
	OwningSIG           string
	AffectedSubprojects []string
	ParticipatingSIGs   []string
	KubernetesWide      bool
	SIGWide             bool
}

func All() []string {
	return generated.SIGList
}

func Exists(s string) bool {
	return generated.SIGSet[s]
}

func BuildRoutingFromPath(contentRoot string, p string) (*RoutingInfo, error) {
	pathAfterContentRoot, err := filepath.Rel(contentRoot, p)
	if err != nil {
		return nil, err
	}

	pathInfo := generated.InfoForPath[pathAfterContentRoot]
	if pathInfo == nil {
		return nil, fmt.Errorf("unable to determine SIG information for given path: %s", p)
	}

	r := &RoutingInfo{
		OwningSIG:      pathInfo.OwningSIG,
		KubernetesWide: pathInfo.KubernetesWide,
		SIGWide:        pathInfo.SIGWide,
	}

	if pathInfo.Subproject != "" {
		r.AffectedSubprojects = []string{pathInfo.Subproject}
	}

	return r, nil
}
