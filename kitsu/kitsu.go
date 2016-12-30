package kitsu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
// should always be specified without a preceding slash. If body is specified,
// it will be encoded to JSON and used as the request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-type", defaultMediaType)
	}

	return req, nil
}

// Do sends an API request and returns the API response. If an API error has
// occurred both the response and the error will be returned in case the caller
// wishes to further inspect the response. If v is passed as an argument, then
// the API response is JSON decoded and stored to v.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = checkResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

// ErrorResponse reports one or more errors caused by an API request.
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Errors   []Error        `json:"errors"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Errors)
}

// Error holds the details of each invidivual error in an ErrorResponse.
//
// JSON API docs: http://jsonapi.org/format/#error-objects
type Error struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Code   string `json:"code"`
	Status string `json:"status"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v: error %v: %v(%v)",
		e.Status, e.Code, e.Title, e.Detail)
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	body, err := ioutil.ReadAll(r.Body)
	if err == nil && body != nil {
		json.Unmarshal(body, errorResponse)
	}
	return errorResponse
}
