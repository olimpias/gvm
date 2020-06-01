package windows

import (
	"errors"
	"os"
)

const (
	GORooT = "GOROOT"
)

var (
	PathNotFound = errors.New("Path is not found")
)

//TODO: add zip existence check here

func GetGoRoot() (string, error) {
	goPath := os.Getenv(GORooT)
	if goPath == "" {
		return "", PathNotFound
	}
	return goPath, nil
}
