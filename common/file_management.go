package common

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

const (
	storePath = ".gvm/"

	home                = "HOME"
	downloadURL         = "https://dl.google.com/go/%s"
	downloadFileVersion = "go%s"
	downloadFileOSArch  = ".%s-%s.tar.gz"
	downloadFileOSArchW = ".%s-%s.zip"

	unixBashSourceCmd  = "source"
	unixBashProfile    = ".bash_profile"
	unixExportPath     = "export PATH=$PATH:/usr/local/go/bin"
	unixExtractCommand = "tar -C /usr/local -xzf %s"
)

var (
	ErrVersionIsNotFound     = errors.New("version is not found")
	ErrBashProfileAlreadySet = errors.New("bash profile has already been set")
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

func (fm *FileManagement) downloadFile(URL string) (io.ReadCloser, int64, error) {
	// Get the data
	resp, err := http.Get(URL)
	if err != nil {
		return nil, 0, err
	}
	// Check server response
	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("bad status: %s", resp.Status)
	}

	return resp.Body, resp.ContentLength, nil
}

func (fm *FileManagement) DownloadFileWithProgress(version string) error {
	fileName := getCompressedFileName(version)
	downloadDirectoryPath := fmt.Sprintf("%s%s", fm.directoryStorePath, fileName)
	out, err := os.Create(downloadDirectoryPath)
	if err != nil {
		return err
	}
	defer out.Close()
	URL := fmt.Sprintf(downloadURL, fileName)
	reader, fileSize, err := fm.downloadFile(URL)
	if err != nil {
		return err
	}
	defer reader.Close()

	bar := pb.Full.Start64(fileSize)
	defer bar.Finish()
	barReader := bar.NewProxyReader(reader)
	// Writer the body to file
	_, err = io.Copy(out, barReader)
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
	fileName := getCompressedFileName(useVersion)
	filePath := fmt.Sprintf("%s%s", fm.directoryStorePath, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ErrVersionIsNotFound
	}
	command := fmt.Sprintf(unixExtractCommand, filePath)
	if isWindowOS() {
		//TODO handle windows OS.
		/*
			Windows
			The Go project provides two installation options for Windows users (besides installing from source): a zip archive that requires you to set some environment variables and an MSI installer that configures your installation automatically.

			MSI installer
			Open the MSI file and follow the prompts to install the Go tools. By default, the installer puts the Go distribution in c:\Go.

			The installer should put the c:\Go\bin directory in your PATH environment variable. You may need to restart any open command prompts for the change to take effect.

			Zip archive
			Download the zip file and extract it into the directory of your choice (we suggest c:\Go).

			Add the bin subdirectory of your Go root (for example, c:\Go\bin) to your PATH environment variable.

			Setting environment variables under Windows
			Under Windows, you may set environment variables through the "Environment Variables" button on the "Advanced" tab of the "System" control panel. Some versions of Windows provide this control panel through the "Advanced System Settings" option inside the "System" control panel.
		*/
	}
	commands := strings.Split(command, " ")
	cmd := exec.Command(commands[0], commands[1:]...)
	return cmd.Run()
}

func (fm *FileManagement) SetEnvVariable() error {
	homePath := os.Getenv(home)
	bashProfilePath := fmt.Sprintf("%s/%s", homePath, unixBashProfile)
	f, err := os.OpenFile(bashProfilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	values, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if strings.Contains(string(values), unixExportPath) {
		return ErrBashProfileAlreadySet
	}
	if _, err := f.WriteString(unixExportPath); err != nil {
		return err
	}
	cmd := exec.Command(unixBashSourceCmd, bashProfilePath)
	return cmd.Run()
}

func isWindowOS() bool {
	return runtime.GOOS == "windows"
}
