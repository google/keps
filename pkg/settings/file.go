package settings

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func writeSettingsFile(loc string, u *User) error {
	settingsBytes, err := yaml.Marshal(u)
	if err != nil {
		log.Error("marshalling user settings file")
		return err
	}

	err = ioutil.WriteFile(loc, settingsBytes, os.ModePerm)
	if err != nil {
		log.Error("writing user settings file")
		return err
	}

	return nil
}

func findSettingsFile() (string, error) {
	d, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	if d == "" {
		log.Info("no user cache directory found, aborting")
		return "", nil
	}

	settingsFileLocation := filepath.Join(d, Dirname, Filename)
	return settingsFileLocation, nil
}

func readSettingsFile(loc string, u *User) error {
	settingsBytes, err := ioutil.ReadFile(loc)
	if err != nil {
		log.Errorf("reading settings location %s: %s", loc, err)
		return err
	}

	err = yaml.Unmarshal(settingsBytes, u)
	if err != nil {
		log.Error("unmarshalling settings YAML")
		return err
	}

	return nil
}
