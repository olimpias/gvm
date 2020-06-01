package commands

import (
	"github.com/olimpias/gvm/filesystem"
)

//go:generate mockgen -source=delete.go -destination=mock/deleter_mock.go -package mock

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
	return l.deleter.DeleteGoPackage(l.version)
}
