# CHANGELOG

## v2.1.0

- Add `GetContext` and `PostContext`; `Get` and `Post` now call them with `context.Background()` [#14](https://github.com/astrocode-id/go-flaresolverr/pull/14)
- Add `Config.HTTPClient`, and reuse one `http.Client` per `Client` instead of creating one per call [#14](https://github.com/astrocode-id/go-flaresolverr/pull/14)
- Fix README install command and imports to use the `/v2` module path [#14](https://github.com/astrocode-id/go-flaresolverr/pull/14)

## v2.0.2

- Add the `/v2` suffix to the module path in go.mod. v2.0.0 and v2.0.1 declared no suffix, so Go could not resolve them as v2 [#11](https://github.com/astrocode-id/go-flaresolverr/pull/11)
- Run integration tests over plain HTTP against a FlareSolverr service container in CI, replacing testcontainers-go [#11](https://github.com/astrocode-id/go-flaresolverr/pull/11)

## v2.0.1

- Bump indirect dependency `github.com/moby/go-archive` from 0.2.0 to 0.3.0 [#10](https://github.com/astrocode-id/go-flaresolverr/pull/10)

## v2.0.0

Breaking: `Get` and `Post` now take functional options and return the full `Response`.

- Replace the struct-param API with `Get(url, ...Option)` and `Post(url, postData, ...Option)`, removing the old `Get`/`Post` and renaming `GetRaw`/`PostRaw` [#7](https://github.com/astrocode-id/go-flaresolverr/pull/7)
- Add `WithMaxTimeout`, `WithCookies`, `WithReturnOnlyCookies`, `WithWaitInSeconds`, and `WithDisableMedia` [#7](https://github.com/astrocode-id/go-flaresolverr/pull/7)
- Export `Cookie` and use `Cookies []Cookie` instead of `[]http.Cookie`; rename the cookie `expiry` field to `expires` to match FlareSolverr [#7](https://github.com/astrocode-id/go-flaresolverr/pull/7)
- Decode `Solution.Response` as a `string` instead of `json.RawMessage` [#7](https://github.com/astrocode-id/go-flaresolverr/pull/7)
- Move integration tests from dockertest to testcontainers-go and require Go 1.25 [#8](https://github.com/astrocode-id/go-flaresolverr/pull/8)

## v1.1.0

- Migrate CI from CircleCI to GitHub Actions (unit and integration tests, gofmt, actionlint, golangci-lint) [#6](https://github.com/astrocode-id/go-flaresolverr/pull/6)

## v1.0.0

- Add command `request.get` and `request.post` [#3](https://github.com/astrocode-id/go-flaresolverr/pull/3)
- Add integration tests folder [#3](https://github.com/astrocode-id/go-flaresolverr/pull/3)
- Add CircleCI workflow
