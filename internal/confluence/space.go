package confluence

import "path/filepath"

type space struct {
	key            string
	homepageID     string
	publishedPages map[string]page // keyed by filepath
}

// newSpace initialises a new space.
func newSpace(key string) space {
	return space{key, "", make(map[string]page)}
}

func (s *space) parentPageIDFor(path string) string {
	parentDir := filepath.Dir(path)
	parentPageID := s.publishedPages[parentDir].id

	if parentPageID == "" {
		return s.homepageID
	}

	return parentPageID
}

// hasDuplicateTitle is true if the given page has the same title as an already published page.
func (s *space) hasDuplicateTitle(page page) bool {
	for _, p := range s.publishedPages {
		if p.title == page.title {
			return true
		}
	}

	return false
}
