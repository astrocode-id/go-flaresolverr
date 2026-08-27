//go:build integration

package integration

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/astrocode-id/go-flaresolverr"
)

// ipv4Pattern matches an IPv4 address anywhere in a page's text, regardless
// of the surrounding HTML/JSON structure.
var ipv4Pattern = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)

// flaresolverrVersion is pinned (instead of "latest") so the pulled image
// is reproducible and cacheable in CI.
// https://github.com/FlareSolverr/FlareSolverr/releases
const flaresolverrVersion = "v3.5.0"

func TestFlareSolverr(t *testing.T) {
	ctx := context.Background()

	container, err := testcontainers.Run(ctx, "ghcr.io/flaresolverr/flaresolverr:"+flaresolverrVersion,
		testcontainers.WithEnv(map[string]string{
			"LOG_LEVEL": "debug",
		}),
		testcontainers.WithExposedPorts("8191/tcp"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("8191/tcp")),
	)
	assert.NoError(t, err)

	baseURL, err := container.PortEndpoint(ctx, "8191/tcp", "http")
	assert.NoError(t, err)
	baseURL += "/v1"

	<-time.After(5 * time.Second)

	fmt.Println("connect to FlareSolverr base URL: " + baseURL)
	c, err := flaresolverr.NewClient(flaresolverr.Config{
		BaseURL: baseURL,
	})
	assert.NoError(t, err)

	// httpbin.org/ip is a well-known, long-standing service that just
	// echoes back the caller's IP, so it's a stable target for testing
	// request.get.
	getResp, err := c.Get("https://httpbin.org/ip")
	assert.NoError(t, err)

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
	assert.NoError(t, err)
	assert.Contains(t, postResp.Solution.Response, "valueA")
	assert.Contains(t, postResp.Solution.Response, "valueB")

	err = container.Terminate(ctx)
	assert.NoError(t, err)
}
