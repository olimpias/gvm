package commands

import (
	"errors"
	"testing"

	"github.com/olimpias/gvm/commands/mock"
	"github.com/olimpias/gvm/filesystem"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	dummyErr = errors.New("dummy err")
)

func TestNewDelCommand(t *testing.T) {
	tests := []struct {
		inputVersion string
		deleter      Deleter
	}{
		{
			inputVersion: "1.4.1",
			deleter:      &filesystem.FileManagement{},
		},
	}

	for _, test := range tests {
		delCommand := NewDelCommand(test.deleter, test.inputVersion)
		assert.Equal(t, test.deleter, delCommand.deleter, "Got %v, but expected %v", delCommand.deleter, test.deleter)
		assert.Equal(t, test.inputVersion, delCommand.version, "Got %v, but expected %v", delCommand.version, test.inputVersion)
	}
}

func TestDelValidate(t *testing.T) {
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
			delCommand := DelCommand{version: test.inputVersion}
			outputErr := delCommand.Validate()
			assert.Equal(t, test.expected, outputErr, "Got %v, but expected", outputErr, test.expected)
		})
	}
}

func TestDelApply(t *testing.T) {
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
			mockDeleter := mock.NewMockDeleter(controller)
			mockDeleter.EXPECT().DeleteGoPackage(test.inputVersion).Return(test.expectedErr).Times(1)
			delCom := NewDelCommand(mockDeleter, test.inputVersion)
			outputErr := delCom.Apply()
			assert.Equal(t, test.expectedErr, outputErr, "Got %v, but expected %v", outputErr, test.expectedErr)
		})
	}
}
