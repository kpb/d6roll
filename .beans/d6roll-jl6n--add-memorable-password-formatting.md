---
# d6roll-jl6n
title: Add memorable password formatting
status: completed
type: feature
priority: normal
created_at: 2026-10-27T23:03:09Z
updated_at: "2026-10-27T23:09:38Z"
---
## Summary
Add `-m` / `--memorable` CLI mode that formats the generated Diceware words as a memorable password.

## Acceptance criteria
- Add `-m` and `--memorable` aliases.
- Capitalize the first letter of every selected word.
- Join words with dashes by default; honor an explicit `--separator`.
- Uniformly select one word and append a random digit followed by a random symbol to it.
- Use secure, unbiased randomness and a documented symbol set.
- Preserve word-count options and existing behavior without the flag.
- Keep formatting in the CLI and passphrases out of diagnostics.
- Keep reported entropy based on word selection.
- Update help, documentation, and tests, including random-source failure handling.

Example: `Maple-Orbit7!-Velvet-Cabin-Meadow-Quartz`

## Outcome
- Implemented the CLI preset, help, and README documentation without changing the library API.
- Review follow-up added deterministic tests for interior and last-word placement, maximum digit and symbol values, rejection sampling, explicit space and empty separators through both aliases, and failures at each random draw without passphrase or metadata output.
- Validation passed: `go build ./...`, `go test ./...`, `go vet ./...`, and `git diff --check`.
