package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetConfigDir() (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting work directory: %s", err)
	}
	return filepath.Join(workDir, "config"), nil
}

func GetIconPath() (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting work directory: %s", err)
	}
	return filepath.Join(workDir, "assets", "xonotic.png"), nil
}

func getWorkDir() (string, error) {
	return os.Getwd()
}
