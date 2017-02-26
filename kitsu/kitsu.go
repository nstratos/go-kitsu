package kitsu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultBaseURL    = "https://kitsu.io/"
	defaultAPIVersion = "api/edge/"

	defaultMediaType = "application/vnd.api+json"
)

// Client manages communication with the kitsu.io API.
type Client struct {
	client *http.Client

	BaseURL *url.URL

	common service

	Anime *AnimeService
	User  *UserService
}

type service struct {
	client *Client
}

// Resource represent a JSON API resource object. It contains common fields
// used by the Kitsu API resources like Anime and Manga.
//
// JSON API docs: http://jsonapi.org/format/#document-resource-objects
type Resource struct {
	ID    string `json:"id"`
	Type  string `json:"type,omitempty"`
	Links Link   `json:"links,omitempty"`
}

// Link represent links that may be contained by resource objects. According to
// the current Kitsu API documentation, links are represented as a string.
//
// JSON API docs: http://jsonapi.org/format/#document-links
type Link struct {
	Self string `json:"self"`
}

// Options specifies the optional parameters to various List methods that
// support them.
//
// Pagination
//
// PageLimit and PageOffset provide pagination support. If PageLimit is not
// specified they are both ignored.
//
// Filtering
//
// If Filter is specified (e.g. genres) then one or more filter values can be
// passed in FilterVal (e.g. sports, sci-fi etc).
//
// Sorting
//
// Sort can be specified to provide sorting for one or more attributes (e.g.
// averageRating for Anime). By default, sorts are applied in ascending order.
// For descending order you can prepend a - to the sort parameter (e.g.
// -averageRating for Anime).
//
// Includes
//
// You can include one or more related resources by specifying the
// relationships in Include. You can also specify successive relationships
// using a . (e.g. media.genres for library entries).
type Options struct {
	PageLimit  int
	PageOffset int
	Filter     string
	FilterVal  []string
	Sort       []string
	Include    []string
}

func addOptions(s string, opt *Options) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	v := u.Query()
	if opt.PageLimit != 0 {
		v.Set("page[limit]", strconv.Itoa(opt.PageLimit))
		v.Set("page[offset]", strconv.Itoa(opt.PageOffset))
	}
	if opt.Filter != "" && opt.FilterVal != nil {
		v.Set(fmt.Sprintf("filter[%s]", opt.Filter),
			strings.Join(opt.FilterVal, ","))
	}
	if opt.Sort != nil {
		v.Set("sort", strings.Join(opt.Sort, ","))
	}
	if opt.Sort != nil {
		v.Set("include", strings.Join(opt.Include, ","))
	}

	u.RawQuery = v.Encode()
	return u.String(), nil
}

// NewClient returns a new kitsu.io API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL}

	c.common.client = c

	c.Anime = (*AnimeService)(&c.common)
	c.User = (*UserService)(&c.common)

	return c
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
		encErr := json.NewEncoder(buf).Encode(body)
		if encErr != nil {
			return nil, encErr
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
	req.Header.Set("Accept", defaultMediaType)

	return req, nil
}

// Response is a Kitsu API response. It wraps the standard http.Response
// returned from the request and provides access to pagination offsets for
// responses that return an array of results.
type Response struct {
	*http.Response

	NextOffset  int
	PrevOffset  int
	FirstOffset int
	LastOffset  int
}

func newResponse(r *http.Response) *Response {
	return &Response{Response: r}
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
