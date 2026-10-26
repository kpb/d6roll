# Go library

Import `github.com/kpb/d6roll` to generate passphrases without the CLI. The
module requires Go 1.22 or later and has no third-party dependencies.

## Generate a passphrase

```go
phrase, err := d6roll.Generate(6) // crypto/rand.Reader by default
if err != nil {
    return err
}
words := phrase.Words() // independent []string in selection order
bits := phrase.EntropyBits
roll := phrase.Selections[0].Roll
// Treat words and rolls as secrets; do not log them.
```

For just a `[]string`, call `d6roll.GenerateWords(6)`. `Passphrase.Selections`
holds each `Selection`'s `Word` and associated `Roll`; `Passphrase.Wordlist` is
`"eff-large"`. `Words()` returns a separate slice, so changing it does not
change the selections. The numeric `EntropyBits` value describes the generation
process, not the visual complexity of the result.

The library accepts **1–1,024** words. It returns an error for invalid counts,
a nil random source, or a random-source failure, without returning a successful
partial passphrase. The CLI has its own narrower 3–32-word range.

## Supply randomness for tests

```go
testRandom := bytes.NewReader(make([]byte, 256)) // deterministic, not secure
generator, err := d6roll.NewGenerator(testRandom)
if err != nil {
    return err
}
phrase, err := generator.Generate(6)
if err != nil {
    return err
}
// generator.GenerateWords(6) also returns only the words.
```

`NewGenerator` accepts an `io.Reader` so tests can control randomness. In
production, an injected reader **must** provide cryptographically secure
randomness; use `Generate` or `GenerateWords` for the secure default. Selections
use uniform rejection sampling, independently and with replacement. Repeated
words are valid. Random-source errors are propagated without an insecure
fallback.

## Wordlist, rolls, and entropy

The embedded [EFF long Diceware wordlist](https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt)
contains 7,776 unique words and their original five-die roll identifiers.
Credit: Electronic Frontier Foundation. The bundled file's SHA-256 is
`addd35536511597a02fa0a9ff1e5284677b8883b83e986e43f15a3db996b903e`.
The library validates its size, word and roll uniqueness, and roll syntax at
startup. No external wordlist file is needed at runtime.

A selected word's `Roll` is its existing five-digit list identifier, not an
additional random draw or additional entropy. For wordlists without roll
notation, `Roll` can be empty; custom wordlists are not supported in v0.1.

For `N` independently and uniformly selected words from `W` unique entries,
entropy is `N × log2(W)` bits. The EFF list contributes about 12.925 bits per
word, or 77.5 bits for six words. These estimates require independent uniform
selection and a secure random source. Treat generated passphrases and roll
metadata as secrets; the library does not save or transmit them.

Architectural decisions are documented in [adr/](adr/).
