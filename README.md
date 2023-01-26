# FlareSolverr v3 Go Client

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Code Climate](https://codeclimate.com/github/astrocode-id/go-flaresolverr.png)](https://codeclimate.com/github/astrocode-id/go-flaresolverr)
[![Test Coverage](https://api.codeclimate.com/v1/badges/c8eaaff0f761d4d1f09f/test_coverage)](https://codeclimate.com/github/astrocode-id/go-flaresolverr/test_coverage)
[![GitHub issues](https://img.shields.io/github/issues/astrocode-id/go-flaresolverr)](https://github.com/astrocode-id/go-flaresolverr/issues)
[![CircleCI](https://circleci.com/gh/astrocode-id/go-flaresolverr.svg?style=shield)](https://circleci.com/gh/astrocode-id/go-flaresolverr)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fastrocode-id%2Fgo-flaresolverr.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fastrocode-id%2Fgo-flaresolverr?ref=badge_shield)

[go-flaresolverr](https://github.com/astrocode-id/go-flaresolverr) is Golang client for [FlareSolverr](https://github.com/FlareSolverr/FlareSolverr) v3.

## Installation
1. Install [FlareSolverr](https://github.com/FlareSolverr/FlareSolverr#installation)
2. Get [go-flaresolverr](https://github.com/astrocode-id/go-flaresolverr)
```shell
go get github.com/astrocode-id/go-flaresolverr
```

## Examples

### Get Page
Retrieves webpage using [`request.get`](https://github.com/FlareSolverr/FlareSolverr#-requestget) command.

```go
package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"

	"github.com/astrocode-id/go-flaresolverr"
)

func main() {
	c, err := flaresolverr.NewClient(flaresolverr.Config{
		BaseURL: baseURL,
	})
	if err != nil {
		log.Fatal(err)
	}

	b, err := c.Get(flaresolverr.GetParams{
		URL: "https://ifconfig.me",
	})
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	ipAddress := doc.Find("strong").First().Text()
	fmt.Println(ipAddress)
}
```

### Post Page
Retrieves webpage using [`request.post`](https://github.com/FlareSolverr/FlareSolverr#-requestpost) command.

_TODO_

## Note

- :warning: Currently, [FlareSolverr v3](https://github.com/FlareSolverr/FlareSolverr/releases)
doesn't support `session` and `proxy`.
For more detail, see [ChangeLog](https://github.com/FlareSolverr/FlareSolverr/blob/master/CHANGELOG.md#v300-20230104).


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fastrocode-id%2Fgo-flaresolverr.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fastrocode-id%2Fgo-flaresolverr?ref=badge_large)