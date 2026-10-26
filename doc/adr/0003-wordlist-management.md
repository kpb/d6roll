# ADR 0003: Wordlist management

## Status

Accepted

## Context

The CLI must work without runtime data files while retaining authentic word-to-roll mappings, and the selection algorithm must not depend on a specific list size or roll format.

## Decision

Embed the official EFF long list with `go:embed`. Parse and validate exactly 7,776 entries, unique nonempty words and rolls, and five-digit rolls using only 1–6. Store each word with its supplied roll; select by index and return that entry's roll without another random draw. Keep the wordlist internal to v0.1; the index-selection core accepts any list size.

## Consequences

Distributed binaries need no wordlist path. A corrupted bundled list fails early rather than producing misleading entropy or mappings. A later custom list could omit rolls while reusing the same selector, but loading and validating custom lists are not implemented now.
