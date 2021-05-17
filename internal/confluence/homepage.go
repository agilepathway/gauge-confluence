package confluence

import (
	"fmt"

	"github.com/agilepathway/gauge-confluence/internal/confluence/api"
	"github.com/agilepathway/gauge-confluence/internal/confluence/time"
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

	h := homepage{id: id, created: time.NewTime(created), childless: children == 0, spaceKey: spaceKey, apiClient: a}

	if id == "" {
		return h, fmt.Errorf("could not obtain a homepage ID for space: %s", spaceKey)
	}

	return h, err
}

func (h *homepage) cqlTimeOffset() int {
	// nolint:gomnd
	minOffset := -12 // the latest time zone on earth, 12 hours behind UTC
	maxOffset := 14  // the earliest time zone on earth, 14 hours ahead of UTC

	for o := minOffset; o <= maxOffset; o++ {
		cqlTime := h.created.CQLFormat(o)
		pages := h.apiClient.PagesCreatedAt(cqlTime)

		for _, pg := range pages {
			if pg == h.id {
				return o
			}
		}
	}

	panic("Could not calculate the time offset")
}
