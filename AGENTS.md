# Agent Instructions

## Project Context

IMPORTANT: Before you do anything else, run beanstalk prime and follow its output.

d6roll is a Go module (`github.com/kpb/d6roll`) with a reusable root package and a thin CLI in `cmd/d6roll`. The release work is tracked in the v0.1.0 Beanstalk milestone; claim tasks before implementing them and record their outcome. Preserve existing GitHub history and the repository's GPL-3.0 license.

The initial public API was approved for incremental review. Keep it small; explain any necessary material API change before making it and record changes to accepted decisions in a superseding ADR. Architecture decisions live in `doc/adr/`.

Use only secure, unbiased random selection by default, keep formatting in the CLI, and keep passphrases out of logs and errors. CLI stdout contains only the passphrase. Run `go build ./...`, `go test ./...`, and `go vet ./...` before completing the release.
