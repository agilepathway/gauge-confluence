package confluence

import (
	"fmt"
	"path/filepath"

	"github.com/agilepathway/gauge-confluence/internal/confluence/api"
	"github.com/agilepathway/gauge-confluence/internal/confluence/time"
)

type space struct {
	key                        string
	homepage                   homepage
	publishedPages             map[string]page // Pages published by current invocation of the plugin, keyed by filepath
	lastPublished              time.LastPublished
	modifiedSinceLastPublished bool
	apiClient                  api.Client
}

// newSpace initialises a new space.
func newSpace(key string, apiClient api.Client) space {
	return space{key: key, publishedPages: make(map[string]page), apiClient: apiClient}
}

func (s *space) setup() {
	h, ch, cr, err := s.apiClient.SpaceHomepage(s.key)
	if err != nil {
		s.printFailureMessage(err)
		return
	}

	s.homepage.id = h
	s.homepage.childless = ch == 0
	s.homepage.created = time.NewTime(cr)

	l, err := s.apiClient.LastPublished(s.homepage.id)
	if err != nil {
		s.printFailureMessage(err)
		return
	}

	s.lastPublished = l

	if l.Version == 0 || s.homepage.childless {
		return
	}

	cqlTime := s.lastPublished.Time.FormatTimeForCQLSearch(s.cqlTimeOffset())

	m, err := s.apiClient.IsSpaceModifiedSinceLastPublished(s.key, cqlTime)
	if err != nil {
		s.printFailureMessage(err)
		return
	}

	s.modifiedSinceLastPublished = m
}

func (s *space) cqlTimeOffset() int {
	// nolint:gomnd
	minOffset := -12 // the latest time zone on earth, 12 hours behind UTC
	maxOffset := 14  // the earliest time zone on earth, 14 hours ahead of UTC

	for o := minOffset; o <= maxOffset; o++ {
		cqlTime := s.homepage.created.FormatTimeForCQLSearch(o)
		pages := s.apiClient.PagesCreatedAt(cqlTime)

		for _, pg := range pages {
			if pg == s.homepage.id {
				return o
			}
		}
	}

	panic("Could not calculate the time offset")
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

func (s *space) updateLastPublished() error {
	return s.apiClient.UpdateLastPublished(s.homepage.id, s.lastPublished.Version)
}

func (s *space) printFailureMessage(i interface{}) {
	fmt.Printf("Failed: %v", i)
}
