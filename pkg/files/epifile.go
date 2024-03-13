package files

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type EpiFileBase interface {
	ReadContents() ([]byte, error)
	CopyToDestination(destinationAbsPath string) (*EpiFileBase, error)
	Rename(newFilename string) error
}

type EpiFile struct {
	fileInfo                 os.FileInfo
	FileName                 string // Filename includes the extension
	FileNameWithoutExtension string
	AbsolutePath             string
	Extension                string // includes a dot, i.e ".jpg"
	Directory                string
	Type                     string
}

func NewEpiFile(absolutePath string) (*EpiFile, error) {
	fileInfo, err := os.Lstat(absolutePath)
	if err != nil {
		return nil, err
	}

	extension := filepath.Ext(absolutePath)
	fileName := fileInfo.Name()
	fileNameWithoutExtension := strings.ReplaceAll(fileName, extension, "")

	return &EpiFile{
		fileInfo:                 fileInfo,
		FileName:                 fileName,
		Extension:                extension,
		FileNameWithoutExtension: fileNameWithoutExtension,
		AbsolutePath:             absolutePath,
		Directory:                filepath.Dir(absolutePath),
	}, nil
}

// Returns the contents of the file as raw bytes
func (epiFile *EpiFile) ReadContents() ([]byte, error) {
	fmt.Println(epiFile.AbsolutePath)
	bytes, err := os.ReadFile(epiFile.AbsolutePath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// Copies file contents to destination, returning the new copied file as a new instance of EpiFile
func (epiFile *EpiFile) CopyToDestination(destinationAbsPath string) (*EpiFile, error) {
	filename := epiFile.FileName

	contents, err := os.Open(epiFile.AbsolutePath)
	if err != nil {
		return nil, err
	}
	defer contents.Close()

	cloneAbsPath := filepath.Join(destinationAbsPath, filename)
	destFile, err := os.Create(cloneAbsPath)
	if err != nil {
		return nil, err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, contents)
	if err != nil {
		return nil, err
	}

	flushErr := destFile.Sync()
	if flushErr != nil {
		return nil, flushErr
	}

	cloneEpiFile, epiFileErr := NewEpiFile(cloneAbsPath)

	return cloneEpiFile, epiFileErr
}

// Renames this epifile "in place"
func (epiFile *EpiFile) Rename(newFilename string) error {
	originalPath := epiFile.AbsolutePath
	newFilenameWithExtension := fmt.Sprintf("%s%s", newFilename, epiFile.Extension)
	newAbsPath := path.Join(epiFile.Directory, newFilenameWithExtension)

	err := os.Rename(originalPath, newAbsPath)
	if err != nil {
		return err
	}

	newEpiFile, err := NewEpiFile(newAbsPath)
	if err != nil {
		return err
	}

	epiFile.fileInfo = newEpiFile.fileInfo
	epiFile.FileName = newEpiFile.FileName
	epiFile.Extension = newEpiFile.Extension
	epiFile.FileNameWithoutExtension = newEpiFile.FileNameWithoutExtension
	epiFile.AbsolutePath = newEpiFile.AbsolutePath
	epiFile.Directory = newEpiFile.Directory

	return nil
}
