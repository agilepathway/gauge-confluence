package confluence

import (
	"fmt"
	"io/ioutil"
	"regexp"

	gauge "github.com/agilepathway/gauge-confluence/gauge_messages"
	"github.com/agilepathway/gauge-confluence/internal/regex"
	"github.com/agilepathway/gauge-confluence/util"
)

// Spec decorates a Gauge specification so it can be published to Confluence.
type Spec struct {
	path      string // absolute path to the specification file, including the filename
	protoSpec *gauge.ProtoSpec
	markdown  markdown // the spec contents
	gitURL    string   // the URL for the spec on e.g. GitHub, GitLab, Bitbucket
}

// NewSpec returns a new Spec for the spec at the given absolute path.
func NewSpec(absolutePath string, protoSpec *gauge.ProtoSpec, gitURL string) Spec {
	return Spec{absolutePath, protoSpec, readMarkdown(absolutePath), gitURL}
}

func (s *Spec) validate() error {
	if s.heading() == "" {
		return &invalidSpecError{*s}
	}

	return nil
}

func (s *Spec) heading() string {
	return s.protoSpec.SpecHeading
}

func (s *Spec) confluenceFmt() string {
	return s.markdown.confluenceFmt()
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

func readMarkdown(absolutePath string) markdown {
	specBytes, err := ioutil.ReadFile(absolutePath) //nolint:gosec
	util.Fatal(fmt.Sprintf("Error while reading %s file", absolutePath), err)

	return markdown(specBytes)
}
