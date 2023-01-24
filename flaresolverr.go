package flaresolverr

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
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
	URL       string          `json:"url"`
	Status    int             `json:"status"`
	Cookies   Cookies         `json:"cookies"`
	UserAgent string          `json:"userAgent"`
	Response  json.RawMessage `json:"response"`
}

// Cookies is a helper to manage cookie from Flaresolverr API.
type Cookies []http.Cookie

type cookie struct {
	Name     string `json:"name,omitempty"`
	Value    string `json:"value,omitempty"`
	Path     string `json:"path,omitempty"`
	Domain   string `json:"domain,omitempty"`
	Expiry   int64  `json:"expiry,omitempty"`
	HTTPOnly bool   `json:"httpOnly,omitempty"`
	Secure   bool   `json:"secure,omitempty"`
	SameSite string `json:"sameSite,omitempty"`
}

// MarshalJSON is a custom function for marshalling Cookies.
func (c Cookies) MarshalJSON() ([]byte, error) {
	if len(c) == 0 {
		return []byte(`""`), nil
	}

	cookies := make([]cookie, 0, len(c))
	for _, cs := range c {
		var sameSite string
		switch cs.SameSite {
		case http.SameSiteStrictMode:
			sameSite = "Strict"
		case http.SameSiteNoneMode:
			sameSite = "None"
		case http.SameSiteLaxMode:
			sameSite = "Lax"
		}
		cookies = append(cookies, cookie{
			Name:     cs.Name,
			Value:    cs.Value,
			Path:     cs.Path,
			Domain:   cs.Domain,
			Expiry:   cs.Expires.Unix(),
			HTTPOnly: cs.HttpOnly,
			Secure:   cs.Secure,
			SameSite: sameSite,
		})
	}

	return json.Marshal(cookies)
}

// UnmarshalJSON is a custom function for UnmarshalJSON Cookies.
func (c *Cookies) UnmarshalJSON(b []byte) error {
	var cookies []cookie
	err := json.Unmarshal(b, &cookies)
	if err != nil {
		return err
	}

	*c = make(Cookies, 0, len(cookies))
	for _, cs := range cookies {
		t := time.Unix(cs.Expiry, 0)
		*c = append(*c, http.Cookie{
			Name:     cs.Name,
			Value:    cs.Value,
			Path:     cs.Path,
			Domain:   cs.Domain,
			Expires:  t,
			Secure:   cs.Secure,
			HttpOnly: cs.HTTPOnly,
			SameSite: 0,
		})
	}

	return nil
}

// GetParams holds parameters for calling Get.
type GetParams struct {
	URL string
	// MaxTimeout to solve the challenge in milliseconds for current API
	// It replaces global timeout from Config.Timeout. Default: No timeout.
	MaxTimeout        int
	Cookies           Cookies
	ReturnOnlyCookies bool
}

// Get requests web page with method http.Get and returns Solution.Response as raw bytes.
// For more detail, refer to https://github.com/FlareSolverr/FlareSolverr#-requestget.
func (c *Client) Get(p GetParams) ([]byte, error) {
	r, err := c.GetRaw(p)
	if err != nil {
		return nil, err
	}

	if r.Status != responseOK {
		return nil, errors.New(r.Message)
	}

	return r.Solution.Response, nil
}

// GetRaw requests web page with method http.Get and returns whole Response.
// For more detail, refer to https://github.com/FlareSolverr/FlareSolverr#-requestget.
func (c *Client) GetRaw(p GetParams) (Response, error) {
	var timeout int
	switch {
	case p.MaxTimeout > 0:
		timeout = p.MaxTimeout
	case c.timeout > 0:
		timeout = c.timeout
	}
	b, err := json.Marshal(requestParams{
		Cmd:               get,
		URL:               p.URL,
		MaxTimeout:        timeout,
		Cookies:           p.Cookies,
		ReturnOnlyCookies: p.ReturnOnlyCookies,
	})
	if err != nil {
		return Response{}, err
	}

	return c.requestURL(b)
}

// PostParams holds parameters for calling Post.
type PostParams struct {
	URL               string
	PostData          url.Values
	MaxTimeout        int
	Cookies           Cookies
	ReturnOnlyCookies bool
}

// Post requests web page with method http.Post and returns Solution.Response as raw bytes.
// For more detail, refer to https://github.com/FlareSolverr/FlareSolverr#-requestpost.
func (c *Client) Post(p PostParams) ([]byte, error) {
	r, err := c.PostRaw(p)
	if err != nil {
		return nil, err
	}

	return r.Solution.Response, nil
}

// PostRaw requests web page with method http.Post and returns whole Response.
// For more detail, refer to https://github.com/FlareSolverr/FlareSolverr#-requestpost.
func (c *Client) PostRaw(p PostParams) (Response, error) {
	var timeout int
	switch {
	case p.MaxTimeout > 0:
		timeout = p.MaxTimeout
	case c.timeout > 0:
		timeout = c.timeout
	}
	b, err := json.Marshal(requestParams{
		Cmd:               post,
		URL:               p.URL,
		PostData:          p.PostData.Encode(),
		MaxTimeout:        timeout,
		Cookies:           p.Cookies,
		ReturnOnlyCookies: p.ReturnOnlyCookies,
	})
	if err != nil {
		return Response{}, err
	}

	return c.requestURL(b)
}

type requestParams struct {
	Cmd               command `json:"cmd"`
	URL               string  `json:"url"`
	PostData          string  `json:"postData,omitempty"`
	MaxTimeout        int     `json:"maxTimeout,omitempty"`
	Cookies           Cookies `json:"cookies,omitempty"`
	ReturnOnlyCookies bool    `json:"returnOnlyCookies,omitempty"`
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
