package commands

import (
	"fmt"
	"github.com/olimpias/gvm/internal/filesystem"
)

type Lister interface {
	ListGoPackageVersions() ([]string, error)
}

type ListCommand struct {
	lister Lister
}

func NewListCommand(lister Lister) *ListCommand {
	return &ListCommand{lister: lister}
}

func (l *ListCommand) Validate() error {
	return nil
}

func (l *ListCommand) Apply() error {
	versions, err := l.fileManager.ListGoPackageVersions()
	if err != nil {
		return err
	}
	for _, version := range versions {
		fmt.Println(version)
	}
	return nil
}
