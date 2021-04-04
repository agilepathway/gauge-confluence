// Package confluence publishes Gauge specifications to Confluence.
package confluence

import (
	"fmt"
)

// Specs is a collection of Gauge Specifications which can be published to Confluence
type Specs []Spec

// PublishToConfluence publishes Gauge specifications to Confluence
func (s Specs) PublishToConfluence() {
	var unpublishedSpecs []Spec

	for _, spec := range s {
		err := s.publishSpec(spec)
		if err != nil {
			unpublishedSpecs = append(unpublishedSpecs, spec)
			fmt.Printf("Failed to publish spec %s: %s\n", spec.path, err)
		}
	}

	switch len(s) - len(unpublishedSpecs) {
	case 0:
		fmt.Println("No specifications were found - so nothing to publish to Confluence")
	case 1:
		fmt.Println("Published 1 specification to Confluence")
	default:
		fmt.Printf("Published %d specifications to Confluence\n", len(s))
	}
}

func (s Specs) publishSpec(spec Spec) error {
	return nil
}
