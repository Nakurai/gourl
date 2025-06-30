package utils

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
)

var DataDirPath string

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func CreateDataDir() error {
	// creating the data directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error while fetching home dir path: %v", err)
	}

	DataDirPath = path.Join(homeDir, ".gourl")
	dataDirExist, err := dirExists(DataDirPath)
	if err != nil {
		return fmt.Errorf("error while checking if data dir exist: %v", err)
	}

	if !dataDirExist {
		err := os.MkdirAll(DataDirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error while creating data directory: %v", err)
		}
	}
	return nil
}
