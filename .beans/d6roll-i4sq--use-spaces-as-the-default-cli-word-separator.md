---
# d6roll-i4sq
title: Use spaces as the default CLI word separator
status: completed
type: feature
priority: normal
created_at: 2026-10-26T23:04:59Z
updated_at: "2026-10-27T02:54:49Z"
---
Change the d6roll CLI default separator from a hyphen to a space for both `-s` and `--separator` flag defaults. Preserve explicit custom separators, including `-s -` for the old behavior. Update default, verbose, and rolls output tests; cover the explicit hyphen override; update CLI README examples and help defaults. Keep the library API, word selection, entropy, and metadata unchanged. Note that scripts relying on hyphen-separated default output must pass `-s -`. Do not implement until requested.
