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
		downloader   Downloader
	}{
		{
			inputVersion: "1.4.1",
			packageUser:  &filesystem.FileManagement{},
			downloader:   &filesystem.FileManagement{},
		},
	}

	for _, test := range tests {
		useCommand := NewUseCommand(test.packageUser, test.downloader, test.inputVersion)
		assert.Equal(t, test.packageUser, useCommand.packageUser, "Got %v, but expected %v", useCommand.packageUser, test.packageUser)
		assert.Equal(t, test.downloader, useCommand.downloader, "Got %v, but expected %v", useCommand.packageUser, test.packageUser)
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := os.Setenv(filesystem.GORooTEnvVariable, test.goRootEnv); err != nil {
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
		name                           string
		inputVersion                   string
		numCheckGoPackageExistenceCall int
		errCheckGoPackageExistence     error
		numDownloaderCall              int
		errDownloader                  error
		numPackageUserCall             int
		errPackageUser                 error
		expectedErr                    error
	}{
		{
			name:                           "success without err",
			inputVersion:                   "1.4.1",
			numCheckGoPackageExistenceCall: 1,
			numPackageUserCall:             1,
		},
		{
			name:                           "success with download",
			inputVersion:                   "1.4.1",
			numCheckGoPackageExistenceCall: 1,
			errCheckGoPackageExistence:     filesystem.ErrVersionIsNotFound,
			numDownloaderCall:              1,
			numPackageUserCall:             1,
		},
		{
			name:                           "failed with file check",
			inputVersion:                   "1.4.1",
			numCheckGoPackageExistenceCall: 1,
			errCheckGoPackageExistence:     dummyErr,
			expectedErr:                    dummyErr,
		},
		{
			name:                           "failed with err from download",
			inputVersion:                   "1.4.1",
			numCheckGoPackageExistenceCall: 1,
			errCheckGoPackageExistence:     filesystem.ErrVersionIsNotFound,
			numDownloaderCall:              1,
			errDownloader:                  dummyErr,
			expectedErr:                    dummyErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			mockPackageUser := mock.NewMockPackageUser(controller)
			mockPackageUser.EXPECT().UseGoPackage(test.inputVersion).Return(test.errPackageUser).Times(test.numPackageUserCall)
			mockPackageUser.EXPECT().CheckGoPackageExistence(test.inputVersion).Return(test.errCheckGoPackageExistence).
				Times(test.numCheckGoPackageExistenceCall)
			mockDownloader := mock.NewMockDownloader(controller)
			mockDownloader.EXPECT().DownloadGoPackage(test.inputVersion).Return(test.errDownloader).Times(test.numDownloaderCall)
			useCom := NewUseCommand(mockPackageUser, mockDownloader, test.inputVersion)
			outputErr := useCom.Apply()
			assert.Equal(t, test.expectedErr, outputErr, "Got %v, but expected %v", outputErr, test.expectedErr)
		})
	}
}
