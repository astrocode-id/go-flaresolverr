package flaresolverr

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "http://localhost:8191/v1"

	contentApplicationJSON = "application/json"
)

// Config holds parameters for calling Flaresolverr NewClient.
type Config struct {
	// BaseURL is for Flaresolverr URL. Default: http://localhost:8191/v1.
	BaseURL string
	// Global Timeout to solve the challenge in milliseconds. Default: No timeout.
	Timeout int
}

// Client is a simple wrapper around the general Flaresolverr and represents
// a client to talk with Flaresolverr API.
type Client struct {
	baseURL string
	timeout int
}

// NewClient is the constructor for Flaresolverr API Client.
func NewClient(c Config) (*Client, error) {
	var baseURL string
	switch c.BaseURL {
	case "":
		baseURL = defaultBaseURL
	default:
		if _, err := url.Parse(c.BaseURL); err != nil {
			return nil, err
		}

		baseURL = c.BaseURL
	}

	return &Client{
		baseURL: baseURL,
		timeout: c.Timeout,
	}, nil
}

type command string

const (
	get  command = "request.get"
	post command = "request.post"
)

type status string

const (
	responseOK    status = "ok"
	responseError status = "error"
)

// Response holds raw response from Flaresolverr API.
type Response struct {
	Status         status   `json:"status,omitempty"`
	Message        string   `json:"message,omitempty"`
	Solution       Solution `json:"solution"`
	StartTimestamp int64    `json:"startTimestamp"`
	EndTimestamp   int64    `json:"endTimestamp"`
	Version        string   `json:"version"`
}

// Solution holds scraped web page from Flaresolverr API.
type Solution struct {
	URL       string  `json:"url"`
	Status    int     `json:"status"`
	Cookies   Cookies `json:"cookies"`
	UserAgent string  `json:"userAgent"`
	Response  string  `json:"response"`
}

// Cookie represents a single browser cookie, as sent in a request or
// returned in a Solution. For more detail, refer to
// https://github.com/FlareSolverr/FlareSolverr#-requestget.
type Cookie struct {
	Name     string  `json:"name,omitempty"`
	Value    string  `json:"value,omitempty"`
	Path     string  `json:"path,omitempty"`
	Domain   string  `json:"domain,omitempty"`
	Expires  float64 `json:"expires,omitempty"`
	HTTPOnly bool    `json:"httpOnly,omitempty"`
	Secure   bool    `json:"secure,omitempty"`
	SameSite string  `json:"sameSite,omitempty"`
}

// Cookies is a helper to manage cookies from the Flaresolverr API.
type Cookies []Cookie

// requestParams holds the base parameters supported by both the request.get
// and request.post commands.
// For more detail, refer to https://github.com/FlareSolverr/FlareSolverr#-requestget.
type requestParams struct {
	Cmd               command `json:"cmd"`
	URL               string  `json:"url"`
	PostData          string  `json:"postData,omitempty"`
	MaxTimeout        int     `json:"maxTimeout,omitempty"`
	Cookies           Cookies `json:"cookies,omitempty"`
	ReturnOnlyCookies bool    `json:"returnOnlyCookies,omitempty"`
	WaitInSeconds     int     `json:"waitInSeconds,omitempty"`
	DisableMedia      bool    `json:"disableMedia,omitempty"`
}

// Option configures an optional Get/Post request parameter.
// For more detail, refer to https://github.com/FlareSolverr/FlareSolverr#-requestget.
type Option func(*requestParams)

// WithMaxTimeout sets the timeout to solve the challenge in milliseconds for
// the current request. It replaces the global timeout from Config.Timeout.
func WithMaxTimeout(ms int) Option {
	return func(p *requestParams) {
		p.MaxTimeout = ms
	}
}

// WithCookies sets the browser cookies to send with the request.
func WithCookies(cookies Cookies) Option {
	return func(p *requestParams) {
		p.Cookies = cookies
	}
}

// WithReturnOnlyCookies strips the page content from the response, returning
// only the resulting cookies.
func WithReturnOnlyCookies(v bool) Option {
	return func(p *requestParams) {
		p.ReturnOnlyCookies = v
	}
}

// WithWaitInSeconds delays the response by the given number of seconds, to
// allow dynamic page content to finish loading.
func WithWaitInSeconds(s int) Option {
	return func(p *requestParams) {
		p.WaitInSeconds = s
	}
}

// WithDisableMedia prevents images, CSS, and fonts from loading.
func WithDisableMedia(v bool) Option {
	return func(p *requestParams) {
		p.DisableMedia = v
	}
}

// Get requests a web page with the request.get command and returns the
// whole Response.
// For more detail, refer to https://github.com/FlareSolverr/FlareSolverr#-requestget.
//
// Sessions and proxies are not yet supported by this client; see the README
// for details.
func (c *Client) Get(u string, opts ...Option) (Response, error) {
	p := requestParams{
		Cmd:        get,
		URL:        u,
		MaxTimeout: c.timeout,
	}
	for _, opt := range opts {
		opt(&p)
	}

	b, err := json.Marshal(p)
	if err != nil {
		return Response{}, err
	}

	return c.requestURL(b)
}

// Post requests a web page with the request.post command and returns the
// whole Response.
// For more detail, refer to https://github.com/FlareSolverr/FlareSolverr#-requestpost.
//
// Sessions and proxies are not yet supported by this client; see the README
// for details.
func (c *Client) Post(u string, postData url.Values, opts ...Option) (Response, error) {
	p := requestParams{
		Cmd:        post,
		URL:        u,
		PostData:   postData.Encode(),
		MaxTimeout: c.timeout,
	}
	for _, opt := range opts {
		opt(&p)
	}

	b, err := json.Marshal(p)
	if err != nil {
		return Response{}, err
	}

	return c.requestURL(b)
}

func (c *Client) requestURL(cmd []byte) (Response, error) {
	client := new(http.Client)
	r, err := client.Post(c.baseURL, contentApplicationJSON, bytes.NewReader(cmd))
	if err != nil {
		return Response{}, err
	}

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return Response{}, err
	}

	var resp Response
	if err := json.Unmarshal(b, &resp); err != nil {
		return Response{}, err
	}

	return resp, nil
}
