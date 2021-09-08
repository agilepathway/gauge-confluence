// Package api provides a Confluence API client.
package api

import (
	"fmt"
	"strings"

	"github.com/agilepathway/gauge-confluence/internal/confluence/api/http"
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

// DeleteAllPagesInSpaceExceptHomepage deletes all the pages in the given Space,
// apart from the Space home page
func (c *Client) DeleteAllPagesInSpaceExceptHomepage(spaceKey string, homepageID string) (err error) {
	pageLimit := 10
	res, err := c.goconfluenceClient.GetContent(goconfluence.ContentQuery{
		SpaceKey: spaceKey,
		Limit:    pageLimit,
	})

	if err != nil {
		return err
	}

	for _, page := range res.Results {
		if page.ID != homepageID && page.Type == "page" {
			err = c.DeletePage(page.ID)

			if err != nil {
				return err
			}
		}
	}

	// the results are paginated so call the method again recursively until all pages are deleted
	if res.Size == res.Limit {
		return c.DeleteAllPagesInSpaceExceptHomepage(spaceKey, homepageID)
	}

	return nil
}

// DeletePage deletes a page from Confluence
func (c *Client) DeletePage(pageID string) (err error) {
	return c.httpClient.DeleteContent(pageID)
}

// CreateSpace creates a Confluence Space with the given key
func (c *Client) CreateSpace(key, name, description string) error {
	type Plain struct {
		Value          string `json:"value"`
		Representation string `json:"representation"`
	}

	type Description struct {
		Plain `json:"plain"`
	}

	type Data struct {
		Key         string `json:"key"`
		Name        string `json:"name"`
		Description `json:"description"`
	}

	path := "space"
	requestBody := Data{
		key,
		name,
		Description{
			Plain{
				description,
				"plain",
			},
		},
	}

	return c.httpClient.PostJSON(path, requestBody)
}

// DoesSpaceExist indicates if the Confluence Space with the given key exists.
func (c *Client) DoesSpaceExist(spaceKey string) (bool, error) {
	path := fmt.Sprintf("space/%s", spaceKey)

	var emptyStruct struct{}

	err := c.httpClient.GetJSON(path, &emptyStruct)

	if err != nil {
		e, ok := err.(*http.RequestError)
		if ok && e.StatusCode == 404 { //nolint:gomnd
			return false, nil
		}

		return false, err
	}

	return true, nil
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

// WasPageCreatedAt returns true if the page was created at the given time.
func (c *Client) WasPageCreatedAt(cqlTime string, pageID string) bool {
	query := goconfluence.SearchQuery{
		CQL: fmt.Sprintf("created=\"%s\" AND ID=\"%s\"", cqlTime, pageID),
	}
	result, _ := c.goconfluenceClient.Search(query)

	for _, r := range result.Results {
		if pageID == r.Content.ID {
			return true
		}
	}

	return false
}

// TotalPagesInSpace returns the number of pages (and blogposts) in the given Space
func (c *Client) TotalPagesInSpace(spaceKey string) (int, error) {
	query := goconfluence.SearchQuery{
		CQL: fmt.Sprintf("space=\"%s\" AND type IN (page, blogpost)", spaceKey),
	}
	result, err := c.goconfluenceClient.Search(query)

	if err != nil {
		return 0, err
	}

	return result.TotalSize, nil
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

// Version represents a Confluence version
type Version struct {
	Number    int  `json:"number"`
	MinorEdit bool `json:"minorEdit"`
}

// Data represents a last published request
type Data struct {
	Value   interface{} `json:"value"`
	Version Version     `json:"version"`
}

// SetContentProperty sets a content property with the provided key for the page with the provided ID.
func (c *Client) SetContentProperty(pageID string, propertyKey string, value interface{}, version int) error {
	path := fmt.Sprintf("content/%s/property/%s", pageID, propertyKey)
	requestBody := Data{
		value,
		Version{
			Number:    version,
			MinorEdit: true,
		},
	}

	return c.httpClient.PutJSON(path, requestBody)
}

// LastPublished returns the last time Confluence specs were published for the space with the given homepageID
func (c *Client) LastPublished(spaceHomepageID, propertyKey string) (string, int, error) {
	path := fmt.Sprintf("content/%s/property/%s", spaceHomepageID, propertyKey)

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
		// a 404 error means the lastPublished property has never been set, so return 0 as the version number
		if strings.Contains(err.Error(), "404") {
			return "", 0, nil
		}

		return "", 0, err
	}

	return resp.Value.LastPublished, resp.Version.Number, nil
}
