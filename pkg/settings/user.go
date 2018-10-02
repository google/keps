package settings

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
	log "github.com/sirupsen/logrus"
)

const (
	Dirname  = "kep"
	Filename = "settings.yaml"
)

type User struct {
	ContentRoot  string `yaml:"content_root"`
	GitHubHandle string `yaml:"github_handle"`
}

func ContentRoot() (string, error) {
	settingsFileLocation, err := findSettingsFile()
	if err != nil {
		return "", err
	}

	s := &User{}
	err = readSettingsFile(settingsFileLocation, s)
	if err != nil {
		log.Warn("reading user settings file")
		return "", nil
	}

	return s.ContentRoot, nil
}

func SaveContentRoot(p string) error {
	settingsFileLocation, err := findSettingsFile()
	if err != nil {
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

func readSettingsFile(loc string, u *User) error {
	settingsBytes, err := ioutil.ReadFile(loc)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(settingsBytes, u)
	if err != nil {
		log.Error("unmarshalling settings YAML")
		return err
	}

	return nil
}
