package filesystem

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/badgerodon/penv"
)

const (
	GORooTEnvVariable = "GOROOT"
	pathEnvVariable   = "PATH"
	homeEnvVariable   = "HOME"

	defaultPath = "go"

	defaultGoPathSuffix = "go/bin"
	trimGoPathSuffix    = "bin"
)

var (
	ErrGORootIsNotFound = errors.New("$GOROOT is not found in environmental variables")
	ErrGORootMustBeSet  = errors.New("$GOROOT must be set by user")

	ErrGoRootNotFoundInPath = errors.New("")
)

type EnvConfigurator interface {
	GetGoRoot() (string, error)
	GetHomePath() string
	ShouldSetInPathVariable() bool
	SetFilePathToPathVariable(path string) error
}

type EnvVariableManager struct {
}

func (m *EnvVariableManager) GetGoRoot() (string, error) {
	value, err := m.getCurrentGORoot()
	switch {
	case err == ErrGORootIsNotFound:
		if isWindowOS() {
			return "", ErrGORootMustBeSet
		}
	case err != nil:
		return "", err
	}
	if value != "" {
		return value, nil
	}

	return m.decideGoRootPath(), nil
}

func (m *EnvVariableManager) GetHomePath() string {
	return os.Getenv(homeEnvVariable)
}

func (m *EnvVariableManager) decideGoRootPath() string {
	path, err := m.getGoRootFromPathEnv()
	if err == nil {
		return path
	}
	return fmt.Sprintf("%s/.gvm/%s", m.GetHomePath(), defaultPath)
}

func (m *EnvVariableManager) getGoRootFromPathEnv() (string, error) {
	pathValues := os.Getenv(pathEnvVariable)
	paths := strings.Split(pathValues, ":")
	for _, path := range paths {
		if strings.HasSuffix(path, defaultGoPathSuffix) {
			return strings.TrimSuffix(path, trimGoPathSuffix), nil
		}
	}

	return "", ErrGoRootNotFoundInPath
}

func (m *EnvVariableManager) getCurrentGORoot() (string, error) {
	goPath := os.Getenv(GORooTEnvVariable)
	if goPath == "" {
		return "", ErrGORootIsNotFound
	}
	return goPath, nil
}

func (m *EnvVariableManager) ShouldSetInPathVariable() bool {
	_, err := m.getGoRootFromPathEnv()
	return err != nil
}

func (m *EnvVariableManager) SetFilePathToPathVariable(path string) error {
	return penv.AppendEnv(pathEnvVariable, path)
}
