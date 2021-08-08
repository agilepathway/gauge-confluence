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
	logger.Debugf(true, "Space homepage id: %s", id)
	logger.Debugf(true, "Space homepage number of children: %d", children)
	logger.Debugf(true, "Space homepage created: %v", created)

	h := homepage{id: id, created: time.NewTime(created), childless: children == 0, spaceKey: spaceKey, apiClient: a}

	if id == "" {
		return h, fmt.Errorf("could not obtain a homepage ID for space: %s", spaceKey)
	}

	return h, err
}

// cqlTimeOffset calculates the number of hours that CQL queries are to be offset (against UTC) by,
// for the Confluence instance that specs are being published to.  We calculate this on the fly each
// time because it is not easy at all for the user of the plugin to know the time offset for CQL
// queries required by their Confluence instance - see:
// https://community.atlassian.com/t5/Confluence-questions/How-do-I-pass-a-UTC-time-as-the-value-of-lastModified-in-a-REST/qaq-p/1557903
func (h *homepage) cqlTimeOffset() (int, error) {
	logger.Debugf(true, "Confluence homepage ID is %s for space %s", h.spaceKey, h.id)
	logger.Debugf(true, "Homepage created at: %v (UTC)", h.created)
	// nolint:gomnd
	minOffset := -12 // the latest time zone on earth, 12 hours behind UTC
	maxOffset := 14  // the earliest time zone on earth, 14 hours ahead of UTC

	var offset int

	err := errors.Retry(5, 1000, func() (err error) { //nolint:gomnd
		for o := minOffset; o <= maxOffset; o++ {
			cqlTime := h.created.CQLFormat(o)

			if h.apiClient.WasPageCreatedAt(cqlTime, h.id) {
				logger.Debugf(true, "Successfully calculated time offset for Confluence CQL searches: UTC %+d hours", o)
				offset = o
				return
			}
		}
		return fmt.Errorf("could not calculate time offset for Confluence CQL searches")
	})

	return offset, err
}
