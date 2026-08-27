//go:build integration

package integration

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	"github.com/astrocode-id/go-flaresolverr"
)

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

	b, err := c.Get(flaresolverr.GetParams{
		URL: "https://ifconfig.me",
	})
	assert.NoError(t, err)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	assert.NoError(t, err)

	ipAddress := doc.Find("strong").First().Text()
	// FlareSolverr returns the response body JSON-escaped but undecoded, so
	// literal `\n` (backslash+n) sequences show up in the text alongside real
	// whitespace; strip both before parsing.
	ipAddress = strings.Trim(ipAddress, "\\n \t\r\n")
	assert.NotNil(t, net.ParseIP(ipAddress))

	err = pool.Purge(resource)
	assert.NoError(t, err)
}
