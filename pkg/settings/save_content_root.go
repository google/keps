package settings

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func SaveContentRoot(p string) error {
	settingsFileLocation, err := findSettingsFile()
	if err != nil {
		return err
	}

	settingsDir := filepath.Dir(settingsFileLocation)
	err = os.MkdirAll(settingsDir, os.ModePerm)
	if err != nil {
		log.Error("failed to create settings directory")
		return err
	}

	s := &User{}
	err = readSettingsFile(settingsFileLocation, s)
	switch {
	case os.IsNotExist(err):
		s.ContentRoot = p
		return writeSettingsFile(settingsFileLocation, s)

	case err == nil:
		s.ContentRoot = p
		return writeSettingsFile(settingsFileLocation, s)

	default:
		log.Error("unexpected error saving user settings")
		return err
	}

	return nil
}
