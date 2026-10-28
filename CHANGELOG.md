# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.3.4, 2026-10-28

- Documentation only. The doc comments on `VideoLibraryModel` and `VideoLibraryUpdateModel` go back to the shorter upstream wording and no longer describe how DRM interacts with early play and original files.
- `VideoLibraryModel.EnableDRM`, `AllowEarlyPlay` and `ExposeOriginals`: the notes that `AllowEarlyPlay` and `ExposeOriginals` are always false while `EnableDRM` is true are removed.
- `VideoLibraryUpdateModel.EnableDRM`, `AllowEarlyPlay` and `ExposeOriginals`: the notes that `EnableDRM` cannot be true together with `AllowEarlyPlay` or `ExposeOriginals` are removed. The note that enabling `AllowEarlyPlay` enables `ExposeOriginals` stays in shorter form.

## v0.3.3, 2026-10-28

- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.3.2, 2026-10-27

- Documentation only. The doc comments on `VideoLibraryModel` and `VideoLibraryUpdateModel` now describe how DRM interacts with early play and original files.
- `VideoLibraryModel.EnableDRM`, `AllowEarlyPlay` and `ExposeOriginals`: while `EnableDRM` is true, `AllowEarlyPlay` and `ExposeOriginals` are always false.
- `VideoLibraryUpdateModel.EnableDRM` cannot be set to true while `AllowEarlyPlay` or `ExposeOriginals` is true. Setting `AllowEarlyPlay` to true also sets `ExposeOriginals` to true. Original files are not protected by DRM.

## v0.3.1, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.
- Raised the minimum Go version in `go.mod` from 1.26.9 to 1.26.10, which fixes the `crypto/x509` vulnerability GO-2026-6809 found by govulncheck.

## v0.3.0, 2026-10-21

- The Bunny.net spec went back to the `StorageZonePublic` operation ID for `GET /storagezone/regions`. `(*Client).GetStoragezoneRegionsEndpointGetStorageZoneRegions` and `(*ClientWithResponses).GetStoragezoneRegionsEndpointGetStorageZoneRegionsWithResponse` are removed. Use `(*Client).StorageZonePublicGetStorageZoneRegions` and `(*ClientWithResponses).StorageZonePublicGetStorageZoneRegionsWithResponse`.
- The matching request builder, response type and parse function are renamed: `NewGetStoragezoneRegionsEndpointGetStorageZoneRegionsRequest` is now `NewStorageZonePublicGetStorageZoneRegionsRequest`, `GetStoragezoneRegionsEndpointGetStorageZoneRegionsResponse` is now `StorageZonePublicGetStorageZoneRegionsResponse`, and `ParseGetStoragezoneRegionsEndpointGetStorageZoneRegionsResponse` is now `ParseStorageZonePublicGetStorageZoneRegionsResponse`. `ClientInterface` and `ClientWithResponsesInterface` change the same way.
- `StorageRegionModel` no longer has `Id` or `Url`. It now has `Code`, `Continent`, `Endpoint`, `IsMainRegionSupported`, `IsS3Supported`, `Latitude`, `Longitude`, `PricePerGigabyte`, `S3Endpoint` and `Tier`.
- The storage regions response no longer has a 401 case with no content type.
- `DnsZonePublicDeleteResponse` and `DnsZonePublicDeleteRecordResponse` now have `JSON409` and `GetJSON409`.
- Added the `OptimizerHtmlToMarkdownTypePrerendered` constant.
- `VideoLibraryModel` has new fields `LiveStreamCount` and `TotalCount`.
- The development image in `Dockerfile.dev` moves to Go 1.26.10, which fixes GO-2026-6809 in `crypto/x509` reported by `make sec`.

## v0.2.0, 2026-10-19

- `(*Client).StorageZonePublicGetStorageZoneRegions` and `(*ClientWithResponses).StorageZonePublicGetStorageZoneRegionsWithResponse` are gone. Use `(*Client).GetStoragezoneRegionsEndpointGetStorageZoneRegions` and `(*ClientWithResponses).GetStoragezoneRegionsEndpointGetStorageZoneRegionsWithResponse` instead, matching a new operation ID in the Bunny.net spec for `GET /storagezone/regions`.
- The matching request builder, response type and parse function are renamed too: `NewStorageZonePublicGetStorageZoneRegionsRequest` is now `NewGetStoragezoneRegionsEndpointGetStorageZoneRegionsRequest`, `StorageZonePublicGetStorageZoneRegionsResponse` is now `GetStoragezoneRegionsEndpointGetStorageZoneRegionsResponse`, and `ParseStorageZonePublicGetStorageZoneRegionsResponse` is now `ParseGetStoragezoneRegionsEndpointGetStorageZoneRegionsResponse`.
- `StorageRegionModel` no longer has `Code`, `Continent`, `Hostname`, `IsS3Supported`, `Latitude`, `Longitude`, `PricePerGigabyte`, `S3Hostname` or `Tier`. It now has `Id` and `Url` instead. `Name` is unchanged.
- `GetStoragezoneRegionsEndpointGetStorageZoneRegionsResponse` now also recognizes a 401 status with no content type, instead of falling through unhandled.

## v0.1.1, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.0, 2026-10-17

- Added `bunnynet`, a REST client generated with oapi-codegen from the Bunny.net core API spec version 1.0.0, OpenAPI 3.0.0, at https://core-api-public-docs.b-cdn.net/docs/v3/public.json. All 111 operations across 88 paths are included: CDN pull zones, cache purges, DNS, storage and video library configuration, statistics, load balancers and billing.
- Number fields use `float64` instead of the default `float32`.
- `make generate` downloads the current spec and regenerates the client without spec patches. The spec is not redistributed.
- Added local HTTP contract tests, the compiled README example, development tooling, CI workflows and the `go-bunnynet` agent skill.
