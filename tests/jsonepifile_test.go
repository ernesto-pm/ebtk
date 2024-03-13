package tests

import (
	"fmt"
	"github.com/ernesto-pm/ebtk/pkg/files"
	"os"
	"path/filepath"
	"testing"
)

func TestNewJSONEpiFile(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}

	filePath := filepath.Join(homeDir, "Desktop", "newFile.json")

	fmt.Println("Trying to load file at path: ", filePath)
	epiFile, err := files.NewEpiFile(filePath)
	if err != nil {
		t.Error(err)
	}

	jsonFile, err := files.TransformFile[files.JSONEpiFile](*epiFile)
	if err != nil {
		t.Error(err)
	}

	// try to parse the epiFile
	parsedMap, err := jsonFile.ParseJSON()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(parsedMap)
}
