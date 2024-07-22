package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultPGVectorsCloudURL = "https://cloud.pgvecto.rs/api/v1"
)

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	apiKey     string
	apiUrl     string
	userAgent  string
	HttpClient HttpClient
}

var (
	errApiKeyRequired = fmt.Errorf("ApiKey is required")
)

func checkApiKey(c *Client) func() error {
	return func() error {
		if c.apiKey == "" {
			return errApiKeyRequired
		}
		return nil
	}
}

func checkApiUrl(c *Client) func() error {
	return func() error {
		if c.apiUrl == "" {
			return fmt.Errorf("ApiUrl is required")
		}
		return nil
	}
}

func (client *Client) Clone(opts ...Option) (*Client, error) {
	clone := func(c Client) *Client {
		return &c
	}

	c := clone(*client)

	for _, opt := range opts {
		opt(c)
	}

	applyOverride(c)

	if err := validate(c); err != nil {
		return nil, err
	}

	return c, nil
}

func validate(c *Client) error {
	checkFns := []func() error{
		checkApiKey(c),
		checkApiUrl(c),
	}
	for _, fn := range checkFns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func applyOverride(c *Client) {
	overrides := []Option{
		OverrideApiUrl(c.apiUrl),
	}
	for _, opt := range overrides {
		opt(c)
	}
}

func applyDefaults(c *Client) {
	defaultOptions := []Option{
		WithDefaultClient(),
		WithDefaultApiUrl(),
		WithDefaultUserAgent(),
	}
	for _, opt := range defaultOptions {
		opt(c)
	}
}

func NewClient(opts ...Option) (*Client, error) {

	// create a new client with options
	c := new(Client)
	for _, opt := range opts {
		opt(c)
	}

	applyDefaults(c)

	if err := validate(c); err != nil {
		return nil, err
	}

	return c, nil
}

type Option func(*Client)

func WithHTTPClient(client HttpClient) Option {
	return func(c *Client) {
		c.HttpClient = client
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

func WithApiKey(apiKey string) Option {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

func WithDefaultApiUrl() Option {
	return func(c *Client) {
		if c.apiKey == "" {
			c.apiUrl = defaultPGVectorsCloudURL
		}
	}
}

func OverrideApiUrl(apiUrl string) Option {
	return func(c *Client) {
		if apiUrl != "" {
			c.apiUrl = apiUrl
		}
	}
}

func WithDefaultClient() Option {
	return func(c *Client) {
		if c.HttpClient == nil {
			c.HttpClient = &http.Client{}
		}
	}
}

func WithDefaultUserAgent() Option {
	return func(c *Client) {
		if c.userAgent == "" {
			c.userAgent = "tensorchord/terraform-provider-pgvecto-rs-cloud"
		}
	}
}

func (c *Client) do(method string, path string, body interface{}, result interface{}) error {

	u, err := c.getURL(path)
	if err != nil {
		return err
	}
	req, err := c.newRequest(method, u, body)
	if err != nil {
		return err
	}
	return c.doRequest(req, result)
}

func (c *Client) newRequest(method string, u *url.URL, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request, v any) error {
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		return parseError(res.Body)
	}

	return decodeResponse(res.Body, v)
}

func parseError(body io.Reader) error {

	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	var e Error
	err = json.Unmarshal(b, &e)
	if err != nil {
		return err
	}

	return e
}

func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	var apierr Error
	err = json.Unmarshal(b, &apierr)
	if err == nil && apierr.HTTPStatusCode != 0 {
		return &apierr
	}
	err = json.Unmarshal(b, v)
	return err
}

func (c *Client) getURL(path string) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("%s/%s", c.apiUrl, path))
}
