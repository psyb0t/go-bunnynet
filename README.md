# go-bunnynet

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-bunnynet.svg)](https://pkg.go.dev/github.com/psyb0t/go-bunnynet)
[![CI](https://github.com/psyb0t/go-bunnynet/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-bunnynet/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-bunnynet/badges/version.svg)](https://github.com/psyb0t/go-bunnynet/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-bunnynet/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-bunnynet/badges/importers.svg)](https://github.com/psyb0t/go-bunnynet/blob/badges/importers.md)

Unofficial Go client for the [Bunny.net core API](https://core-api-public-docs.b-cdn.net/docs/v3/public.json), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI document Bunny.net publishes. I am not affiliated with Bunny.net.

The package is `bunnynet`, in the repository root. Everything in `client.gen.go` is generated. Don't edit it by hand, regenerate it. All 111 operations in the core spec are included: CDN pull zones, cache purges, DNS, storage and video library configuration, statistics, load balancers and billing. Storage file transfers and the separate Stream API are outside this spec.

## Install

```sh
go get github.com/psyb0t/go-bunnynet
```

## Usage

Every endpoint has a method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-bunnynet"
)

func main() {
	auth := bunnynet.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("AccessKey", os.Getenv("BUNNYNET_API_KEY"))

			return nil
		},
	)

	client, err := bunnynet.NewClientWithResponses(
		"https://api.bunny.net",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.PullZonePublicIndexAllWithResponse(
		context.Background(),
		nil,
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles. Method names come from operation IDs: `PullZonePublic_IndexAll` becomes `PullZonePublicIndexAllWithResponse`.

Check `StatusCode()` and the decoded status field before using a response. The error return covers transport and decoding failures. An HTTP error can return a nil Go error, with its status and raw `Body` intact. Success responses with no schema, such as adding a pull zone, expose status and raw body rather than a typed success field.

Pull-zone listing declares an array schema, but its prose describes a paginated object for `page > 0`. Leave `Page` unset for typed array decoding. For paginated listing, use the raw `PullZonePublicIndexAll` method and decode the body yourself. The client does not paginate or retry automatically.

## Auth

The default scheme puts the account API key in the `AccessKey` header. Set it through `WithRequestEditorFn`, as above, using `BUNNYNET_API_KEY`. The base URL is `https://api.bunny.net`. Storage zone passwords and video library keys belong to separate APIs.

Some operations also declare an alternative `Authorization` header scheme named `bearer`. For that scheme, use `req.Header.Set("Authorization", token)`, supplying the full expected header value. The client adds no prefix automatically. Never log either header.

## Regenerating

The spec is not committed. `make generate` downloads it into the gitignored `api/` directory and regenerates the client:

```sh
make generate
```

`gen.go` declares the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `https://core-api-public-docs.b-cdn.net/docs/v3/public.json` to `api/openapi.json` and checks the document.
- oapi-codegen runs with `config.yaml`. Unformatted `number` fields map to `float64` instead of `float32`. No spec patches are needed.

`.github/workflows/regen.yml` repeats this daily. It checks lint, tests and security, permits bounded repair attempts outside its protected paths, and releases passing updates with a CHANGELOG entry.

## Agent integrations

The [skill](.agents/skills/go-bunnynet) teaches a coding agent how to use go-bunnynet when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-bunnynet@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-bunnynet@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-bunnynet
```

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the spec and regenerate
```

## License

[MIT](LICENSE). Vendored dependencies retain their own licenses, including GPL development tooling. The generated code comes from Bunny.net's published API specification. The specification is not redistributed here. It declares no license and links to [Bunny.net's terms of service](https://bunny.net/tos); this repository's MIT license does not grant rights to the provider's specification or service.
