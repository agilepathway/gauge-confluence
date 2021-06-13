package confluence

import (
	"regexp"

	bfconfluence "github.com/kentaro-m/blackfriday-confluence"
	"github.com/russross/blackfriday/v2"
)

type markdown string

// gaugeMarkdownToConfluence converts Gauge specifications
// (which are a flavour of markdown) into Confluence's own wiki format.
// https://github.com/kentaro-m/blackfriday-confluence
// https://support.atlassian.com/confluence-cloud/docs/insert-confluence-wiki-markup/
func (md markdown) confluenceFmt() string {
	mdWithPlaceholders := md.replaceTagsWithPlaceholders()
	confluenceFormatted := md.standardMarkdownToConfluence(mdWithPlaceholders)

	return md.replacePlaceholdersWithTags(confluenceFormatted)
}

func (md markdown) standardMarkdownToConfluence(markdown string) string {
	renderer := &bfconfluence.Renderer{}
	bf := blackfriday.New(blackfriday.WithRenderer(renderer), blackfriday.WithExtensions(blackfriday.CommonExtensions))

	return string(renderer.Render(bf.Parse([]byte(markdown))))
}

func (md markdown) replaceTagsWithPlaceholders() string {
	return regexp.MustCompile(`<(.*?)>`).ReplaceAllString(string(md), "openingtagplaceholder${1}closingtagplaceholder")
}

func (md markdown) replacePlaceholdersWithTags(specWithPlaceholders string) string {
	opening := regexp.MustCompile(`openingtagplaceholder`).ReplaceAllString(specWithPlaceholders, "<")
	return regexp.MustCompile(`closingtagplaceholder`).ReplaceAllString(opening, ">")
}
