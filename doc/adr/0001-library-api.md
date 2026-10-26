# ADR 0001: Library API

## Status

Accepted

## Context

Callers need both simple word slices and richer results with entropy and Diceware rolls. Tests need controlled randomness without making insecure randomness the default. Custom wordlists are out of scope for v0.1.

## Decision

Expose `Generate(n)` and `GenerateWords(n)` using secure defaults, and `NewGenerator(io.Reader)` with corresponding methods for explicit randomness injection. `Passphrase` contains selections, entropy bits, and wordlist identity; `Words()` returns a separate slice. Each `Selection` contains a word and an optional roll (empty string when not applicable). Keep generator state and wordlist representation private.

## Consequences

Simple callers avoid configuration; deterministic tests use a generator without functional options. The paired convenience methods add a little API surface but avoid forcing callers to handle unused metadata. Public result fields are plain, mutable data, not promises of immutable or independently verified security. Future wordlists can have different sizes and no rolls without changing the result types; loading them is not part of v0.1.
