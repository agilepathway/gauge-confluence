package confluence

import (
	"fmt"
	"path/filepath"

	"github.com/agilepathway/gauge-confluence/internal/confluence/api"
	"github.com/agilepathway/gauge-confluence/internal/confluence/time"
	"github.com/agilepathway/gauge-confluence/internal/logger"
)

type space struct {
	key                        string
	homepage                   homepage
	publishedPages             map[string]page // Pages published by current invocation of the plugin, keyed by filepath
	lastPublished              time.LastPublished
	modifiedSinceLastPublished bool
	apiClient                  api.Client
	cqlOffset                  int // Number of hours that CQL queries are to be offset (against UTC) by
}

// newSpace initialises a new space.
func newSpace(key string, apiClient api.Client) space {
	return space{key: key, publishedPages: make(map[string]page), apiClient: apiClient}
}

func (s *space) setup() error {
	h, err := newHomepage(s.key, s.apiClient)
	if err != nil {
		return err
	}

	s.homepage = h
	s.cqlOffset, err = s.homepage.cqlTimeOffset()

	if err != nil {
		return err
	}

	lastPublishedString, version, err := s.apiClient.LastPublished(s.homepage.id, time.LastPublishedPropertyKey)
	if err != nil {
		return err
	}

	logger.Debugf(false, "Last published: %s", lastPublishedString)
	logger.Debugf(false, "Last published version: %d", version)

	s.lastPublished = time.NewLastPublished(lastPublishedString, version)

	if s.lastPublished.Version == 0 {
		blankSpace, err := s.isBlank()

		if err != nil {
			return err
		}

		if blankSpace {
			return nil
		}

		return fmt.Errorf("the space must be empty when you publish for the first time. "+
			"It can contain a homepage but no other pages. Space key: %s", s.key)
	}

	cqlTime := s.lastPublished.Time.CQLFormat(s.cqlOffset)

	m, err := s.apiClient.IsSpaceModifiedSinceLastPublished(s.key, cqlTime)
	if err != nil {
		return err
	}

	s.modifiedSinceLastPublished = m

	if s.modifiedSinceLastPublished {
		return fmt.Errorf("the space has been modified since the last publish. Space key: %s", s.key)
	}

	return nil
}

func (s *space) isBlank() (bool, error) {
	totalPagesInSpace, err := s.apiClient.TotalPagesInSpace(s.key)

	logger.Debugf(false, "Total pages in Confluence space prior to publishing: %d", totalPagesInSpace)

	if err != nil {
		return false, err
	}

	return totalPagesInSpace <= 1, nil
}

func (s *space) parentPageIDFor(path string) string {
	parentDir := filepath.Dir(path)
	parentPageID := s.publishedPages[parentDir].id

	if parentPageID == "" {
		return s.homepage.id
	}

	return parentPageID
}

// Value contains the LastPublished time
type Value struct {
	LastPublished string `json:"lastPublished"`
}

// updateLastPublished stores the time of publishing as a Confluence content property,
// so that in the next run of the plugin it can check that the Confluence space has not
// been edited manually in the meantime.
//
// The content property is attached to the Space homepage rather than to the Space itself, as
// attaching the property to the Space requires admin permissions and we want to allow the
// plugin to be run by non-admin users too.
func (s *space) updateLastPublished() error {
	value := Value{
		LastPublished: time.Now().String(),
	}

	logger.Debugf(false, "updating last published version to: %d", s.lastPublished.Version+1)

	return s.apiClient.SetContentProperty(s.homepage.id, time.LastPublishedPropertyKey, value, s.lastPublished.Version+1)
}

func (s *space) deleteAllPagesExceptHomepage() (err error) {
	return s.apiClient.DeleteAllPagesInSpaceExceptHomepage(s.key, s.homepage.id)
}

// deleteEmptyDirPages deletes any pages that the plugin has published to in this run
// that are empty directories
func (s *space) deleteEmptyDirPages() (err error) {
	for key, page := range s.publishedPages {
		if s.isEmptyDir(page) {
			err = s.apiClient.DeletePage(page.id)
			if err != nil {
				return err
			}

			delete(s.publishedPages, key)
		}
	}

	return nil
}

func (s *space) isEmptyDir(p page) bool {
	return p.isDir && s.isChildless(p)
}

func (s *space) isChildless(p page) bool {
	return len(s.children(p)) == 0
}

func (s *space) children(page page) []string {
	var children []string

	for _, p := range s.publishedPages {
		if page.id == p.parentID {
			children = append(children, p.id)
		}
	}

	return children
}
