package confluence

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/agilepathway/gauge-confluence/internal/regex"
	"github.com/agilepathway/gauge-confluence/util"
)

// Spec decorates a Gauge specification so it can be published to Confluence.
type Spec struct {
	path     string // absolute path to the specification file, including the filename
	markdown string // the spec contents
	gitURL   string // the URL for the specs directory on e.g. GitHub, GitLab, Bitbucket
}

// NewSpec returns a new Spec object for the spec at the given absolute path
func NewSpec(absolutePath string, gitURL string) Spec {
	return Spec{absolutePath, readMarkdown(absolutePath), gitURL}
}

func (s *Spec) addGitLinkAfterSpecHeading(spec string) string {
	if s.gitURL == "" {
		return spec
	}

	replacement := fmt.Sprintf("${1}\n%s\n", s.gitLinkInConfluenceFormat())

	return regex.ReplaceFirstMatch(spec, replacement, regexp.MustCompile(`(h1.*)`))
}

func (s *Spec) gitLinkInConfluenceFormat() string {
	// TODO: change the implementation to be Confluence format, not Jira
	//nolint:godox
	return fmt.Sprintf("[View or edit this spec in Git|%s]", s.gitURL)
}

func readMarkdown(absolutePath string) string {
	specBytes, err := ioutil.ReadFile(absolutePath) //nolint:gosec
	util.Fatal(fmt.Sprintf("Error while reading %s file", absolutePath), err)

	return string(specBytes)
}
