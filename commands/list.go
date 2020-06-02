package commands

import (
	"fmt"
	"os"
)

//go:generate mockgen -source=list.go -destination=mock/lister_mock.go -package mock

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
	versions, err := l.lister.ListGoPackageVersions()
	if err != nil {
		return err
	}
	for _, version := range versions {
		os.Stdout.WriteString(fmt.Sprintf("%s\n", version))
	}
	return nil
}
