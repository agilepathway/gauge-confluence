package confluence

import (
	bfconfluence "github.com/kentaro-m/blackfriday-confluence"
	"github.com/russross/blackfriday/v2"
)

// markdownToConfluence converts GitHub Flavored Markdown,
// which Gauge specifications are written in,
// into Confluence's own wiki format.
// https://github.com/kentaro-m/blackfriday-confluence
// https://support.atlassian.com/confluence-cloud/docs/insert-confluence-wiki-markup/
func markdownToConfluence(markdown string) string {
	renderer := &bfconfluence.Renderer{}
	bf := blackfriday.New(blackfriday.WithRenderer(renderer), blackfriday.WithExtensions(blackfriday.CommonExtensions))

	return string(renderer.Render(bf.Parse([]byte(markdown))))
}
