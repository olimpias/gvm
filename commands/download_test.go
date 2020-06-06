package commands

import (
	"testing"

	"github.com/olimpias/gvm/commands/mock"
	"github.com/olimpias/gvm/filesystem"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewDLCommand(t *testing.T) {
	tests := []struct {
		inputVersion string
		downloader   Downloader
	}{
		{
			inputVersion: "1.4.1",
			downloader:   &filesystem.FileManagement{},
		},
	}

	for _, test := range tests {
		dlCommand := NewDLCommand(test.downloader, test.inputVersion)
		assert.Equal(t, test.downloader, dlCommand.downloader, "Got %v, but expected %v", dlCommand.downloader, test.downloader)
		assert.Equal(t, test.inputVersion, dlCommand.version, "Got %v, but expected %v", dlCommand.version, test.inputVersion)
	}
}

func TestDLValidate(t *testing.T) {
	tests := []struct {
		name         string
		inputVersion string
		expected     error
	}{
		{
			name:         "success",
			inputVersion: "1.4.1",
		},
		{
			name:         "failed due to invalid version",
			inputVersion: "1.asd.4",
			expected:     filesystem.ErrInvalidVersion,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dlCommand := DLCommand{version: test.inputVersion}
			outputErr := dlCommand.Validate()
			assert.Equal(t, test.expected, outputErr, "Got %v, but expected", outputErr, test.expected)
		})
	}
}

func TestDLApply(t *testing.T) {
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
			mockDownloader := mock.NewMockDownloader(controller)
			mockDownloader.EXPECT().DownloadGoPackage(test.inputVersion).Return(test.expectedErr).Times(1)
			dlCom := NewDLCommand(mockDownloader, test.inputVersion)
			outputErr := dlCom.Apply()
			assert.Equal(t, test.expectedErr, outputErr, "Got %v, but expected %v", outputErr, test.expectedErr)
		})
	}
}
