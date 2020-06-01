package commands

import "github.com/olimpias/gvm/filesystem"

//go:generate mockgen -source=download.go -destination=mock/downloader_mock.go -package mock

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
	return i.downloader.DownloadGoPackage(i.version)
}
