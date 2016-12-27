package kitsu

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://kitsu.io/api/17/"
	edgeBaseURL    = "https://kitsu.io/api/edge/"

	defaultMediaType = "application/vnd.api+json"
)

// Client manages communication with the kitsu.io API.
type Client struct {
	client *http.Client

	BaseURL *url.URL

	common service
}

// NewClient returns a new kitsu.io API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(edgeBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL}

	c.common.client = c

	return c
}

type service struct {
	client *Client
}
