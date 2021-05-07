// Package confluence publishes Gauge specifications to Confluence.
package confluence

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/agilepathway/gauge-confluence/gauge_messages"
	"github.com/agilepathway/gauge-confluence/internal/confluence/api"
	"github.com/agilepathway/gauge-confluence/internal/env"
	"github.com/agilepathway/gauge-confluence/internal/gauge"
	"github.com/agilepathway/gauge-confluence/internal/git"
)

// Publisher publishes Gauge specifications to Confluence.
type Publisher struct {
	apiClient api.Client
	space     space
	specs     map[string]Spec // keyed by filepath
}

// NewPublisher instantiates a new Publisher.
func NewPublisher(m *gauge_messages.SpecDetails) Publisher {
	spaceKey := env.GetRequired("CONFLUENCE_SPACE_KEY")
	return Publisher{api.NewClient(), newSpace(spaceKey), makeSpecsMap(m)}
}

func makeSpecsMap(m *gauge_messages.SpecDetails) map[string]Spec {
	specs := make(map[string]Spec)

	for _, s := range m.Details {
		path := s.Spec.FileName
		specs[path] = NewSpec(path, s.Spec, git.SpecGitURL(path, projectRoot))
	}

	return specs
}

// Publish publishes all Gauge specifications under the given paths to Confluence.
func (p *Publisher) Publish(specPaths []string) {
	var err error

	p.setupSpace()

	if !p.space.isValid() {
		return
	}

	for _, specPath := range specPaths {
		err = p.publishAllSpecsUnder(specPath)
		if err != nil {
			break
		}
	}

	if err != nil {
		p.printFailureMessage(err)
		return
	}

	fmt.Printf("Success: published %d specs and directory pages to Confluence", len(p.space.publishedPages))
}

func (p *Publisher) setupSpace() {
	spaceHomepageID, err := p.apiClient.SpaceHomepageID(p.space.key)
	if err != nil {
		p.printFailureMessage(err)
		return
	}

	p.space.homepageID = spaceHomepageID
}

func (p *Publisher) printFailureMessage(err error) {
	fmt.Printf("Failed: %v", err)
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
	pg, err := newPage(entry, p.space.parentPageIDFor(entry.Path), p.specs[entry.Path])
	if err != nil {
		return err
	}

	return p.publishPage(pg)
}

func (p *Publisher) publishPage(pg page) (err error) {
	err = p.space.checkForDuplicateTitle(pg)
	if err != nil {
		return err
	}

	publishedPageID, err := p.apiClient.PublishPage(p.space.key, pg.title, pg.body, pg.parentID)

	if err != nil {
		return err
	}

	pg.id = publishedPageID

	p.space.publishedPages[pg.path] = pg

	return nil
}
