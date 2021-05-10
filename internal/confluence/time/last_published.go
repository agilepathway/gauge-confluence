package time

// LastPublishedPropertyKey is the name of the Confluence content property key used to store the last published time
const LastPublishedPropertyKey = "lastPublished"

// LastPublished represents the last time Confluence specs were published for a space
type LastPublished struct {
	Time    Time
	Version int // version is the Confluence content property version, which is incremented on every publish
}

// NewLastPublished creates a new LastPublished from the given time in Confluence format and version.
func NewLastPublished(lastPublishedInConfluenceFormat string, version int) LastPublished {
	if lastPublishedInConfluenceFormat == "" {
		return LastPublished{}
	}

	return LastPublished{Time: NewTime(lastPublishedInConfluenceFormat), Version: version}
}
