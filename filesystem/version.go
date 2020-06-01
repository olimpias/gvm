package filesystem

import (
	"errors"
	"regexp"
)

var (
	versionValidation = regexp.MustCompile("(\\d+\\.)(\\d+\\.)(\\d+)")

	ErrInvalidVersion = errors.New("invalid version")
)

func ValidateVersion(version string) error {
	if !versionValidation.MatchString(version) {
		return ErrInvalidVersion
	}
	return nil
}
