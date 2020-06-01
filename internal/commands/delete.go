package commands

import "github.com/olimpias/gvm/common"

type DelCommand struct {
	fileManager *common.FileManagement
	version     string
}

func NewDelCommand(fileManager *common.FileManagement, version string) *DelCommand {
	return &DelCommand{fileManager: fileManager, version: version}
}

func (l *DelCommand) Validate() error {
	return common.ValidateVersion(l.version)
}

func (l *DelCommand) Apply() error {
	return l.fileManager.DeleteFile(l.version)
}
