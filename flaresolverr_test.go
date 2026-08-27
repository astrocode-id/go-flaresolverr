package flaresolverr

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc              string
		method            command
		url               string
		maxTimeout        int
		cookies           Cookies
		returnOnlyCookies bool
		handlerFunc       func(t *testing.T, p requestParams) http.HandlerFunc
		expected          Response
		isError           assert.ErrorAssertionFunc
	}{
		{
			desc:       "GET url with no error returns response",
			method:     get,
			url:        "https://try.me",
			maxTimeout: 5000,
			cookies: Cookies{
				{
					Name:    "OGPC",
					Value:   "19033459-1:",
					Path:    "/",
					Domain:  ".try.me",
					Expires: 1679759834,
				},
			},
			returnOnlyCookies: false,
			handlerFunc: func(t *testing.T, p requestParams) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					b, err := io.ReadAll(r.Body)
					defer r.Body.Close()
					assert.NoError(t, err)

					var test requestParams
					err = json.Unmarshal(b, &test)
					assert.NoError(t, err)
					assert.Equal(t, p, test)

					b, err = os.ReadFile("./testdata/response_get_success.json")
					assert.NoError(t, err)

					_, err = w.Write(b)
					assert.NoError(t, err)
				}
			},
			expected: Response{
				Status:  responseOK,
				Message: "Challenge not detected!",
				Solution: Solution{
					URL:    "https://try.me/",
					Status: 200,
					Cookies: Cookies{
						{
							Name:    "OGPC",
							Value:   "19033459-1:",
							Path:    "/",
							Domain:  ".try.me",
							Expires: 1679759834,
						},
						{
							Name:     "1P_JAR",
							Value:    "2023-01-24-15",
							Path:     "/",
							Domain:   ".try.me",
							Secure:   true,
							Expires:  1677167834,
							SameSite: "None",
						},
					},
					UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
					Response:  `<html lang="en">TRY_ME</html>`,
				},
				StartTimestamp: 1674575494857,
				EndTimestamp:   1674575499113,
				Version:        "3.0.2",
			},
			isError: assert.NoError,
		},
		{
			desc:       "GET url with error returns error",
			method:     get,
			url:        "https://try.me",
			maxTimeout: 5000,
			cookies: Cookies{
				{
					Name:    "OGPC",
					Value:   "19033459-1:",
					Path:    "/",
					Domain:  ".try.me",
					Expires: 1679759834,
				},
			},
			returnOnlyCookies: false,
			handlerFunc: func(t *testing.T, p requestParams) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					b, err := io.ReadAll(r.Body)
					defer r.Body.Close()
					assert.NoError(t, err)

					var test requestParams
					err = json.Unmarshal(b, &test)
					assert.NoError(t, err)
					assert.Equal(t, p, test)

					b, err = os.ReadFile("./testdata/response_error.json")
					assert.NoError(t, err)

					_, err = w.Write(b)
					assert.NoError(t, err)
				}
			},
			expected: Response{
				Status:         responseError,
				Message:        "Error: Not implemented yet.",
				Solution:       Solution{},
				StartTimestamp: 1674574599950,
				EndTimestamp:   1674574599952,
				Version:        "3.0.2",
			},
			isError: assert.NoError,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ts := httptest.NewServer(test.handlerFunc(t, requestParams{
				Cmd:               test.method,
				URL:               test.url,
				MaxTimeout:        test.maxTimeout,
				Cookies:           test.cookies,
				ReturnOnlyCookies: test.returnOnlyCookies,
			}))
			defer ts.Close()

			c, err := NewClient(Config{
				BaseURL: ts.URL,
				Timeout: 0,
			})
			assert.NoError(t, err)

			resp, err := c.Get(test.url,
				WithMaxTimeout(test.maxTimeout),
				WithCookies(test.cookies),
				WithReturnOnlyCookies(test.returnOnlyCookies),
			)
			test.isError(t, err)
			assert.Equal(t, test.expected, resp)
		})
	}
}

