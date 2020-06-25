package commands

import (
	"github.com/olimpias/gvm/filesystem"
)

//go:generate mockgen -source=use.go -destination=mock/user_mock.go -package mock

type PackageUser interface {
	UseGoPackage(version string) error
	CheckGoPackageExistence(version string) error
}

type UseCommand struct {
	packageUser PackageUser
	downloader  Downloader
	version     string
}

func NewUseCommand(fileManager PackageUser, downloader Downloader, version string) *UseCommand {
	return &UseCommand{packageUser: fileManager, downloader: downloader, version: version}
}

func (u *UseCommand) Validate() error {
	err := filesystem.ValidateOperation()
	if err != nil {
		return err
	}

	return filesystem.ValidateVersion(u.version)
}

func (u *UseCommand) Apply() error {
	err := u.packageUser.CheckGoPackageExistence(u.version)
	switch {
	case err == filesystem.ErrVersionIsNotFound:
		if err := u.downloader.DownloadGoPackage(u.version); err != nil {
			return err
		}
	case err != nil:
		return err
	}

	return u.packageUser.UseGoPackage(u.version)
}
