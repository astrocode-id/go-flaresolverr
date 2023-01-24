//go:build integration

package integration

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	"github.com/fahimbagar/go-flaresolverr"
)

const (
	containerName = "flaresolverr"
)

var (
	isCircleCI bool
	resource   *dockertest.Resource
)

func init() {
	isCircleCI, _ = strconv.ParseBool(os.Getenv("CIRCLECI"))
}

func TestFlareSolverr(t *testing.T) {
	pool, err := dockertest.NewPool("")
	assert.NoError(t, err)

	var baseURL string

	if isCircleCI {
		baseURL = fmt.Sprintf("http://localhost:8191/v1")
	} else {
		_ = pool.RemoveContainerByName(containerName)
		resource, err = pool.RunWithOptions(&dockertest.RunOptions{
			Name:       containerName,
			Repository: "ghcr.io/flaresolverr/flaresolverr",
			Tag:        "latest",
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

		baseURL = fmt.Sprintf("http://%s/v1", resource.GetHostPort("8191/tcp"))

		<-time.After(5 * time.Second)
	}

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
	assert.NotNil(t, net.ParseIP(ipAddress))

	if !isCircleCI {
		err := pool.Purge(resource)
		assert.NoError(t, err)
	}
}
