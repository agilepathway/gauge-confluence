// Package api provides a Confluence API client
package api

import (
	"fmt"
	"strings"

	"github.com/agilepathway/gauge-confluence/internal/confluence/api/http"
	"github.com/agilepathway/gauge-confluence/internal/env"
	"github.com/agilepathway/gauge-confluence/util"

	goconfluence "github.com/virtomize/confluence-go-api"
)

// Client is a Confluence API client
type Client struct {
	goconfluenceClient *goconfluence.API
	httpClient         http.Client
}

// NewClient initialises a new Client
func NewClient() Client {
	return Client{confluenceClient(), httpClient()}
}

// PublishPage publishes a page to Confluence as a child of the given parent page
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

// SpaceHomepageID retrieves the page ID for the given Space's homepage
func (c *Client) SpaceHomepageID(spaceKey string) (string, error) {
	path := fmt.Sprintf("space?spaceKey=%s&expand=homepage", spaceKey)

	var homepageResponse struct {
		Results []struct {
			Homepage struct {
				ID string `json:"id"`
			} `json:"homepage"`
		} `json:"results"`
	}

	err := c.httpClient.GetJSON(path, &homepageResponse)

	if err != nil {
		return "", err
	}

	return homepageResponse.Results[0].Homepage.ID, nil
}

func confluenceClient() *goconfluence.API {
	api, err := goconfluence.NewAPI(baseEndpoint(), username(), token())
	util.Fatal("Error while creating Confluence API Client", err)

	return api
}

func httpClient() http.Client {
	return http.NewClient(baseEndpoint(), username(), token())
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
