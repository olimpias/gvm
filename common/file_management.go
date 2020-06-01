package common

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

const (
	storePath = ".gvm/"
	tmpFile   = ".tmp"

	home                = "HOME"
	downloadURL         = "https://dl.google.com/go/%s"
	downloadFileVersion = "go%s"
	downloadFileOSArch  = ".%s-%s.tar.gz"
	downloadFileOSArchW = ".%s-%s.zip"

	unixBashSourceCmd = "source"
	unixBashProfile   = ".bash_profile"
	unixExportPath    = "export PATH=$PATH:/usr/local/go/bin"

	extractCommand = "tar -C %s -zxvf %s"
)

var (
	ErrNotFound = errors.New("Goroot is not found")

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
	storePath := fmt.Sprintf("%s/%s", homePath, storePath)
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
	goroot, err := getGORoot()
	if err != nil {
		return err
	}
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
	tmpFilePath := fmt.Sprintf("%s%s", fm.directoryStorePath, tmpFile)
	if err := fm.createADirectory(tmpFilePath); err != nil {
		return err
	}
	//TODO: find a lib that gonna do tar operation... So we can visualize the progress in stdout
	command := fmt.Sprintf(extractCommand, tmpFilePath, filePath)
	commands := strings.Split(command, " ")
	extractCommand := exec.Command(commands[0], commands[1:]...)
	if err := extractCommand.Run(); err != nil {
		return fmt.Errorf("tar command failed %s", err)
	}

	if err := CopyDir(fmt.Sprintf("%s/go/", tmpFilePath), goroot); err != nil {
		return fmt.Errorf("CopyDir failed %s", err)
	}

	if err := fm.removeADirectory(tmpFilePath); err != nil {
		return fmt.Errorf("removeADirectory failed %s", err)
	}
	//TODO add for other operating systems...
	return nil
}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}
	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func (fm *FileManagement) createADirectory(directoryPath string) error {
	return os.MkdirAll(directoryPath, os.ModePerm)
}

func (fm *FileManagement) removeADirectory(directoryPath string) error {
	return os.RemoveAll(directoryPath)
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
