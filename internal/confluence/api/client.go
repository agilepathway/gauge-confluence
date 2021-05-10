// Package api provides a Confluence API client.
package api

import (
	"fmt"
	"strings"

	"github.com/agilepathway/gauge-confluence/internal/confluence/api/http"
	"github.com/agilepathway/gauge-confluence/internal/confluence/time"
	"github.com/agilepathway/gauge-confluence/internal/env"
	"github.com/agilepathway/gauge-confluence/util"

	goconfluence "github.com/virtomize/confluence-go-api"
)

// Client is a Confluence API client.
type Client struct {
	httpClient         http.Client
	goconfluenceClient *goconfluence.API
}

// NewClient initialises a new Client.
func NewClient() Client {
	httpClient := http.NewClient(baseEndpoint(), username(), token())
	goconfluenceClient, err := goconfluence.NewAPI(baseEndpoint(), username(), token())
	util.Fatal("Error while creating Confluence API Client", err)

	return Client{httpClient, goconfluenceClient}
}

// PublishPage publishes a page to Confluence as a child of the given parent page.
func (c *Client) PublishPage(spaceKey, title, body, parentPageID string) (pageID string, err error) {
	requestContent := &goconfluence.Content{
		Type:  "page",
		Title: title,
		Body: goconfluence.Body{
			Storage: goconfluence.Storage{
				Value:          body,
				Representation: "wiki",
			},
		},
		Space:   goconfluence.Space{Key: spaceKey},
		Version: &goconfluence.Version{Number: 1},
	}

	if parentPageID != "" {
		requestContent.Ancestors = []goconfluence.Ancestor{
			{ID: parentPageID},
		}
	}

	responseContent, err := c.goconfluenceClient.CreateContent(requestContent)

	if err != nil {
		return "", err
	}

	return responseContent.ID, nil
}

// SpaceHomepage provides the page ID, no. of children and created time for the given Space's homepage.
func (c *Client) SpaceHomepage(spaceKey string) (string, int, string, error) {
	path := fmt.Sprintf("space?spaceKey=%s&expand=homepage.children.page,homepage.history", spaceKey)

	var homepageResponse struct {
		Results []struct {
			Homepage struct {
				ID      string `json:"id"`
				History struct {
					CreatedDate string `json:"createdDate"`
				} `json:"history"`
				Children struct {
					Page struct {
						Size int `json:"size"`
					} `json:"page"`
				} `json:"children"`
			} `json:"homepage"`
		} `json:"results"`
	}

	err := c.httpClient.GetJSON(path, &homepageResponse)

	if err != nil {
		return "", 0, "", err
	}

	h := homepageResponse.Results[0].Homepage

	return h.ID, h.Children.Page.Size, h.History.CreatedDate, nil
}

// IsSpaceModifiedSinceLastPublished returns true if any page was modified since the last publish
//
// The lastPublished parameter is a string in Confluence CQL format.
func (c *Client) IsSpaceModifiedSinceLastPublished(spaceKey string, lastPublished string) (bool, error) {
	query := goconfluence.SearchQuery{
		CQL: fmt.Sprintf("space.key=\"%s\" and lastModified>\"%s\"", spaceKey, lastPublished),
	}
	result, err := c.goconfluenceClient.Search(query)

	if err != nil {
		return true, err
	}

	return result.TotalSize > 0, nil
}

// PagesCreatedAt returns the pageIDs for pages created at the given cqlTime.
func (c *Client) PagesCreatedAt(cqlTime string) []string {
	query := goconfluence.SearchQuery{
		CQL: fmt.Sprintf("created=\"%s\"", cqlTime),
	}
	result, _ := c.goconfluenceClient.Search(query)

	pages := make([]string, len(result.Results))

	for _, r := range result.Results {
		pages = append(pages, r.Content.ID)
	}

	return pages
}

func baseEndpoint() string {
	return fmt.Sprintf("%s/rest/api", baseURL())
}

func username() string {
	return env.GetRequired("CONFLUENCE_USERNAME")
}

func token() string {
	return env.GetRequired("CONFLUENCE_TOKEN")
}

func baseURL() string {
	confluenceBaseURL := strings.TrimSuffix(env.GetRequired("CONFLUENCE_BASE_URL"), "/")
	if strings.HasSuffix(confluenceBaseURL, "atlassian.net") {
		return fmt.Sprintf("%s/wiki", confluenceBaseURL)
	}

	return confluenceBaseURL
}

// Value contains the LastPublished time
type Value struct {
	LastPublished string `json:"lastPublished"`
}

// Version represents a Confluence version
type Version struct {
	Number    int  `json:"number"`
	MinorEdit bool `json:"minorEdit"`
}

// Data represents a last published request
type Data struct {
	Value   Value   `json:"value"`
	Version Version `json:"version"`
}

// UpdateLastPublished stores the time of publishing as a Confluence content property,
// so that in the next run of the plugin it can check that the Confluence space has not
// been edited manually in the meantime.
//
// The content property is attached to the Space homepage rather than to the Space itself, as
// attaching the property to the Space requires admin permissions and we want to allow the
// plugin to be run by non-admin users too.
func (c *Client) UpdateLastPublished(homepageID string, currentVersion int) error {
	path := fmt.Sprintf("content/%s/property/%s", homepageID, time.LastPublishedPropertyKey)
	requestBody := Data{
		Value{
			LastPublished: time.Now().String(),
		},
		Version{
			Number:    currentVersion + 1,
			MinorEdit: true,
		},
	}

	return c.httpClient.PutJSON(path, requestBody)
}

// LastPublished returns the last time Confluence specs were published for the space with the given homepageID
func (c *Client) LastPublished(spaceHomepageID string) (time.LastPublished, error) {
	path := fmt.Sprintf("content/%s/property/%s", spaceHomepageID, time.LastPublishedPropertyKey)

	var resp struct {
		Value struct {
			LastPublished string `json:"lastPublished"`
		} `json:"value"`
		Version struct {
			Number int `json:"number"`
		} `json:"version"`
	}

	err := c.httpClient.GetJSON(path, &resp)

	if err != nil && !strings.Contains(err.Error(), "404") {
		return time.LastPublished{}, err
	}

	return time.NewLastPublished(resp.Value.LastPublished, resp.Version.Number), nil
}