func TestClient_Post(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc              string
		method            command
		url               string
		postDate          url.Values
		maxTimeout        int
		cookies           Cookies
		returnOnlyCookies bool
		handlerFunc       func(t *testing.T, p requestParams) http.HandlerFunc
		expected          Response
		isError           assert.ErrorAssertionFunc
	}{
		{
			desc:   "POST url with no error returns response",
			method: post,
			url:    "https://try.me/form-post-tester.php",
			postDate: map[string][]string{
				"q": {"test1"},
				"v": {"test2"},
			},
			maxTimeout: 5000,
			cookies: Cookies{
				{
					Name:    "OGPC",
					Value:   "19033459-1:",
					Path:    "/",
					Domain:  ".try.me",
					Expires: 1679759834,
				},
			},
			returnOnlyCookies: false,
			handlerFunc: func(t *testing.T, p requestParams) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					b, err := io.ReadAll(r.Body)
					defer r.Body.Close()
					assert.NoError(t, err)

					var test requestParams
					err = json.Unmarshal(b, &test)
					assert.NoError(t, err)
					assert.Equal(t, p, test)

					b, err = os.ReadFile("./testdata/response_post_success.json")
					assert.NoError(t, err)

					_, err = w.Write(b)
					assert.NoError(t, err)
				}
			},
			expected: Response{
				Status:  responseOK,
				Message: "Challenge not detected!",
				Solution: Solution{
					URL:       "https://try.me/form-post-tester.php",
					Status:    200,
					Cookies:   Cookies{},
					UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
					Response: `<html><head><meta name="color-scheme" content="light dark"></head><body><pre style="word-wrap: break-word; white-space: pre-wrap;">The POSTed value is:
********************
key1=valueA&amp;key2=valueB
********************


The Headers are:
********************
Host: www.hashemian.com
Connection: Keep-Alive
Accept-Encoding: gzip
X-Forwarded-For: 45.77.173.90
CF-RAY: 78ea349188884cdd-EWR
Content-Length: 23
X-Forwarded-Proto: https
CF-Visitor: {"scheme":"https"}
cache-control: max-age=0
sec-ch-ua: "Not?A_Brand";v="8", "Chromium";v="108"
sec-ch-ua-mobile: ?0
sec-ch-ua-platform: "Linux"
upgrade-insecure-requests: 1
origin: null
content-type: application/x-www-form-urlencoded
user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36
accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
sec-fetch-site: cross-site
sec-fetch-mode: navigate
sec-fetch-dest: document
accept-language: en-US,en;q=0.9
CF-Connecting-IP: 45.77.173.90
CF-IPCountry: SG
CDN-Loop: cloudflare
********************

</pre></body></html>`,
				},
				StartTimestamp: 1674578364632,
				EndTimestamp:   1674578368933,
				Version:        "3.0.2",
			},
			isError: assert.NoError,
		},
		{
			desc:       "POST url with error returns error",
			method:     post,
			url:        "https://try.me/form-post-tester.php",
			maxTimeout: 5000,
			cookies: Cookies{
				{
					Name:    "OGPC",
					Value:   "19033459-1:",
					Path:    "/",
					Domain:  ".try.me",
					Expires: 1679759834,
				},
			},
			returnOnlyCookies: false,
			handlerFunc: func(t *testing.T, p requestParams) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					b, err := io.ReadAll(r.Body)
					defer r.Body.Close()
					assert.NoError(t, err)

					var test requestParams
					err = json.Unmarshal(b, &test)
					assert.NoError(t, err)
					assert.Equal(t, p, test)

					b, err = os.ReadFile("./testdata/response_error.json")
					assert.NoError(t, err)

					_, err = w.Write(b)
					assert.NoError(t, err)
				}
			},
			expected: Response{
				Status:         responseError,
				Message:        "Error: Not implemented yet.",
				Solution:       Solution{},
				StartTimestamp: 1674574599950,
				EndTimestamp:   1674574599952,
				Version:        "3.0.2",
			},
			isError: assert.NoError,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ts := httptest.NewServer(test.handlerFunc(t, requestParams{
				Cmd:               test.method,
				URL:               test.url,
				PostData:          test.postDate.Encode(),
				MaxTimeout:        test.maxTimeout,
				Cookies:           test.cookies,
				ReturnOnlyCookies: test.returnOnlyCookies,
			}))
			defer ts.Close()

			c, err := NewClient(Config{
				BaseURL: ts.URL,
				Timeout: 0,
			})
			assert.NoError(t, err)

			resp, err := c.Post(test.url, test.postDate,
				WithMaxTimeout(test.maxTimeout),
				WithCookies(test.cookies),
				WithReturnOnlyCookies(test.returnOnlyCookies),
			)
			test.isError(t, err)
			assert.Equal(t, test.expected, resp)
		})
	}
}
