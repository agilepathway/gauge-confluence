package confluence

import (
	"fmt"

	"github.com/agilepathway/gauge-confluence/internal/confluence/api"
	"github.com/agilepathway/gauge-confluence/internal/confluence/time"
	"github.com/agilepathway/gauge-confluence/internal/errors"
	"github.com/agilepathway/gauge-confluence/internal/git"
	"github.com/agilepathway/gauge-confluence/internal/logger"
)

type homepage struct {
	id        string
	created   time.Time
	childless bool
	spaceKey  string
	title     string
	apiClient api.Client
	version   int
}

func newHomepage(s *space) (homepage, error) {
	a := s.apiClient
	id, children, created, version, err := a.SpaceHomepage(s.key)
	logger.Debugf(false, "Space homepage id: %s", id)
	logger.Debugf(false, "Space homepage number of children: %d", children)
	logger.Debugf(false, "Space homepage created: %v", created)
	logger.Debugf(false, "Space homepage version: %d", version)

	if err != nil {
		return homepage{}, err
	}

	if id == "" {
		return homepage{}, fmt.Errorf("the Confluence space with key %s has no homepage - "+
			"add a homepage manually in Confluence to the space, then try again", s.key)
	}

	title, err := title(s)

	h := homepage{
		id:        id,
		created:   time.NewTime(created),
		childless: children == 0,
		spaceKey:  s.key,
		title:     title,
		version:   version,
		apiClient: a}

	return h, err
}

func (h *homepage) publish() error {
	newVersion := h.version + 1
	logger.Debugf(true, "Updating Space homepage to version %d ...", newVersion)

	body, err := h.body()
	if err != nil {
		return err
	}

	return h.apiClient.UpdatePage(h.spaceKey, h.id, h.title, body, newVersion)
}

func title(s *space) (string, error) {
	n, err := s.name()
	return fmt.Sprintf("%s Home", n), err
}

func (h *homepage) body() (string, error) {
	gitWebURL, err := git.WebURL()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\\\\ [Gauge|https://gauge.org] specifications from %s, "+
		"published automatically by the "+
		"[Gauge Confluence plugin tool|https://github.com/agilepathway/gauge-confluence] as living documentation.\\\\ \\\\"+
		"Do not edit this Space manually.\\\\ \\\\"+
		"You can use "+
		"[Confluence's Include Macro|https://confluence.atlassian.com/doc/include-page-macro-139514.html] "+
		"to include these specifications in as many of your existing Confluence Spaces as you wish.\\\\ \\\\ \\\\%s",
		gitWebURL, childrenMacro), nil
}

// cqlTimeOffset calculates the number of hours that CQL queries are to be offset (against UTC) by,
// for the Confluence instance that specs are being published to.  We calculate this on the fly each
// time because it is not easy at all for the user of the plugin to know the time offset for CQL
// queries required by their Confluence instance - see:
// https://community.atlassian.com/t5/Confluence-questions/How-do-I-pass-a-UTC-time-as-the-value-of-lastModified-in-a-REST/qaq-p/1557903
func (h *homepage) cqlTimeOffset() (int, error) {
	logger.Debugf(false, "Confluence homepage ID is %s for space %s", h.spaceKey, h.id)
	logger.Debugf(false, "Homepage created at: %v (UTC)", h.created)
	// nolint:gomnd
	minOffset := -12 // the latest time zone on earth, 12 hours behind UTC
	maxOffset := 14  // the earliest time zone on earth, 14 hours ahead of UTC

	var offset int

	err := errors.Retry(5, 1000, func() (err error) { //nolint:gomnd
		for o := minOffset; o <= maxOffset; o++ {
			cqlTime := h.created.CQLFormat(o)
			wasPageCreatedAtCQLTime, err := h.apiClient.WasPageCreatedAt(cqlTime, h.id)
			if err != nil {
				return err
			}

			if wasPageCreatedAtCQLTime {
				logger.Debugf(false, "Successfully calculated time offset for Confluence CQL searches: UTC %+d hours", o)
				offset = o
				return nil
			}
		}
		return fmt.Errorf("could not calculate time offset for Confluence CQL searches")
	})

	return offset, err
}
