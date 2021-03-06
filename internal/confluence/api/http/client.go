// Package http provides a low-level Confluence API HTTP client.
package http

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a low-level Confluence API HTTP client.
type Client struct {
	restEndpoint string
	username     string
	token        string
	httpClient   *http.Client
}

// NewClient initialises a new low-level Confluence API HTTP client.
func NewClient(restEndpoint string, username string, token string) Client {
	httpClient := &http.Client{
		Timeout: time.Second * 10, //nolint: gomnd
	}

	return Client{restEndpoint, username, token, httpClient}
}

// GetJSON sends an HTTP GET request and stores the JSON response in the value pointed to by v.
func (c *Client) GetJSON(path string, v interface{}) error {
	responseBody, err := c.httpGet(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBody, &v)
}

// PutJSON sends an HTTP PUT request with the given URL path and request body
func (c *Client) PutJSON(path string, requestBody interface{}) error {
	b, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	return c.httpPut(path, b)
}

// PostJSON sends an HTTP POST request with the given URL path and request body
func (c *Client) PostJSON(path string, requestBody interface{}) error {
	b, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	return c.httpPost(path, b)
}

// DeleteContent deletes the Confluence content with the given contentID
// See e.g. https://developer.atlassian.com/server/confluence/confluence-rest-api-examples/#delete-a-page
func (c *Client) DeleteContent(contentID string) error {
	path := fmt.Sprintf("content/%s", contentID)
	return c.httpDelete(path)
}

func (c *Client) basicAuth() string {
	auth := c.username + ":" + c.token
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Client) httpGet(path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.restEndpoint, path)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+c.basicAuth())

	return c.do(req)
}

func (c *Client) httpPut(path string, requestBody []byte) error {
	url := fmt.Sprintf("%s/%s", c.restEndpoint, path)
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", "Basic "+c.basicAuth())
	_, err := c.do(req)

	return err
}

func (c *Client) httpPost(path string, requestBody []byte) error {
	url := fmt.Sprintf("%s/%s", c.restEndpoint, path)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", "Basic "+c.basicAuth())
	_, err := c.do(req)

	return err
}

func (c *Client) httpDelete(path string) error {
	url := fmt.Sprintf("%s/%s", c.restEndpoint, path)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Basic "+c.basicAuth())
	_, err := c.do(req)

	return err
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.Body != nil {
		defer response.Body.Close() //nolint: errcheck
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode > 299 { //nolint: gomnd
		err = newRequestError(response.StatusCode, string(body))
	}

	return body, err
}
