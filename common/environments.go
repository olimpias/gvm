package common

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

func getGORoot() (string, error) {
	goPath := os.Getenv(GORooT)
	if goPath == "" {
		return "", PathNotFound
	}
	return goPath, nil
}
