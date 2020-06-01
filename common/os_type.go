package common

import "runtime"

func isWindowOS() bool {
	return runtime.GOOS == "windows"
}
