package commands

import (
	"testing"

	"github.com/olimpias/gvm/commands/mock"
	"github.com/olimpias/gvm/filesystem"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewListCommand(t *testing.T) {
	tests := []struct {
		lister Lister
	}{
		{
			lister: &filesystem.FileManagement{},
		},
	}

	for _, test := range tests {
		lsCommand := NewListCommand(test.lister)
		assert.Equal(t, test.lister, lsCommand.lister, "Got %v, but expected %v", lsCommand.lister, test.lister)
	}
}

func TestListValidate(t *testing.T) {
	lsCommand := ListCommand{}
	assert.Equal(t, nil, lsCommand.Validate())
}

func TestListApply(t *testing.T) {
	tests := []struct {
		name             string
		expectedVersions []string
		expectedErr      error
	}{
		{
			name:             "success without err",
			expectedVersions: []string{"1.4.1", "1.4.0"},
		},
		{
			name:             "failed with err",
			expectedVersions: []string{"1.4.1", "1.4.0"},
			expectedErr:      dummyErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			mockLister := mock.NewMockLister(controller)
			mockLister.EXPECT().ListGoPackageVersions().Return(test.expectedVersions, test.expectedErr).Times(1)
			listCommand := NewListCommand(mockLister)
			outputErr := listCommand.Apply()
			assert.Equal(t, test.expectedErr, outputErr, "Got %v, but expected %v", outputErr, test.expectedErr)
		})
	}
}
