//go:build integration

package integration

import (
	"net"
	"net/url"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	flaresolverr "github.com/astrocode-id/go-flaresolverr/v2"
)

// ipv4Pattern matches an IPv4 address anywhere in a page's text, regardless
// of the surrounding HTML/JSON structure.
var ipv4Pattern = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)

// flaresolverrURLEnv points the test at a running FlareSolverr instance.
// It's provided by CI as a service container and, for local runs, is
// expected to be started separately (e.g. `docker run -p 8191:8191
// ghcr.io/flaresolverr/flaresolverr:v3.5.0`).
const flaresolverrURLEnv = "FLARESOLVERR_URL"

func TestFlareSolverr(t *testing.T) {
	baseURL := os.Getenv(flaresolverrURLEnv)
	if baseURL == "" {
		baseURL = "http://localhost:8191/v1"
	}

	c, err := flaresolverr.NewClient(flaresolverr.Config{
		BaseURL: baseURL,
	})
	require.NoError(t, err)

	// httpbin.org/ip is a well-known, long-standing service that just
	// echoes back the caller's IP, so it's a stable target for testing
	// request.get.
	getResp, err := c.Get("https://httpbin.org/ip")
	require.NoError(t, err)

	ipAddress := ipv4Pattern.FindString(getResp.Solution.Response)
	assert.NotEmpty(t, ipAddress)
	assert.NotNil(t, net.ParseIP(ipAddress))

	// httpbin.org/post is a well-known, long-standing service that just
	// echoes back the posted form data, so it's a stable target for testing
	// request.post.
	postResp, err := c.Post("https://httpbin.org/post", url.Values{
		"key1": {"valueA"},
		"key2": {"valueB"},
	})
	require.NoError(t, err)
	assert.Contains(t, postResp.Solution.Response, "valueA")
	assert.Contains(t, postResp.Solution.Response, "valueB")
}
