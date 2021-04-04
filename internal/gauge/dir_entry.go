// Package gauge provides functionality around files in a Gauge project
package gauge

import (
	"io/fs"

	"github.com/agilepathway/gauge-confluence/util"
)

// DirEntry is a directory or file read from a Gauge directory
type DirEntry struct {
	Path string
	d    fs.DirEntry
}

// NewDirEntry initialises a new Gauge DirEntry.
func NewDirEntry(path string, d fs.DirEntry) DirEntry {
	return DirEntry{path, d}
}

// IsDirOrSpec is true if the entry is either a directory or Gauge spec file
func (g *DirEntry) IsDirOrSpec() bool {
	return (g.d.IsDir()) || (g.IsSpec())
}

// IsDir reports whether the entry represents a directory.
func (g *DirEntry) IsDir() bool {
	return (g.d.IsDir())
}

// IsSpec is true if the entry is a Gauge spec file
func (g *DirEntry) IsSpec() bool {
	return util.FileExists(g.Path) && util.IsValidSpecExtension(g.Path) && !util.IsConceptFile(g.Path)
}
