package filesystem

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateOperation(t *testing.T) {
	tests := []struct {
		name      string
		goRootEnv string
		expected  error
	}{
		{
			name:      "success",
			goRootEnv: "/test/test",
		},
		{
			name:      "fail with non set goRoot",
			goRootEnv: "",
			expected:  ErrGORootIsNotFound,
		},
	}

	for _, test := range tests {
		if err := os.Setenv(GORooT, test.goRootEnv); err != nil {
			t.Error(err)
		}
		outputErr := ValidateOperation()
		assert.Equal(t, test.expected, outputErr, "Got %v, but expected %v", outputErr, test.expected)

	}
}
