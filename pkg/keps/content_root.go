package keps

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/calebamiles/keps/pkg/settings"
	"github.com/calebamiles/keps/pkg/sigs"
)

const (
	contentRootDirectory = "content"
	contentRootHelper    = ".keps"
	ContentRootEnv       = "KEP_CONTENT_ROOT"
)

// FindContentRoot looks for a location with the following structure
// <some-containing-dir>/content/<all-known-sigs>
// returning `content`
func FindContentRoot() (string, error) {
	var foundRoot string

	startLocation, err := contentSearchRoot()
	if err != nil {
		return "", err
	}

	envRoot := os.Getenv(ContentRootEnv)
	cachedRoot, err := settings.ContentRoot()
	if err != nil {
		return "", err
	}

	switch {
	case hasDirForEachSIG(envRoot):
		foundRoot = envRoot
	case hasDirForEachSIG(cachedRoot):
		foundRoot = cachedRoot
	default:
		err = filepath.Walk(startLocation, func(path string, info os.FileInfo, wErr error) error {
			if foundRoot != "" {
				return nil
			}

			if wErr != nil {
				return wErr
			}

			if !info.IsDir() {
				return nil
			}

			if hasDirForEachSIG(path) {
				foundRoot = path
			}

			return nil
		})

		if err != nil {
			return "", err
		}

	}

	if foundRoot == "" {
		return "", fmt.Errorf("could not find KEP content under: %s", startLocation)
	}

	return foundRoot, nil
}

func hasDirForEachSIG(p string) bool {
	knownSIGs := sigs.All()
	for _, s := range knownSIGs {
		if _, sErr := os.Stat(filepath.Join(p, s)); os.IsNotExist(sErr) {
			return false
		}
	}

	return true
}

func contentSearchRoot() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Error("could not get current user information")
		return "", err
	}

	homeDir := usr.HomeDir
	invokedDir, err := os.Getwd()
	if err != nil {
		log.Error("could not get current directory information")
		return "", err
	}

	pathRelativeToHome, err := filepath.Rel(homeDir, invokedDir)
	if err != nil {
		log.Error("finding $PWD relative to $HOME")
		return "", err
	}

	prependedSlashToPath, err := filepath.Abs(pathRelativeToHome)
	if err != nil {
		log.Error("checking that $PWD is under $HOME")
		return "", err
	}

	if prependedSlashToPath == invokedDir {
		log.Error("$PWD seems to not share elements with $HOME")
		return "", fmt.Errorf("file search must start at location under $HOME: %s, not: %s", homeDir, invokedDir)
	}

	pathRelativeToHomeComponents := strings.Split(pathRelativeToHome, string(filepath.Separator))
	startLocation := filepath.Join(homeDir, pathRelativeToHomeComponents[0])

	return filepath.Clean(startLocation), nil
}
