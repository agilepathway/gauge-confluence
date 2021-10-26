package time

import (
	"github.com/agilepathway/gauge-confluence/internal/confluence/api"
	"github.com/agilepathway/gauge-confluence/internal/logger"
)

// LastPublishedPropertyKey is the name of the Confluence content property key used to store the last published time
const LastPublishedPropertyKey = "lastPublished"

// LastPublished represents the last time Confluence specs were published for a space
type LastPublished struct {
	Time    Time
	Version int // version is the Confluence content property version, which is incremented on every publish
	Err     error
}

// NewLastPublished creates a new LastPublished from the property key stored in Confluence
func NewLastPublished(apiClient api.Client, homepageID string) LastPublished {
	lastPublishedInConfluenceFormat, version, err := apiClient.LastPublished(homepageID, LastPublishedPropertyKey)

	logger.Debugf(false, "Last published: %s", lastPublishedInConfluenceFormat)
	logger.Debugf(false, "Last published version: %d", version)

	if lastPublishedInConfluenceFormat == "" {
		return LastPublished{Err: err}
	}

	return LastPublished{Time: NewTime(lastPublishedInConfluenceFormat), Version: version, Err: err}
}
