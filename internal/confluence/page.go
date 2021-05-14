package confluence

import (
	"path/filepath"

	"github.com/agilepathway/gauge-confluence/internal/gauge"
	"github.com/agilepathway/gauge-confluence/util"
)

const (
	childrenMacro = "{children:excerpt=none|all=true}"
)

var projectRoot = util.GetProjectRoot() //nolint:gochecknoglobals

// page encapsulates a Confluence page.
type page struct {
	id       string
	title    string
	path     string
	body     string
	parentID string
	isDir    bool
}

func newPage(entry gauge.DirEntry, parentID string, spec Spec) (page, error) {
	if entry.IsDir() {
		return newDirPage(entry.Path, parentID), nil
	}

	return newSpecPage(parentID, spec)
}

// newDirPage initialises a new page encapsulating a directory.
func newDirPage(path, parentID string) page {
	return page{title: filepath.Base(path), path: path, body: childrenMacro, parentID: parentID, isDir: true}
}

// newSpecPage initialises a new page encapsulating a Gauge specifiction.
func newSpecPage(parentID string, spec Spec) (page, error) {
	err := spec.validate()
	if err != nil {
		return page{}, err
	}

	return page{title: spec.heading(), path: spec.path, body: spec.confluenceFmt(), parentID: parentID, isDir: false}, nil
}

func (p *page) isSpec() bool {
	return !p.isDir
}
