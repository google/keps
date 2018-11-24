package workflow

import (
	"path/filepath"
	"strings"

	"github.com/calebamiles/keps/pkg/keps"
	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/keps/sections"
	"github.com/calebamiles/keps/pkg/keps/skeleton"
	"github.com/calebamiles/keps/pkg/settings"
	"github.com/calebamiles/keps/pkg/sigs"
)

// Init is responsible for minimizing the busy work of creating a KEP.
// Currently it:
//  * creates the initial directory structure
//    - placing KEPs created at the top level of KEP content in a `kubernetes-wide` directory
//    - places KEPs created at the top level of a SIG directory in a `sig-wide` directory
//  * creates initial metadata with required sections
// Unlike other functions in workflow/ we need to return the path explicitly as it may have
// changed from Runtime.TargetDir() for SIG or Kubernetes wide KEPs
func Init(runtime settings.Runtime) (string, error) {
	routingInfo, err := sigs.BuildRoutingFromPath(runtime.ContentRoot(), runtime.TargetDir())
	if err != nil {
		return "", err
	}

	authors := []string{runtime.Principal()}
	title := buildTitleFromPath(filepath.Base(runtime.TargetDir()))

	kepMetadata, err := metadata.New(authors, title, routingInfo)
	if err != nil {
		return "", err
	}

	// kepMetadata satisfies the requirements for rendering sections
	sectionContent, err := sections.ForProvisionalState(kepMetadata)
	if err != nil {
		return "", err
	}

	err = skeleton.Init(kepMetadata)
	if err != nil {
		return "", err
	}

	kep, err := keps.New(kepMetadata, sectionContent)
	if err != nil {
		//TODO erase skeleton if an error occurred
		return "", err
	}

	err = kep.Persist()
	if err != nil {
		//TODO erase skeleton if an error occurred
		return "", err
	}

	return routingInfo.ContentDir(), nil
}

func buildTitleFromPath(p string) string {
	return strings.Title(strings.Replace(strings.Replace(p, "-", " ", -1), "_", " ", -1))

}
