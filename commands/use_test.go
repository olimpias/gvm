package commands

import (
	"os"
	"testing"

	"github.com/olimpias/gvm/commands/mock"
	"github.com/olimpias/gvm/filesystem"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewUseCommand(t *testing.T) {
	tests := []struct {
		inputVersion string
		packageUser  PackageUser
	}{
		{
			inputVersion: "1.4.1",
			packageUser:  &filesystem.FileManagement{},
		},
	}

	for _, test := range tests {
		useCommand := NewUseCommand(test.packageUser, test.inputVersion)
		assert.Equal(t, test.packageUser, useCommand.packageUser, "Got %v, but expected %v", useCommand.packageUser, test.packageUser)
		assert.Equal(t, test.inputVersion, useCommand.version, "Got %v, but expected %v", useCommand.version, test.inputVersion)
	}
}

func TestUseValidate(t *testing.T) {
	tests := []struct {
		name         string
		inputVersion string
		goRootEnv    string
		expected     error
	}{
		{
			name:         "success",
			inputVersion: "1.4.1",
			goRootEnv:    "/test/test",
		},
		{
			name:         "failed due to invalid version",
			inputVersion: "1.asd.4",
			goRootEnv:    "/test/test",
			expected:     filesystem.ErrInvalidVersion,
		},
		{
			name:         "failed due to empty goRoot",
			inputVersion: "1.1.4",
			goRootEnv:    "",
			expected:     filesystem.ErrGORootIsNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := os.Setenv(filesystem.GORooT, test.goRootEnv); err != nil {
				t.Error(err)
			}
			useCommand := UseCommand{version: test.inputVersion}
			outputErr := useCommand.Validate()
			assert.Equal(t, test.expected, outputErr, "Got %v, but expected", outputErr, test.expected)
		})
	}
}

func TestUseApply(t *testing.T) {
	tests := []struct {
		name         string
		inputVersion string
		expectedErr  error
	}{
		{
			name:         "success without err",
			inputVersion: "1.4.1",
		},
		{
			name:         "failed with err",
			inputVersion: "1.4.1",
			expectedErr:  dummyErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			mockPackageUser := mock.NewMockPackageUser(controller)
			mockPackageUser.EXPECT().UseGoPackage(test.inputVersion).Return(test.expectedErr).Times(1)
			useCom := NewUseCommand(mockPackageUser, test.inputVersion)
			outputErr := useCom.Apply()
			assert.Equal(t, test.expectedErr, outputErr, "Got %v, but expected %v", outputErr, test.expectedErr)
		})
	}
}
