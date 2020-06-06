package commands

import (
	"github.com/olimpias/gvm/filesystem"
)

//go:generate mockgen -source=use.go -destination=mock/user_mock.go -package mock

type PackageUser interface {
	UseGoPackage(version string) error
}

type UseCommand struct {
	packageUser PackageUser
	version     string
}

func NewUseCommand(fileManager PackageUser, version string) *UseCommand {
	return &UseCommand{packageUser: fileManager, version: version}
}

func (u *UseCommand) Validate() error {
	err := filesystem.ValidateOperation()
	if err != nil {
		return err
	}

	return filesystem.ValidateVersion(u.version)
}

func (u *UseCommand) Apply() error {
	return u.packageUser.UseGoPackage(u.version)
}
