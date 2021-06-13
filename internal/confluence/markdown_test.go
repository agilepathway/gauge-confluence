package confluence

import (
	"testing"
)

const (
	exampleMarkdown = `
# This is the spec heading

    |id| name      |
    |--|-----------|
    |1 | Alice     |
    |2 | Bob       |
    |3 | Eve       |

## This is the scenario heading
* Say "hello" to <name>.
`
	expectedConfluenceFmt = `h1. This is the spec heading
{code}
|id| name      |
|--|-----------|
|1 | Alice     |
|2 | Bob       |
|3 | Eve       |
{code}

h2. This is the scenario heading
* Say "hello" to <name>.

`
)

var markdownTests = []struct { //nolint:gochecknoglobals
	input    markdown
	expected string
}{
	{exampleMarkdown, expectedConfluenceFmt},
}

func TestConfluenceFmt(t *testing.T) {
	for _, tt := range markdownTests {
		expected := tt.expected
		actual := tt.input.confluenceFmt()

		if expected != actual {
			t.Fatalf(`
	Expected
	%s
	
	but got:
	%s`, expected, actual)
		}
	}
}
