// Package confluence publishes Gauge specifications to Confluence.
package confluence

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/agilepathway/gauge-confluence/gauge_messages"
	"github.com/agilepathway/gauge-confluence/internal/confluence/api"
	"github.com/agilepathway/gauge-confluence/internal/confluence/time"
	"github.com/agilepathway/gauge-confluence/internal/env"
	"github.com/agilepathway/gauge-confluence/internal/errors"
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
	return Publisher{apiClient: api.NewClient(), space: newSpace(spaceKey), specs: makeSpecsMap(m)}
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

	isSpaceValid, msg := p.space.isValid()

	if !isSpaceValid {
		p.printFailureMessage(msg)
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

	err = p.updateLastPublished()
	if err != nil {
		p.printFailureMessage(err)
		return
	}

	fmt.Printf("Success: published %d specs and directory pages to Confluence", len(p.space.publishedPages))
}

func (p *Publisher) setupSpace() {
	h, ch, cr, err := p.apiClient.SpaceHomepage(p.space.key)
	if err != nil {
		p.printFailureMessage(err)
		return
	}

	p.space.homepageID = h

	p.space.homepageNumberOfChildren = ch

	p.space.homepageCreated = time.NewTime(cr)

	l, err := p.apiClient.LastPublished(p.space.homepageID)
	if err != nil {
		p.printFailureMessage(err)
		return
	}

	p.space.lastPublished = l

	if l.Version == 0 || p.space.homepageNumberOfChildren == 0 {
		return
	}

	cqlTime := p.space.lastPublished.Time.FormatTimeForCQLSearch(p.cqlTimeOffset())

	m, err := p.apiClient.IsSpaceModifiedSinceLastPublished(p.space.key, cqlTime)
	if err != nil {
		p.printFailureMessage(err)
		return
	}

	p.space.modifiedSinceLastPublished = m
}

func (p *Publisher) cqlTimeOffset() int {
	// nolint:gomnd
	minOffset := -12 // the latest time zone on earth, 12 hours behind UTC
	maxOffset := 14  // the earliest time zone on earth, 14 hours ahead of UTC

	for o := minOffset; o <= maxOffset; o++ {
		cqlTime := p.space.homepageCreated.FormatTimeForCQLSearch(o)
		pages := p.apiClient.PagesCreatedAt(cqlTime)

		for _, pg := range pages {
			if pg == p.space.homepageID {
				return o
			}
		}
	}

	panic("Could not calculate the time offset")
}

func (p *Publisher) printFailureMessage(s interface{}) {
	fmt.Printf("Failed: %v", s)
}

func (p *Publisher) publishAllSpecsUnder(baseSpecPath string) (err error) {
	return filepath.WalkDir(baseSpecPath, p.publishIfDirOrSpec)
}

func (p *Publisher) publishIfDirOrSpec(path string, d fs.DirEntry, err error) error {
	var e error

	entry := gauge.NewDirEntry(path, d)

	if entry.IsDirOrSpec() {
		e = p.publishDirOrSpec(entry)
	}

	if errors.IsNonfatal(e) {
		fmt.Printf("Skipping file: %v", e)
		return nil
	}

	return e
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

func (p *Publisher) updateLastPublished() error {
	return p.apiClient.UpdateLastPublished(p.space.homepageID, p.space.lastPublished.Version)
}
