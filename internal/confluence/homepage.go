package confluence

import (
	"fmt"

	"github.com/agilepathway/gauge-confluence/internal/confluence/api"
	"github.com/agilepathway/gauge-confluence/internal/confluence/time"
	"github.com/agilepathway/gauge-confluence/internal/errors"
	"github.com/agilepathway/gauge-confluence/internal/logger"
)

type homepage struct {
	id        string
	created   time.Time
	childless bool
	spaceKey  string
	apiClient api.Client
}

func newHomepage(spaceKey string, a api.Client) (homepage, error) {
	id, children, created, err := a.SpaceHomepage(spaceKey)
	logger.Debugf(false, "Space homepage id: %s", id)
	logger.Debugf(false, "Space homepage number of children: %d", children)
	logger.Debugf(false, "Space homepage created: %v", created)

	if id == "" {
		return homepage{}, fmt.Errorf("the Confluence space with key %s has no homepage - "+
			"add a homepage manually in Confluence to the space, then try again", spaceKey)
	}

	h := homepage{id: id, created: time.NewTime(created), childless: children == 0, spaceKey: spaceKey, apiClient: a}

	return h, err
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
