package workflow_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/calebamiles/keps/pkg/sigs"
)

func createSIGDirs() (string, error) {
	tmpDir, err := ioutil.TempDir("", "kep-content")
	if err != nil {
		return "", err
	}

	err = createSIGDirsAt(tmpDir)
	if err != nil {
		return "", err
	}

	return tmpDir, nil
}

func createSIGDirsAt(p string) error {
	for _, sig := range sigs.All() {
		err := os.MkdirAll(filepath.Join(p, sig), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
