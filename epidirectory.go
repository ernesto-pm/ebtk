package ebtk

import (
	"errors"
	"os"
	"path/filepath"
)

type EpiDirectory struct {
	AbsPath string
}

func NewEpiDirectory(absPath string) (*EpiDirectory, error) {
	fileInfo, err := os.Lstat(absPath)
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, errors.New("That path is not a directory")
	}

	return &EpiDirectory{
		AbsPath: absPath,
	}, nil
}

func (epiDirectory EpiDirectory) GetAllFiles() ([]EpiFile, error) {
	var epiFiles []EpiFile

	err := filepath.Walk(epiDirectory.AbsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			epiFile, err := NewEpiFile(path)
			if err != nil {
				return err
			}

			epiFiles = append(epiFiles, *epiFile)
		}

		return nil
	})
	if err != nil {
		return epiFiles, err
	}

	return epiFiles, nil
}

func (epiDirectory EpiDirectory) GetFilesWithExtension(extension string) ([]EpiFile, error) {
	var epiFiles []EpiFile

	err := filepath.Walk(epiDirectory.AbsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			epiFile, err := NewEpiFile(path)
			if err != nil {
				return err
			}

			if epiFile.Extension == extension {
				epiFiles = append(epiFiles, *epiFile)
			}
		}

		return nil
	})
	if err != nil {
		return epiFiles, err
	}

	return epiFiles, nil
}

func (epiDirectory EpiDirectory) GetFilesWithExtensions(extensions []string) ([]EpiFile, error) {
	var epiFiles []EpiFile

	err := filepath.Walk(epiDirectory.AbsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			epiFile, err := NewEpiFile(path)
			if err != nil {
				return err
			}

			found := false
			for _, extension := range extensions {
				if epiFile.Extension == extension {
					found = true
				}
			}

			if found {
				epiFiles = append(epiFiles, *epiFile)
			}
		}

		return nil
	})
	if err != nil {
		return epiFiles, err
	}

	return epiFiles, nil
}

func (epiDirectory EpiDirectory) CopyFiles(destinationAbsPath string) error {
	allFiles, err := epiDirectory.GetAllFiles()
	if err != nil {
		return err
	}

	for _, file := range allFiles {
		_, copyErr := file.CopyToDestination(destinationAbsPath)
		if copyErr != nil {
			return copyErr
		}
	}

	return nil
}

func (epiDirectory EpiDirectory) DeleteDirectoryAndContents() error {
	err := os.RemoveAll(epiDirectory.AbsPath)
	return err
}
