package commands

import "github.com/olimpias/gvm/internal/filesystem"

type Deleter interface {
	DeleteGoPackage(version string) error
}

type DelCommand struct {
	deleter Deleter
	version string
}

func NewDelCommand(deleter Deleter, version string) *DelCommand {
	return &DelCommand{deleter: deleter, version: version}
}

func (l *DelCommand) Validate() error {
	return filesystem.ValidateVersion(l.version)
}

func (l *DelCommand) Apply() error {
	return l.fileManager.DeleteGoPackage(l.version)
}
