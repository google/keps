package skeleton

import (
	"os"
	"path/filepath"
)

func Init(provider dirProvider) error {
	contentDir := provider.ContentDir()

	err := os.MkdirAll(contentDir, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(contentDir, guidesDir), os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Create(filepath.Join(contentDir, guidesDir, gitkeep))
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(contentDir, experienceReportsDir), os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Create(filepath.Join(contentDir, experienceReportsDir, gitkeep))
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(contentDir, assetsDir), os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Create(filepath.Join(contentDir, assetsDir, gitkeep))
	if err != nil {
		return err
	}

	return nil
}

func Erase(provider dirProvider) error {
	contentDir := provider.ContentDir()

	return os.Remove(contentDir)
}

type dirProvider interface {
	ContentDir() string
}

const (
	guidesDir            = "guides"
	experienceReportsDir = "experience_reports"
	assetsDir            = "assets"
	gitkeep              = ".gitkeep"
)
