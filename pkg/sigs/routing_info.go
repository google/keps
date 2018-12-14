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

func BuildRoutingFromPath(contentRoot string, targetPath string) (RoutingInfo, error) {
	// we assume that targetPath will be of the form
	// sig-node/kubelet/device-plugins
	// where
	// - sig-node/ is the owning SIG
	// - kubelet/ is the subproject
	// - device-plugins/ is the KEP directory

	routingPath := filepath.Dir(targetPath)
	kepDirName := filepath.Base(targetPath)

	pathInfo := generated.InfoForPath[routingPath]
	if pathInfo == nil {
		return nil, fmt.Errorf("unable to determine SIG information for given path: %s", targetPath)
	}

	r := &routingInfo{
		OwningSIGField:      pathInfo.OwningSIG,
		KubernetesWideField: pathInfo.KubernetesWide,
		SIGWideField:        pathInfo.SIGWide,
		contentRoot:         contentRoot,
		kepDirName:          kepDirName,
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
	kepDirName               string
}

func (i *routingInfo) ContentDir() string {
	//TODO save this info when creating the routingInfo
	switch {
	case i.KubernetesWideField:
		return filepath.Join(i.contentRoot, kubernetesWideDir, i.kepDirName)
	case len(i.AffectedSubprojectsField) > 0:
		return filepath.Join(i.contentRoot, i.OwningSIGField, i.AffectedSubprojectsField[0], i.kepDirName)
	}

	return filepath.Join(i.contentRoot, i.OwningSIGField, sigWideDir, i.kepDirName)
}

func (i *routingInfo) OwningSIG() string             { return i.OwningSIGField }
func (i *routingInfo) AffectedSubprojects() []string { return i.AffectedSubprojectsField }
func (i *routingInfo) ParticipatingSIGs() []string   { return i.ParticipatingSIGsField }
func (i *routingInfo) KubernetesWide() bool          { return i.KubernetesWideField }
func (i *routingInfo) SIGWide() bool                 { return i.SIGWideField }
