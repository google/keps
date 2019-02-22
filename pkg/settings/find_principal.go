package settings

import (
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	PrincipalEnv = "KEP_PRINCIPAL_GITHUB_HANDLE"
)

func FindPrincipal() (string, error) {
	envPrincipal := os.Getenv(PrincipalEnv)
	if envPrincipal != "" {
		return envPrincipal, nil
	}

	return principal()
}

func principal() (string, error) {
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

	return s.GitHubHandle, nil
}
