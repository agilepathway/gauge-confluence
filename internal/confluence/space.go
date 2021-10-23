package confluence

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/agilepathway/gauge-confluence/internal/confluence/api"
	"github.com/agilepathway/gauge-confluence/internal/confluence/api/http"
	"github.com/agilepathway/gauge-confluence/internal/confluence/time"
	"github.com/agilepathway/gauge-confluence/internal/env"
	"github.com/agilepathway/gauge-confluence/internal/git"
	"github.com/agilepathway/gauge-confluence/internal/logger"
	str "github.com/agilepathway/gauge-confluence/internal/strings"
)

type space struct {
	key                        string
	homepage                   homepage
	publishedPages             map[string]page // Pages published by current invocation of the plugin, keyed by filepath
	lastPublished              time.LastPublished
	modifiedSinceLastPublished bool
	apiClient                  api.Client
	err                        error
}

// newSpace initialises a new space.
func newSpace(apiClient api.Client) space {
	return space{publishedPages: make(map[string]page), apiClient: apiClient}
}

func (s *space) retrieveOrGenerateKey() {
	retrievedKey := os.Getenv("CONFLUENCE_SPACE_KEY")

	if retrievedKey == "" {
		s.key = s.generateKey()
	} else {
		s.key = retrievedKey
	}
}

func (s *space) generateKey() string {
	gitWebURL, err := git.WebURL()
	if err != nil {
		s.err = err
		return ""
	}

	return keyFmt(gitWebURL)
}

func keyFmt(u *url.URL) string {
	hostAndPath := u.Host + u.Path
	alphanumeric := str.StripNonAlphaNumeric(hostAndPath)

	return strings.ToUpper(alphanumeric)
}

func (s *space) checkRequiredConfigVars() {
	env.GetRequired("CONFLUENCE_BASE_URL")
	env.GetRequired("CONFLUENCE_USERNAME")
	env.GetRequired("CONFLUENCE_TOKEN")
}

func (s *space) setup() {
	s.checkRequiredConfigVars()
	s.retrieveOrGenerateKey()
	s.createIfDoesNotAlreadyExist()
	s.checkUnmodifiedSinceLastPublish()
}

func (s *space) checkUnmodifiedSinceLastPublish() {
	if s.err != nil {
		return
	}

	s.homepage = newHomepage(s)
	s.lastPublished = time.NewLastPublished(s.apiClient, s.homepage.id)

	if s.lastPublished.Version == 0 {
		if s.isBlank() {
			return
		}

		s.err = fmt.Errorf("the space must be empty when you publish for the first time. "+
			"It can contain a homepage but no other pages. Space key: %s", s.key)

		return
	}

	cqlTime := s.lastPublished.Time.CQLFormat(s.homepage.cqlTimeOffset())
	s.modifiedSinceLastPublished, s.err = s.apiClient.IsSpaceModifiedSinceLastPublished(s.key, cqlTime)

	if s.modifiedSinceLastPublished {
		s.err = fmt.Errorf("the space has been modified since the last publish. Space key: %s", s.key)
	}
}

func (s *space) createIfDoesNotAlreadyExist() {
	if (s.err != nil) || (s.exists()) {
		return
	}

	logger.Infof(true, "Space with key %s does not already exist, creating it ...", s.key)

	s.createSpace()
}

func (s *space) createSpace() {
	if s.err != nil {
		return
	}

	s.err = s.apiClient.CreateSpace(s.key, s.name(), s.description())

	if s.err != nil {
		e, ok := s.err.(*http.RequestError)
		if ok && e.StatusCode == 403 { //nolint:gomnd
			s.err = fmt.Errorf("the Confluence user %s does not have permission to create the Confluence Space. "+
				"Either rerun the plugin with a user who does have permissions to create the Space, "+
				"or get someone to create the Space manually and then run the plugin again. "+
				"Also check the password or token you supplied for the Confluence user is correct",
				env.GetRequired("CONFLUENCE_USERNAME"))
		}
	}
}

func (s *space) name() string {
	if s.err != nil {
		return ""
	}

	var gitRemoteURLPath string

	gitRemoteURLPath, s.err = git.RemoteURLPath()

	return fmt.Sprintf("Gauge specs for %s", gitRemoteURLPath)
}

func (s *space) description() string {
	if s.err != nil {
		return ""
	}

	gitWebURL, err := git.WebURL()
	if err != nil {
		s.err = err
		return ""
	}

	return fmt.Sprintf("Gauge (https://gauge.org) specifications from %s, "+
		"published automatically by the Gauge Confluence plugin tool "+
		"(https://github.com/agilepathway/gauge-confluence) as living documentation.  "+
		"Do not edit this Space manually.  "+
		"You can use Confluence's Include Macro (https://confluence.atlassian.com/doc/include-page-macro-139514.html) "+
		"to include these specifications in as many of your existing Confluence Spaces as you wish.", gitWebURL)
}

func (s *space) exists() bool {
	doesSpaceExist, err := s.apiClient.DoesSpaceExist(s.key)
	s.err = err

	return doesSpaceExist
}

func (s *space) isBlank() bool {
	totalPagesInSpace, err := s.apiClient.TotalPagesInSpace(s.key)

	logger.Debugf(false, "Total pages in Confluence space prior to publishing: %d", totalPagesInSpace)

	s.err = err

	return totalPagesInSpace <= 1
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

	logger.Debugf(false, "Updating last published version to: %d", s.lastPublished.Version+1)

	return s.apiClient.SetContentProperty(s.homepage.id, time.LastPublishedPropertyKey, value, s.lastPublished.Version+1)
}

func (s *space) deleteAllPagesExceptHomepage() (err error) {
	return s.apiClient.DeleteAllPagesInSpaceExceptHomepage(s.key, s.homepage.id)
}

// deleteEmptyDirPages deletes any pages that the plugin has published to in this run
// that are empty directories
func (s *space) deleteEmptyDirPages() (err error) {
	for s.hasEmptyDirPages() {
		for key, page := range s.emptyDirPages() {
			err = s.apiClient.DeletePage(page.id)
			if err != nil {
				return err
			}

			delete(s.publishedPages, key)
		}
	}

	return nil
}

func (s *space) hasEmptyDirPages() bool {
	return len(s.emptyDirPages()) > 0
}

func (s *space) emptyDirPages() map[string]page {
	emptyDirPages := make(map[string]page)

	for key, page := range s.publishedPages {
		if s.isEmptyDir(page) {
			emptyDirPages[key] = page
		}
	}

	return emptyDirPages
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
