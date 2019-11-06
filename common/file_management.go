package common

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
)

const (
	storePath = ".gvm/"

	home                = "HOME"
	downloadURL         = "https://dl.google.com/go/%s"
	downloadFileVersion = "go%s"
	downloadFileOSArch  = ".%s-%s.tar.gz"
	downloadFileOSArchW = ".%s-%s.zip"
)

type FileManagement struct {
	directoryStorePath string
	osName             string
	archName           string
}

func New() (*FileManagement, error) {
	homePath := os.Getenv(home)
	storePath := fmt.Sprintf("%s/%s/", homePath, storePath)
	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		if err := os.MkdirAll(storePath, os.ModePerm); err != nil {
			return nil, err
		}
	}
	return &FileManagement{directoryStorePath: storePath, osName: runtime.GOOS, archName: runtime.GOARCH}, nil
}

//TODO add progress bar
func (fm *FileManagement) DownloadFile(version string) error {
	fileName := getCompressedFileName(version)
	downloadDirectoryPath := fmt.Sprintf("%s%s", fm.directoryStorePath, fileName)
	out, err := os.Create(downloadDirectoryPath)
	if err != nil {
		return err
	}
	defer out.Close()

	URL := fmt.Sprintf(downloadURL, fileName)
	// Get the data
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func getCompressedFileName(version string) string {
	fileVersionDef := fmt.Sprintf(downloadFileVersion, version)
	osArchDef := fmt.Sprintf(downloadFileOSArch, runtime.GOOS, runtime.GOARCH)
	if isWindowOS() {
		osArchDef = fmt.Sprintf(downloadFileOSArchW, runtime.GOOS, runtime.GOARCH)
	}
	return fmt.Sprintf("%s%s", fileVersionDef, osArchDef)
}

func (fm *FileManagement) DeleteFile(version string) error {
	fileName := getCompressedFileName(version)
	folderPath := fmt.Sprintf("%s/%s", fm.directoryStorePath, fileName)
	return os.Remove(folderPath)
}

func (fm *FileManagement) ListVersions() ([]string, error) {
	files, err := ioutil.ReadDir(fm.directoryStorePath)
	if err != nil {
		return nil, err
	}

	versionNames := make([]string, len(files))
	for i, file := range files {
		fileExtraction := fmt.Sprintf(downloadFileOSArch, runtime.GOOS, runtime.GOARCH)
		if isWindowOS() {
			fileExtraction = fmt.Sprintf(downloadFileOSArchW, runtime.GOOS, runtime.GOARCH)
		}
		if len(fileExtraction) > len(file.Name()) {
			continue
		}
		fileName := file.Name()
		versionNames[i] = fileName[:len(fileName)-len(fileExtraction)]
	}
	return versionNames, nil
}

func (fm *FileManagement) MoveFiles(useVersion string) error {
	return nil
}

func isWindowOS() bool {
	return runtime.GOOS == "windows"
}
