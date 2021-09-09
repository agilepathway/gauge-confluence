package strings

import (
	"testing"
)

var stripNonAlphaNumericTests = []struct { //nolint:gochecknoglobals
	input    string
	expected string
}{
	{"github.com/example-user/example-repo", "githubcomexampleuserexamplerepo"},
	{"abc*12?3def", "abc123def"},
	{"ABC*12?3DEF1", "ABC123DEF1"},
}

//nolint:errcheck,gosec
func TestStripNonAlphaNumeric(t *testing.T) {
	for _, tt := range stripNonAlphaNumericTests {
		expected := tt.expected
		actual := StripNonAlphaNumeric(tt.input)

		if expected != actual {
			t.Fatalf(`
	Expected
	%s
	
	but got:
	%s`, expected, actual)
		}
	}
}
