# d6roll

d6roll is a small Go library and command-line utility for generating Diceware passphrases from the [EFF long wordlist](https://www.eff.org/dice). It has no third-party dependencies or runtime wordlist file.

## Build and use

Requires Go 1.22 or later.

```sh
go install github.com/kpb/d6roll/cmd/d6roll@latest
# Or from a checkout:
go build -o d6roll ./cmd/d6roll

d6roll                       # six words, hyphen-separated; stdout only
d6roll -w 8 -s ' '          # eight space-separated words
d6roll --verbose             # passphrase on stdout, metadata on stderr
d6roll --rolls               # roll table and metadata on stderr
d6roll --help
```

The CLI accepts 3–32 words; the library accepts any positive word count subject to available resources. `--verbose` reports word count, entropy, a descriptive strength label, and wordlist name. `--rolls` also reports the five-digit roll for each word. Redirecting stdout (for example, `d6roll --verbose > passphrase.txt`) captures only the passphrase. Errors go to stderr and exit nonzero.

## Library

```go
phrase, err := d6roll.Generate(6) // secure default source
if err != nil {
    return err
}
words := phrase.Words()             // []string, separate from the selections
bits := phrase.EntropyBits
roll := phrase.Selections[0].Roll
_, _ = words, bits
_ = roll
```

For just the words, use `d6roll.GenerateWords(6)`. For deterministic tests, use `d6roll.NewGenerator(reader)`, then `generator.Generate(6)` or `generator.GenerateWords(6)`. In production, an injected `io.Reader` **must** provide cryptographically secure randomness; the defaults use `crypto/rand.Reader`. Invalid counts, nil sources, and random-source failures return errors, never successful partial results.

## Wordlist, rolls, and entropy

The bundled [EFF long Diceware wordlist](https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt) contains 7,776 distinct words with their original five-die roll identifiers. Credit: Electronic Frontier Foundation. The embedded source file's SHA-256 is `addd35536511597a02fa0a9ff1e5284677b8883b83e986e43f15a3db996b903e`. The program validates the embedded list's size, unique words and rolls, and roll syntax at startup. A roll identifies the selected word; it is **not** additional randomness or entropy. Future lists could have no rolls, represented by an empty `Selection.Roll`; custom list loading is not supported in v0.1.

For `N` independently selected words from `W` unique entries, entropy is `N × log2(W)` bits. With the EFF list, that is about 12.925 bits per word, or 77.5 bits for six words. This estimate assumes a secure source and unbiased, independent selection with replacement; repeated words are valid. The CLI's display labels are *Low* (<40 bits), *Moderate* (40–<60), *Strong* (60–<75), and *Very strong* (≥75). They are descriptive buckets, not universal security guarantees; the numerical entropy is authoritative.

Treat passphrases and displayed rolls as secrets. Output in a terminal may remain in scrollback or captured logs; redirect deliberately. d6roll does not save or transmit passphrases. Architectural rationale is in [doc/adr/](doc/adr/). The repository license is [GPL-3.0](LICENSE).
