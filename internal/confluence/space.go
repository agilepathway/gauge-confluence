package confluence

import (
	"fmt"
	"path/filepath"

	"github.com/agilepathway/gauge-confluence/internal/confluence/time"
)

type space struct {
	key                        string
	homepage                   homepage
	publishedPages             map[string]page // Pages published by current invocation of the plugin, keyed by filepath
	lastPublished              time.LastPublished
	modifiedSinceLastPublished bool
}

// newSpace initialises a new space.
func newSpace(key string) space {
	return space{key: key, publishedPages: make(map[string]page)}
}

func (s *space) isValid() (bool, string) {
	if s.modifiedSinceLastPublished {
		fmt.Println()
		return false, fmt.Sprintf("the space has been modified since the last publish. Space key: %s", s.key)
	}

	if s.homepage.id == "" {
		return false, fmt.Sprintf("could not obtain a homepage ID for space: %s", s.key)
	}

	return true, ""
}

func (s *space) parentPageIDFor(path string) string {
	parentDir := filepath.Dir(path)
	parentPageID := s.publishedPages[parentDir].id

	if parentPageID == "" {
		return s.homepage.id
	}

	return parentPageID
}

// checkForDuplicateTitle returns an error if the given page has the same title as an already published page.
func (s *space) checkForDuplicateTitle(given page) error {
	for _, p := range s.publishedPages {
		if p.title == given.title {
			return &duplicatePageError{p, given}
		}
	}

	return nil
}
