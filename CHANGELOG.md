# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.1, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.0, 2026-10-17

- Added `bunnynet`, a REST client generated with oapi-codegen from the Bunny.net core API spec version 1.0.0, OpenAPI 3.0.0, at https://core-api-public-docs.b-cdn.net/docs/v3/public.json. All 111 operations across 88 paths are included: CDN pull zones, cache purges, DNS, storage and video library configuration, statistics, load balancers and billing.
- Number fields use `float64` instead of the default `float32`.
- `make generate` downloads the current spec and regenerates the client without spec patches. The spec is not redistributed.
- Added local HTTP contract tests, the compiled README example, development tooling, CI workflows and the `go-bunnynet` agent skill.
