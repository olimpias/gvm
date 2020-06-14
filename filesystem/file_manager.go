package filesystem

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
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

//TODO: Refactor code base for mocking!

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
)

var (
	ErrVersionIsNotFound = errors.New("version is not found")
)

type FileManagement struct {
	directoryStorePath string
	osName             string
	archName           string
}

type mockOSOperation interface {
	deleteFile(folderPath string) error
	validateFileExistence(filePath string) error
	createDirectory(directoryPath string) error
	removeDirectory(directoryPath string) error
	executeCommands(name string, args []string) error
	removeDirectoryContents(path string) error
}

func New() (*FileManagement, error) {
	homePath := os.Getenv(home)
	storePath := fmt.Sprintf("%s/%s", homePath, storePath)
	fm := &FileManagement{directoryStorePath: storePath, osName: runtime.GOOS, archName: runtime.GOARCH}
	if err := fm.validateFileExistence(storePath); err != nil {
		if err := fm.createDirectory(storePath); err != nil {
			return nil, err
		}
	}
	return fm, nil
}

func (fm *FileManagement) downloadFileWithURL(URL string) (io.ReadCloser, int64, error) {
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

func (fm *FileManagement) DownloadGoPackage(version string) error {
	fileName := getCompressedFileName(version)
	downloadDirectoryPath := fmt.Sprintf("%s%s", fm.directoryStorePath, fileName)
	out, err := os.OpenFile(downloadDirectoryPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer out.Close()
	URL := fmt.Sprintf(downloadURL, fileName)
	reader, fileSize, err := fm.downloadFileWithURL(URL)
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

func (fm *FileManagement) DeleteGoPackage(version string) error {
	fileName := getCompressedFileName(version)
	folderPath := fmt.Sprintf("%s/%s", fm.directoryStorePath, fileName)
	return fm.deleteFile(folderPath)
}

func (fm *FileManagement) deleteFile(folderPath string) error {
	return os.Remove(folderPath)
}

func (fm *FileManagement) ListGoPackageVersions() ([]string, error) {
	files, err := ioutil.ReadDir(fm.directoryStorePath)
	if err != nil {
		return nil, err
	}

	return fm.iterateOverPackages(files), nil
}

func (fm *FileManagement) iterateOverPackages(goVersionFileInfos []os.FileInfo) []string {
	versionNames := make([]string, 0, len(goVersionFileInfos))
	for _, file := range goVersionFileInfos {
		fileExtraction := fmt.Sprintf(downloadFileOSArch, runtime.GOOS, runtime.GOARCH)
		if isWindowOS() {
			fileExtraction = fmt.Sprintf(downloadFileOSArchW, runtime.GOOS, runtime.GOARCH)
		}
		if len(fileExtraction) > len(file.Name()) {
			continue
		}
		fileName := file.Name()
		versionNames = append(versionNames, fileName[:len(fileName)-len(fileExtraction)])
	}
	return versionNames
}

func (fm *FileManagement) UseGoPackage(version string) error {
	fileName := getCompressedFileName(version)
	filePath := fmt.Sprintf("%s%s", fm.directoryStorePath, fileName)
	if err := fm.validateFileExistence(filePath); err != nil {
		return err
	}

	goroot, err := getGORoot()
	if err != nil {
		return err
	}

	tmpFilePath := fmt.Sprintf("%s%s", fm.directoryStorePath, tmpFile)
	if err := fm.createDirectory(tmpFilePath); err != nil {
		return err
	}

	//TODO: add progress bar
	if err := fm.extractCompressedFile(filePath, tmpFilePath); err != nil {
		return fmt.Errorf("failed to unzip file. Err: %s", err)
	}

	if err := fm.removeDirectoryContents(fmt.Sprintf("%s/", goroot)); err != nil {
		return fmt.Errorf("removing files and directories in goroot folder failed. Err: %s. The goroot path %s", err, goroot)
	}

	if err := fm.copyDirectory(fmt.Sprintf("%s/go/", tmpFilePath), goroot); err != nil {
		return fmt.Errorf("copying the directory failed from tmp location to GoROOT location. Err: %s", err)
	}

	if err := fm.removeDirectory(tmpFilePath); err != nil {
		return fmt.Errorf("removing temperory folder failed. Err: %s. The tmp path %s", err, tmpFilePath)
	}

	return nil
}

func (fm *FileManagement) extractCompressedFile(srcPath, destPath string) error {
	if strings.HasSuffix(srcPath, ".zip") {
		return fm.unzipFile(srcPath, destPath)
	}

	return fm.unTarFile(srcPath, destPath)
}

func (fm *FileManagement) unzipFile(srcPath, destPath string) error {
	reader, err := zip.OpenReader(srcPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
		return err
	}
	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(destPath, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, f.Mode()); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
				return err
			}

			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range reader.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fm *FileManagement) unTarFile(srcPath, destPath string) error {
	zipFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	uncompressedStream, err := gzip.NewReader(zipFile)
	if err != nil {
		return err
	}
	defer uncompressedStream.Close()
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fPath := filepath.Join(destPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(fPath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(fPath)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknow type: %s in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

func (fm *FileManagement) validateFileExistence(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ErrVersionIsNotFound
	}
	return nil
}

func (fm *FileManagement) executeCommands(name string, args []string) error {
	extractCommand := exec.Command(name, args...)
	if err := extractCommand.Run(); err != nil {
		return err
	}

	return nil
}

func (fm *FileManagement) removeDirectoryContents(path string) error {
	d, err := os.Open(path)
	if err != nil {
		return err
	}
	defer d.Close()
	//Return all directories
	names, err := d.Readdirnames(0)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (fm *FileManagement) copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	if !isWindowOS() {
		if err := out.Chmod(0755); err != nil {
			return err
		}
	}

	return out.Sync()
}

func (fm *FileManagement) copyDirectory(src string, dst string) error {
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
		return err
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err := fm.copyDirectory(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err := fm.copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (fm *FileManagement) createDirectory(directoryPath string) error {
	return os.MkdirAll(directoryPath, os.ModePerm)
}

func (fm *FileManagement) removeDirectory(directoryPath string) error {
	return os.RemoveAll(directoryPath)
}

////Private for now. Maybe in the future if user does not have installed go before. It might be use for initial installation...
//func (fm *FileManagement) setEnvVariable() error {
//	homePath := os.Getenv(home)
//	bashProfilePath := fmt.Sprintf("%s/%s", homePath, unixBashProfile)
//	f, err := os.OpenFile(bashProfilePath, os.O_APPEND|os.O_WRONLY, 0644)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//	values, err := ioutil.ReadAll(f)
//	if err != nil {
//		return err
//	}
//	if strings.Contains(string(values), unixExportPath) {
//		return ErrBashProfileAlreadySet
//	}
//	if _, err := f.WriteString(unixExportPath); err != nil {
//		return err
//	}
//	cmd := exec.Command(unixBashSourceCmd, bashProfilePath)
//	return cmd.Run()
//}
