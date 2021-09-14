package confluence

import (
	"net/url"
	"testing"
)

var keyFmtTests = []struct { //nolint:gochecknoglobals
	input    string
	expected string
}{
	{"https://github.com/example-user/example-repo", "GITHUBCOMEXAMPLEUSEREXAMPLEREPO"},
	{"http://github.com/example-user/example-repo", "GITHUBCOMEXAMPLEUSEREXAMPLEREPO"},
	{"http://github.com:8080/example-user/example-repo", "GITHUBCOM8080EXAMPLEUSEREXAMPLEREPO"},
	{"http://example.com/example-user/example-repo", "EXAMPLECOMEXAMPLEUSEREXAMPLEREPO"},
	{"https://example.com/example-user/example-repo", "EXAMPLECOMEXAMPLEUSEREXAMPLEREPO"},
	{"https://example.com/example-user/example-repo/nested", "EXAMPLECOMEXAMPLEUSEREXAMPLEREPONESTED"},
}

//nolint:errcheck,gosec
func TestKeyFmt(t *testing.T) {
	for _, tt := range keyFmtTests {
		expected := tt.expected
		inputURL, _ := url.Parse(tt.input)
		actual := keyFmt(inputURL)

		if expected != actual {
			t.Fatalf(`
	Expected
	%s
	
	but got:
	%s`, expected, actual)
		}
	}
}
