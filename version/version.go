package version

import (
	"fmt"

	"github.com/olimpias/gvm/logger"
)

const (
	major = 0
	minor = 1
	patch = 1
)

func Print() {
	version := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	logger.Info(fmt.Sprintf("gvm version %s \n", version))
}
