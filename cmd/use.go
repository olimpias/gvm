package cmd

import "github.com/olimpias/gvm/common"

type UseCommand struct {
	fileManager *common.FileManagement
	version     string
}

func NewUseCommand(fileManager *common.FileManagement, version string) *UseCommand {
	return &UseCommand{fileManager: fileManager, version: version}
}

func (u *UseCommand) Validate() error {
	return nil
}

func (u *UseCommand) Apply() error {
	return nil
}
