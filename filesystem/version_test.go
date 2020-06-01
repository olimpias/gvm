package filesystem

import "testing"

func TestValidateVersion(t *testing.T) {
	tests := []struct {
		version  string
		expected error
	}{
		{version: "1.14.1"},
		{version: "1.10.1"},
		{version: "blabla", expected: ErrInvalidVersion},
	}
	for _, test := range tests {
		output := ValidateVersion(test.version)
		if test.expected != output {
			t.Errorf("got %v, but expected %v", output, test.expected)
		}
	}
}
