//go:build integration

package integration

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	"github.com/astrocode-id/go-flaresolverr"
)

// ipv4Pattern matches an IPv4 address anywhere in a page's text, regardless
// of the surrounding HTML/JSON structure.
var ipv4Pattern = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)

const (
	containerName = "flaresolverr"
	// flaresolverrVersion is pinned (instead of "latest") so the pulled image
	// is reproducible and cacheable in CI.
	// https://github.com/FlareSolverr/FlareSolverr/releases
	flaresolverrVersion = "v3.5.0"
)

func TestFlareSolverr(t *testing.T) {
	pool, err := dockertest.NewPool("")
	assert.NoError(t, err)

	_ = pool.RemoveContainerByName(containerName)
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       containerName,
		Repository: "ghcr.io/flaresolverr/flaresolverr",
		Tag:        flaresolverrVersion,
		Env: []string{
			"LOG_LEVEL=debug",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	assert.NoError(t, err)

	baseURL := fmt.Sprintf("http://%s/v1", resource.GetHostPort("8191/tcp"))

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

	err = pool.Purge(resource)
	assert.NoError(t, err)
}
