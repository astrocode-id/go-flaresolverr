# FlareSolverr v3 Go Client

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub issues](https://img.shields.io/github/issues/astrocode-id/go-flaresolverr)](https://github.com/astrocode-id/go-flaresolverr/issues)
[![static analysis](https://github.com/astrocode-id/go-flaresolverr/actions/workflows/static-analysis.yml/badge.svg)](https://github.com/astrocode-id/go-flaresolverr/actions/workflows/static-analysis.yml)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fastrocode-id%2Fgo-flaresolverr.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fastrocode-id%2Fgo-flaresolverr?ref=badge_shield)

[go-flaresolverr](https://github.com/astrocode-id/go-flaresolverr) is Golang client for [FlareSolverr](https://github.com/FlareSolverr/FlareSolverr) v3.

## Installation
1. Install [FlareSolverr](https://github.com/FlareSolverr/FlareSolverr#installation)
2. Get [go-flaresolverr](https://github.com/astrocode-id/go-flaresolverr)
```shell
go get github.com/astrocode-id/go-flaresolverr/v2
```

## Examples

### Get Page
Retrieves webpage using [`request.get`](https://github.com/FlareSolverr/FlareSolverr#-requestget) command.

```go
package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/astrocode-id/go-flaresolverr/v2"
)

func main() {
	c, err := flaresolverr.NewClient(flaresolverr.Config{
		BaseURL: "http://localhost:8191/v1",
	})
	if err != nil {
		log.Fatal(err)
	}

	r, err := c.Get("https://httpbin.org/ip",
		flaresolverr.WithMaxTimeout(60000),
	)
	if err != nil {
		log.Fatal(err)
	}
	if r.Status != "ok" {
		log.Fatal(errors.New(r.Message))
	}

	fmt.Println(r.Solution.Response)
}
```

### Post Page
Retrieves webpage using [`request.post`](https://github.com/FlareSolverr/FlareSolverr#-requestpost) command.

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/astrocode-id/go-flaresolverr/v2"
)

func main() {
	c, err := flaresolverr.NewClient(flaresolverr.Config{
		BaseURL: "http://localhost:8191/v1",
	})
	if err != nil {
		log.Fatal(err)
	}

	r, err := c.Post("https://httpbin.org/post", url.Values{
		"key1": {"valueA"},
		"key2": {"valueB"},
	})
	if err != nil {
		log.Fatal(err)
	}
	if r.Status != "ok" {
		log.Fatal(errors.New(r.Message))
	}

	fmt.Println(r.Solution.Response)
}
```

### Options
Both `Get` and `Post` accept optional parameters, which map to the
[base parameters](https://github.com/FlareSolverr/FlareSolverr#-requestget)
supported by the FlareSolverr API:

- `WithMaxTimeout(ms int)`
- `WithCookies(cookies flaresolverr.Cookies)`
- `WithReturnOnlyCookies(v bool)`
- `WithWaitInSeconds(s int)`
- `WithDisableMedia(v bool)`

`session`, `session_ttl_minutes`, `proxy`, `returnScreenshot`, and
`tabs_till_verify` are part of the FlareSolverr API but aren't supported by
this client yet.

### Context and timeouts
`GetContext` and `PostContext` take a `context.Context` that bounds the HTTP
call to FlareSolverr. Keep the deadline above `WithMaxTimeout`, since
FlareSolverr holds the connection open until the challenge is solved.

```go
ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
defer cancel()

r, err := c.GetContext(ctx, "https://httpbin.org/ip",
	flaresolverr.WithMaxTimeout(60000),
)
```

`Get` and `Post` use `context.Background()`. To bound them too, pass an
`http.Client` with a `Timeout` through `Config.HTTPClient`.

### Error handling
The returned `error` covers transport and decoding failures only. When
FlareSolverr itself fails (for example, the challenge times out), it replies
with `"status": "error"`, and the client returns that `Response` with a nil
error. Check `r.Status` yourself, as the examples above do:

```go
if r.Status != "ok" {
	return errors.New(r.Message)
}
```

Returning an error for these responses would change behavior for callers who
already check `r.Status`, so it is left for a future major version.

### Session
_TODO_: not supported yet.

### Proxy
_TODO_: not supported yet.

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fastrocode-id%2Fgo-flaresolverr.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fastrocode-id%2Fgo-flaresolverr?ref=badge_large)
