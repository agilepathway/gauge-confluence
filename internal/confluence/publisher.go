// Package confluence publishes Gauge specifications to Confluence.
package confluence

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/agilepathway/gauge-confluence/internal/confluence/api"
	"github.com/agilepathway/gauge-confluence/internal/env"
	"github.com/agilepathway/gauge-confluence/internal/gauge"
)

// Publisher publishes Gauge specifications to Confluence.
type Publisher struct {
	apiClient api.Client
	space     space
}

// NewPublisher instantiates a new Publisher.
func NewPublisher() Publisher {
	spaceKey := env.GetRequired("CONFLUENCE_SPACE_KEY")
	return Publisher{api.NewClient(), newSpace(spaceKey)}
}

// Publish publishes all Gauge specifications under the given paths to Confluence.
func (p *Publisher) Publish(specPaths []string) {
	var err error

	spaceHomepageID, err := p.apiClient.SpaceHomepageID(p.space.key)

	if err != nil {
		fmt.Printf("Failed: %v", err)
		return
	}

	p.space.homepageID = spaceHomepageID

	for _, specPath := range specPaths {
		err = p.publishAllSpecsUnder(specPath)
		if err != nil {
			break
		}
	}

	if err == nil {
		fmt.Printf("Success: published %d specs and directory pages to Confluence", len(p.space.publishedPages))
	} else {
		fmt.Printf("Failed: %v", err)
	}
}

func (p *Publisher) publishAllSpecsUnder(baseSpecPath string) (err error) {
	return filepath.WalkDir(baseSpecPath, p.publishIfDirOrSpec)
}

func (p *Publisher) publishIfDirOrSpec(path string, d fs.DirEntry, err error) error {
	entry := gauge.NewDirEntry(path, d)

	if entry.IsDirOrSpec() {
		return p.publishDirOrSpec(entry)
	}

	return nil
}

func (p *Publisher) publishDirOrSpec(entry gauge.DirEntry) error {
	page, err := newPage(entry, p.space.parentPageIDFor(entry.Path))
	if err != nil {
		return err
	}

	return p.publishPage(page)
}

func (p *Publisher) publishPage(page page) (err error) {
	if p.space.hasDuplicateTitle(page) {
		return fmt.Errorf("duplicate page: %s", page.title)
	}

	publishedPageID, err := p.apiClient.PublishPage(p.space.key, page.title, page.body, page.parentPageID)

	if err != nil {
		return err
	}

	page.id = publishedPageID

	p.space.publishedPages[page.specPath] = page

	return nil
}
