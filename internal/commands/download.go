package commands

import "github.com/olimpias/gvm/internal/filesystem"

type Downloader interface {
	DownloadGoPackage(version string) error
}

type DLCommand struct {
	version    string
	downloader Downloader
}

func NewDLCommand(downloader Downloader, version string) *DLCommand {
	return &DLCommand{downloader: downloader, version: version}
}

func (i *DLCommand) Validate() error {
	return filesystem.ValidateVersion(i.version)
}

func (i *DLCommand) Apply() error {
	return i.fileManager.DownloadGoPackage(i.version)
}
