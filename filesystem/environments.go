package filesystem

import (
	"errors"
	"os"
)

const (
	GORooT = "GOROOT"
)

var (
	ErrGORootIsNotFound = errors.New("$GOROOT is not found in environmental variables")
)

func getGORoot() (string, error) {
	goPath := os.Getenv(GORooT)
	if goPath == "" {
		return "", ErrGORootIsNotFound
	}
	return goPath, nil
}
