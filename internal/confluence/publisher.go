// Package confluence publishes Gauge specifications to Confluence.
package confluence

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/agilepathway/gauge-confluence/gauge_messages"
	"github.com/agilepathway/gauge-confluence/internal/confluence/api"
	"github.com/agilepathway/gauge-confluence/internal/env"
	"github.com/agilepathway/gauge-confluence/internal/errors"
	"github.com/agilepathway/gauge-confluence/internal/gauge"
	"github.com/agilepathway/gauge-confluence/internal/git"
	"github.com/agilepathway/gauge-confluence/internal/logger"
)

// Publisher publishes Gauge specifications to Confluence.
type Publisher struct {
	apiClient   api.Client
	space       space           // Represents the Confluence Space that is published to
	specs       map[string]Spec // keyed by filepath
	dryRunPages map[string]page // Used to check for duplicate pages, keyed by filepath
}

// NewPublisher instantiates a new Publisher.
func NewPublisher(m *gauge_messages.SpecDetails) Publisher {
	apiClient := api.NewClient()

	return Publisher{apiClient: apiClient, space: newSpace(apiClient), specs: makeSpecsMap(m),
		dryRunPages: make(map[string]page)}
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
func (p *Publisher) Publish(specPaths []string) (err error) { //nolint:funlen
	logger.Infof(true, "Checking specs are in a valid state for publishing to Confluence ...")

	for _, specPath := range specPaths {
		err = p.dryRunChecks(specPath)
		if err != nil {
			return err
		}
	}

	if env.GetBoolean("DRY_RUN") {
		logger.Infof(true, "Dry run finished successfully")
		return nil
	}

	logger.Infof(true, "Checking finished successfully")
	logger.Infof(true, "Publishing Gauge specs to Confluence ...")

	p.space.setup()

	if p.space.err != nil {
		return p.space.err
	}

	err = p.space.deleteAllPagesExceptHomepage()

	if err != nil {
		return err
	}

	for _, specPath := range specPaths {
		err = p.publishAllSpecsUnder(specPath)
		if err != nil {
			break
		}
	}

	if err != nil {
		p.space.updateLastPublished() //nolint:errcheck,gosec
		return err
	}

	err = p.space.updateLastPublished()

	if err != nil {
		return err
	}

	spaceName := p.space.name()

	if p.space.err != nil {
		return p.space.err
	}

	err = p.space.homepage.publish()
	if err != nil {
		return err
	}

	logger.Infof(true, "Success: published %d specs and directory pages to Confluence Space named: %s",
		len(p.space.publishedPages), spaceName)

	return nil
}

func (p *Publisher) dryRunChecks(baseSpecPath string) (err error) {
	return filepath.WalkDir(baseSpecPath, p.dryRunCheck)
}

func (p *Publisher) publishAllSpecsUnder(baseSpecPath string) (err error) {
	err = filepath.WalkDir(baseSpecPath, p.publishIfDirOrSpec)
	if err != nil {
		return err
	}

	return p.space.deleteEmptyDirPages()
}

func (p *Publisher) dryRunCheck(path string, d fs.DirEntry, err error) error {
	entry := gauge.NewDirEntry(path, d)
	if entry.IsDirOrSpec() {
		pg, err := newPage(entry, "", p.specs[entry.Path])
		if err != nil {
			if errors.IsNonfatal(err) {
				fmt.Printf("Skipping file: %v", err)
				return nil
			}

			return err
		}

		err = p.checkForDuplicateTitle(pg)

		if err != nil {
			return err
		}

		p.dryRunPages[pg.path] = pg
	}

	return err
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
	publishedPageID, err := p.apiClient.CreatePage(p.space.key, pg.title, pg.body, pg.parentID)

	if err != nil {
		return err
	}

	pg.id = publishedPageID

	p.space.publishedPages[pg.path] = pg

	return nil
}

// checkForDuplicateTitle returns an error if the given page has the same title as an already published page.
func (p *Publisher) checkForDuplicateTitle(given page) error {
	for _, pg := range p.dryRunPages {
		if pg.title == given.title {
			return &duplicatePageError{pg, given}
		}
	}

	return nil
}
