# ADR 0002: Randomness and word selection

## Status

Accepted

## Context

Biased indices or a weak fallback would undermine entropy estimates even if the wordlist were correct.

## Decision

Use `crypto/rand.Reader` by default and `crypto/rand.Int` for uniform indices in the actual wordlist size. Draw independently with replacement. Permit an explicit `io.Reader` for reproducible tests; production callers must supply a cryptographically secure source. Propagate read errors, never substitute a source, and never return a successful partial phrase.

## Consequences

There is no bespoke random-number implementation or modulo bias. Repeated words are valid. Entropy is `N × log2(W)` only under independent, uniform selection from a list of `W` unique words with a secure source; roll notation does not add entropy.
