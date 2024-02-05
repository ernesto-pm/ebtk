package ebtk

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// "epi" == above, over, in addition to
type EpiFile struct {
	fileInfo                 os.FileInfo
	FileName                 string // Filename includes the extension
	FileNameWithoutExtension string
	AbsolutePath             string
	Extension                string // includes a dot, i.e ".jpg"
	Directory                string
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

func (epiFile *EpiFile) ReplaceFileTextContents(newContent string) error {
	err := os.WriteFile(epiFile.AbsolutePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (epiFile *EpiFile) IsImageFile() bool {
	extension := epiFile.Extension
	return extension == ".jpg" ||
		extension == ".jpeg" ||
		extension == ".png" ||
		extension == ".heic" ||
		extension == ".JPG"
}

func (epiFile *EpiFile) IsTxtFile() bool {
	return epiFile.Extension == ".txt"
}

// Returns the contents of the file as a string
func (epiFile *EpiFile) ContentsAsString() (*string, error) {
	bytes, err := os.ReadFile(epiFile.AbsolutePath)
	if err != nil {
		return nil, err
	}

	stringContent := string(bytes)
	return &stringContent, nil
}

// Returns the contents of the file as raw bytes
func (epiFile *EpiFile) ToBytes() ([]byte, error) {
	bytes, err := os.ReadFile(epiFile.AbsolutePath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// Returns the contents of the file as a base64 encoded string (only works with images for now...)
func (epiFile *EpiFile) ToBase64String() (*string, error) {
	bytes, err := epiFile.ToBytes()
	if err != nil {
		return nil, err
	}

	mimeType, err := epiFile.GetMimeType()
	if err != nil {
		return nil, err
	}

	var base64Encoding string
	switch *mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}
	base64Encoding += base64.StdEncoding.EncodeToString(bytes)

	return &base64Encoding, nil
}

// Returns the mime type identified for this file
func (epiFile *EpiFile) GetMimeType() (*string, error) {
	bytes, err := epiFile.ToBytes()
	if err != nil {
		return nil, err
	}

	mimeType := http.DetectContentType(bytes)
	return &mimeType, nil
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

// This method does not modify the
func (epiFile *EpiFile) RenameFile(newFilename string) error {
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
