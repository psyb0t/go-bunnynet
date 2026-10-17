---
name: go-bunnynet
description: Use the unofficial Go client for the Bunny.net core API from Go code. Package bunnynet at github.com/psyb0t/go-bunnynet provides generated CDN, DNS, storage and video library configuration, statistics, load balancer and billing operations. Use when a Go project imports go-bunnynet or needs Bunny.net core API calls. Excludes storage file transfers and the separate Stream API.
homepage: https://github.com/psyb0t/go-bunnynet
user-invocable: true
metadata:
  openclaw:
    emoji: "🐰"
    requires:
      bins:
        - go
      env:
        - BUNNYNET_API_KEY
permissions:
  filesystem: "Guides edits to Go source in the current project. The library reads and writes no files."
  shell: "The Go toolchain to add and compile the dependency."
  network: "Module downloads use the configured Go proxy. Runtime core API requests use https://api.bunny.net."
---

# go-bunnynet

Unofficial, typed Go client for the Bunny.net core API, in package `bunnynet` at the module root. Read [references/setup.md](references/setup.md) for installation, methods and response handling.

## Security and safety

- Read `BUNNYNET_API_KEY` from the environment and set `AccessKey` through `WithRequestEditorFn`. Never hardcode or log credentials.
- Use `https://api.bunny.net`. Core account keys differ from storage zone passwords and video library keys.
- Creating, changing or deleting resources, purging caches, rotating keys and closing accounts act on the user's account. Call those operations only when requested.

## Usage

Use the README example, compiled as `Example` in `example_test.go`. Construct `bunnynet.NewClientWithResponses` with `bunnynet.WithRequestEditorFn`; attach `AccessKey` in the editor. Supply `bunnynet.WithHTTPClient` with a timeout when needed.

`PullZonePublicIndexAllWithResponse(ctx, nil)` lists pull zones. `PullZonePublicIndexWithResponse(ctx, id, nil)` fetches one. Operation IDs determine method names, with punctuation removed.

Check the Go error, HTTP status and decoded field separately. API errors can return a nil Go error. Optional fields are pointers. Success responses with no schema expose status and raw body only.

Pull-zone listing declares an array but describes an object when `page > 0`. Leave `Page` unset for typed decoding or use raw `PullZonePublicIndexAll` and decode the paginated body. There are no automatic retries or pagination.

For changes to the generated client, read `AGENTS.md` and regenerate from the spec.
