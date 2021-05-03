package confluence

import (
	"path/filepath"

	"github.com/agilepathway/gauge-confluence/internal/gauge"
	"github.com/agilepathway/gauge-confluence/internal/git"
	"github.com/agilepathway/gauge-confluence/util"
)

const (
	childrenMacro = "{children:excerpt=none|all=true}"
)

var projectRoot = util.GetProjectRoot() //nolint:gochecknoglobals

// page encapsulates a Confluence page.
type page struct {
	id           string
	specPath     string
	title        string
	body         string
	parentPageID string
}

func newPage(entry gauge.DirEntry, parentPageID string) (page, error) {
	if entry.IsDir() {
		return newDirPage(entry.Path, parentPageID), nil
	}

	return newSpecPage(entry, parentPageID)
}

// newDirPage initialises a new page encapsulating a directory.
func newDirPage(path, parentPageID string) page {
	return page{"", path, filepath.Base(path), childrenMacro, parentPageID}
}

// newSpecPage initialises a new page encapsulating a Gauge specifiction.
func newSpecPage(entry gauge.DirEntry, parentPageID string) (page, error) {
	spec := NewSpec(entry.Path, git.SpecGitURL(entry.Path, projectRoot))

	err := spec.validate()
	if err != nil {
		return page{}, err
	}

	return page{"", spec.path, spec.heading(), spec.confluenceFmt(), parentPageID}, nil
}
