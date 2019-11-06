package cmd

import (
	"fmt"
	"github.com/olimpias/gvm/common"
)

type ListCommand struct {
	fileManager *common.FileManagement
}

func NewListCommand(fileManager *common.FileManagement) *ListCommand {
	return &ListCommand{fileManager: fileManager}
}

func (l *ListCommand) Validate() error {
	return nil
}

func (l *ListCommand) Apply() error {
	versions, err := l.fileManager.ListVersions()
	if err != nil {
		return err
	}
	for _, version := range versions {
		fmt.Println(version)
	}
	return nil
}
