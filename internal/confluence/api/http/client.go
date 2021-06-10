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

func (c *Client) basicAuth() string {
	auth := c.username + ":" + c.token
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *Client) httpGet(path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.restEndpoint, path)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+c.basicAuth())

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
		err = fmt.Errorf("HTTP response error: %d %s", response.StatusCode, body)
	}

	return body, err
}

func (c *Client) httpPut(path string, requestBody []byte) error {
	url := fmt.Sprintf("%s/%s", c.restEndpoint, path)
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Authorization", "Basic "+c.basicAuth())

	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if response.Body != nil {
		defer response.Body.Close() //nolint: errcheck
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode > 299 { //nolint: gomnd
		err = fmt.Errorf("HTTP response error: %d %s", response.StatusCode, responseBody)
	}

	return err
}
