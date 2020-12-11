package affise

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/schema"
)

// UserAgent is the value for the library part of the User-Agent header
// that is sent with each request.
const UserAgent = "go-affise/" + Version

// Client is a client for the Affise API
type Client struct {
	apiKey     string
	endpoint   string
	userAgent  string
	encoder    *schema.Encoder
	httpClient *http.Client

	Offers  *OffersClient
	Presets *PresetsService
}

// A ClientOption is used to configure a Client.
type ClientOption func(*Client)

// WithEndpoint configures a Client to use the specified API endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(client *Client) {
		client.endpoint = strings.TrimRight(endpoint, "/")
	}
}

// WithAPIKey configures a Client to use the specified api key for authentication.
func WithAPIKey(token string) ClientOption {
	return func(client *Client) {
		client.apiKey = token
	}
}

// NewClient creates a new client.
func NewClient(options ...ClientOption) *Client {
	client := &Client{
		encoder:    schema.NewEncoder(),
		userAgent:  UserAgent,
		httpClient: &http.Client{},
	}
	client.encoder.SetAliasTag("json")

	for _, option := range options {
		option(client)
	}

	client.Offers = &OffersClient{client: client}
	client.Presets = &PresetsService{client: client}

	return client
}

// NewRequest creates an HTTP request against the API. The returned request
// is assigned with ctx and has all necessary headers
func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	url := c.endpoint + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("API-Key", c.apiKey)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req = req.WithContext(ctx)
	return req, nil
}

// NewRequest creates an HTTP request against the API with query parameters.
func (c *Client) NewRequestQueryParams(ctx context.Context, method, path string,
	queryParams interface{}, body io.Reader) (*http.Request, error) {

	if queryParams != nil {
		values := make(url.Values)
		err := c.encoder.Encode(queryParams, values)
		if err != nil {
			return nil, fmt.Errorf("encode query params err: %v", err)
		}
		path += "?" + values.Encode()
	}
	return c.NewRequest(ctx, method, path, body)
}

// Do performs an HTTP request against the API.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{Response: resp}

	err = c.checkResponse(resp)
	if err != nil {
		return response, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, bytes.NewReader(body))
		} else {
			err = json.Unmarshal(body, v)
		}
	}

	return response, err
}

func (c *Client) checkResponse(r *http.Response) error {
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("%v %v: %s", r.Request.Method, r.Request.URL, r.Status)
	}
	return nil
}

// Response represents a response from the API. It embeds http.Response.
type Response struct {
	*http.Response
}

//todo
// Pagination represents pagination information.
type Pagination struct {
	Page         int
	PerPage      int
	PreviousPage int
	NextPage     int
	LastPage     int
	TotalEntries int
}
