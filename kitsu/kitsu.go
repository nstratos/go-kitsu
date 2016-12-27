package kitsu

import (
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://kitsu.io/"

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
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL}

	c.common.client = c

	return c
}

type service struct {
	client *Client
}

// NewRequest creates an API request. If a relative URL is provided in urlStr,
// it will be resolved relative to the BaseURL of the Client. Relative URLs
// should always be specified without a preceding slash.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
