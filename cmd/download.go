package cmd

import "github.com/olimpias/gvm/common"

type DLCommand struct {
	version     string
	fileManager *common.FileManagement
}

func NewDLCommand(fileManager *common.FileManagement, version string) *DLCommand {
	return &DLCommand{fileManager: fileManager, version: version}
}

func (i *DLCommand) Validate() error {
	return common.ValidateVersion(i.version)
}

func (i *DLCommand) Apply() error {
	return i.fileManager.DownloadFile(i.version)
}
