package metadata

import (
	"io/ioutil"
	"path/filepath"
)

type Routing struct {
	OwningSIG           string   `yaml:"owning_sig"`
	ParticipatingSIGs   []string `yaml:"participating_sigs"`
	AffectedSubprojects []string `yaml:"affected_subprojects"`
}

//
//  The expected structure is one of
//	- keps/content/sig-node/kublet/device-allocation/<KEP Content>
//	- keps/content/sig-node/device-allocaiton/<KEP Content>
//
func ExtractRoutingFromPathString(path string) (Routing, error) {
	var sig, subproject string
	componentsFound := map[string]bool{}

	pathComponents := strings.Split(path, filepath.Separator)

	for _, component := range pathComponents {
		switch {
		case sigs.Exists(component) && !componentsFound[sigComponent]:
			// set sig
		case sigs.Exists(component) && componentsFound[sigComponent]:
			// error
		case sigs.SubprojectExists(component) && !componentsFound[subprojectComponent]:
			// set subproject
		case sigs.SubprojectExists(component) && componentsFound[subprojectComponent]:
			// error
		}
	}
}

// BuildRoutingFromDirectoryWalk starts at `root` and looks for
// 	- .sig
//	- .subproject
// files
func BuildRoutingFromDirectoryWalk(root string) (Routing, error) {
	var sig, subproject string
	var fileContent []byte
	var readErr error

	routingComponents := map[string]bool{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// extract just the filename
		switch filepath.Base(path) {
		case sigMetadataFile:
			if routingComponents[sigMetadataFile] {
				return fmt.Errorf("found duplicate sig routing info while walking tree building metadata at %s, previously found SIG was: %s", path, sig)
			}

			fileContent, readErr = ioutil.ReadFile(path)
			if readErr != nil {
				return readErr
			}

			sig = string(fileContent)
			if !sigs.Exists(sig) {
				return fmf.Errorf("found non existant SIG routing info: %s,  at: %s", sig, path)
			}

			routingComponenets[sigMetadataFile] = true

		case subprojectMetadataFile:
			if routingComponents[subprojectMetadataFile] {
				return fmt.Errorf("found duplicate subproject routing info while walking tree at: %s, previously found subproject was: %s", path, subproject)
			}

			fileContent, readErr = ioutil.ReadFile(path)
			if readErr != nil {
				return readErr
			}

			subproject = string(fileContent)
			if !sigs.SubprojectExists(subproject) {
				return fmf.Errorf("found non existant subproject routing info: %s, at: %s", subproject, path)
			}

			routingComponenets[subprojectMetadataFile] = true

		default:
			return nil
		}
	})

	return Routing{
		OwningSIG:           sig,
		AffectedSubprojects: []string{subproject},
	}
}

const (
	sigMetadataFile        = ".sig"
	subprojectMetadataFile = ".subproject"
	sigComponent           = "sig"
	subprojectComponent    = "subproject"
)
