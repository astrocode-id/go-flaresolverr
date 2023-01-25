# FlareSolverr v3 Go Client

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Code Climate](https://codeclimate.com/github/fahimbagar/go-flaresolverr.png)](https://codeclimate.com/github/fahimbagar/go-flaresolverr)
[![Test Coverage](https://api.codeclimate.com/v1/badges/c8eaaff0f761d4d1f09f/test_coverage)](https://codeclimate.com/github/fahimbagar/go-flaresolverr/test_coverage)
[![GitHub issues](https://img.shields.io/github/issues/fahimbagar/go-flaresolverr)](https://github.com/fahimbagar/go-flaresolverr/issues)
[![CircleCI](https://circleci.com/gh/fahimbagar/go-flaresolverr.svg?style=shield)](https://circleci.com/gh/fahimbagar/go-flaresolverr)

[go-flaresolverr](https://github.com/fahimbagar/go-flaresolverr) is Golang client for [FlareSolverr](https://github.com/FlareSolverr/FlareSolverr) v3. https://github.com/FlareSolverr/FlareSolverr/releases/tag/v3.0.0

## Installation
1. Install [FlareSolverr](https://github.com/FlareSolverr/FlareSolverr#installation)
2. Get [go-flaresolverr](https://github.com/fahimbagar/go-flaresolverr)
```shell
go get github.com/fahimbagar/go-flaresolverr
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

	"github.com/fahimbagar/go-flaresolverr"
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
