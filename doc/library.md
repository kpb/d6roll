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
comes from the [Electronic Frontier Foundation's Dice-Generated Passphrases guide](https://www.eff.org/dice).
It contains 7,776 unique words and their original five-die roll identifiers.
The bundled file's SHA-256 is
`addd35536511597a02fa0a9ff1e5284677b8883b83e986e43f15a3db996b903e`.
The library validates its size, word and roll uniqueness, and roll syntax at
startup. No external wordlist file is needed at runtime.

A selected word's `Roll` is its existing five-digit list identifier, not an
additional random draw or additional entropy. For wordlists without roll
notation, `Roll` can be empty; custom wordlists are not supported in v0.1.

Entropy measures how many possible outputs the generation process could have
produced, not how complicated the resulting words look. For a wordlist of `W`
unique entries, one uniformly selected word has `W` possibilities. Selecting
`N` words independently **with replacement** gives `W^N` ordered passphrases;
any particular sequence has probability `1 / W^N`. Therefore its entropy is
`H = log2(W^N) = N × log2(W)` bits.

For the EFF list, `W = 7,776 = 6^5`, so each word contributes
`log2(7,776) ≈ 12.925` bits. Six words yield about 77.55 bits, displayed by
the CLI as **77.5 bits**. Repeated words are valid: a sequence with repeats has
the same probability as any other specific sequence. Word length, spelling,
separators, and roll notation do not add choices, so they do not add entropy.

This estimate assumes the list has unique words and that each selection is
uniform, independent, and made with a cryptographically secure source. A weak
or predictable injected reader invalidates the security interpretation of the
number, even though the formula still describes the intended selection process.
Entropy is not a guarantee about a particular account or threat: a disclosed
or reused passphrase may be unsafe regardless of its calculated entropy. Treat
generated passphrases and roll metadata as secrets; the library does not save
or transmit them.

Architectural decisions are documented in [adr/](adr/).
