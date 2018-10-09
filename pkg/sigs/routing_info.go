package sigs

import (
	"fmt"
	"path/filepath"

	"github.com/calebamiles/keps/pkg/sigs/internal/generated"
)

var _ RoutingInfo = &routingInfo{}

type RoutingInfo interface {
	OwningSIG() string
	AffectedSubprojects() []string
	ParticipatingSIGs() []string
	KubernetesWide() bool
	SIGWide() bool
	ContentDir() string
}

func BuildRoutingFromPath(contentRoot string, p string) (RoutingInfo, error) {
	pathAfterContentRoot, err := filepath.Rel(contentRoot, p)
	if err != nil {
		return nil, err
	}

	pathInfo := generated.InfoForPath[pathAfterContentRoot]
	if pathInfo == nil {
		return nil, fmt.Errorf("unable to determine SIG information for given path: %s", p)
	}

	r := &routingInfo{
		OwningSIGField:      pathInfo.OwningSIG,
		KubernetesWideField: pathInfo.KubernetesWide,
		SIGWideField:        pathInfo.SIGWide,
	}

	if pathInfo.Subproject != "" {
		r.AffectedSubprojectsField = []string{pathInfo.Subproject}
	}

	return r, nil
}

const (
	kubernetesWideDir = "kubernetes-wide"
	sigWideDir        = "sig-wide"
)

type routingInfo struct {
	OwningSIGField           string
	AffectedSubprojectsField []string
	ParticipatingSIGsField   []string
	KubernetesWideField      bool
	SIGWideField             bool
	contentRoot              string
}

func (i *routingInfo) ContentDir() string {
	switch {
	case i.KubernetesWideField:
		return filepath.Join(i.contentRoot, kubernetesWideDir)
	case len(i.AffectedSubprojectsField) > 0:
		return filepath.Join(i.contentRoot, i.OwningSIGField, i.AffectedSubprojectsField[0])
	}

	return filepath.Join(i.contentRoot, i.OwningSIGField, sigWideDir)
}

func (i *routingInfo) OwningSIG() string             { return i.OwningSIGField }
func (i *routingInfo) AffectedSubprojects() []string { return i.AffectedSubprojectsField }
func (i *routingInfo) ParticipatingSIGs() []string   { return i.ParticipatingSIGsField }
func (i *routingInfo) KubernetesWide() bool          { return i.KubernetesWideField }
func (i *routingInfo) SIGWide() bool                 { return i.SIGWideField }
