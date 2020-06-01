package common

import "github.com/olimpias/gvm/windows"

func ValidateOperation() error {
	if isWindowOS() {
		_, err := windows.GetGoRoot()
		return err
	}
	//TODO: decide how to validate linux and OSx...
	return nil
}
