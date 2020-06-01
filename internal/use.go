package internal

import (
	"github.com/olimpias/gvm/common"
)

type UseCommand struct {
	fileManager *common.FileManagement
	version     string
}

func NewUseCommand(fileManager *common.FileManagement, version string) *UseCommand {
	return &UseCommand{fileManager: fileManager, version: version}
}

func (u *UseCommand) Validate() error {
	return common.ValidateOperation()
}

func (u *UseCommand) Apply() error {
	if err := u.fileManager.MoveFiles(u.version); err != nil {
		return err
	}
	//if err := u.fileManager.SetEnvVariable(); err != nil {
	//	return err
	//}
	return nil
}
