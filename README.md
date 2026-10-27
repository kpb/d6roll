<img src="image/alois-komenda-hAdfTKFXmwM-unsplash.jpg" alt="Three white dice on a dark surface" width="720">

Photo by [Alois Komenda][alois-komenda] on [Unsplash][unsplash-photo].

# d6roll

[![License][license-badge]][license]
[![Latest release][release-badge]][latest-release]
[![Tests][tests-badge]][tests-workflow]
[![Go version][go-badge]][go-mod]

d6roll is a small command-line utility for generating [Diceware][eff-dice] passphrases from the [EFF long
wordlist][eff-wordlist]. It has no third-party dependencies or runtime wordlist file. Would you like to use it as a Go
library? See the [library guide][library-guide].

## Install and use

With [Homebrew][homebrew] on macOS (Intel or Apple Silicon) or Linux
(amd64), install from the [kpb tap][kpb-tap]:

```sh
brew install kpb/tap/d6roll
```

On Linux (amd64) or macOS (amd64 or arm64), inspect [install.sh][install-script]
before running the installer:

```sh
curl -fsSL https://raw.githubusercontent.com/kpb/d6roll/main/install.sh | sh
```

The installer downloads the latest release to `~/.local/bin` by default, verifies the archive against the release's
SHA-256 `checksums.txt`, and requires no root access. Set `VERSION` and `INSTALL_DIR` on the right of the pipe to
override either default:

```sh
curl -fsSL https://raw.githubusercontent.com/kpb/d6roll/main/install.sh |
  VERSION=v0.2.1 INSTALL_DIR="$HOME/bin" sh
```

Add the installation directory to your `PATH` if needed (use `$HOME/bin` instead when following the override example):

```sh
export PATH="$HOME/.local/bin:$PATH"
```

Prebuilt archives for Linux (amd64), macOS (amd64 and arm64), and Windows (amd64) are available on the [latest
release][latest-release] page. Download the archive for your platform, verify it against `checksums.txt`, extract it,
and put the binary on your `PATH`. On Windows, use the ZIP archive instead of the shell installer. Alternatively,
install with Go (Go 1.22 or later):

```sh
go install github.com/kpb/d6roll/cmd/d6roll@latest
# Or from a checkout:
go build -o d6roll ./cmd/d6roll

d6roll                       # six space-separated words; stdout only
d6roll -w 8 -s -             # eight hyphen-separated words
d6roll --memorable           # capitalized words with a digit and symbol
d6roll --verbose             # passphrase on stdout, metadata on stderr
d6roll --rolls               # roll table and metadata on stderr
d6roll help                  # usage on stdout
d6roll --help
```

The CLI accepts 3–32 words and separates words with spaces by default. Pass `-s -` or `--separator=-` for
hyphen-separated output. `--verbose` reports word count, entropy, a descriptive strength label, and wordlist
name. `--rolls` also reports the five-digit roll for each word.

Use `-m` or `--memorable` for capitalized words joined by dashes, with a random digit and symbol appended to one
randomly selected word (for example, `Maple-Orbit7!-Velvet-Cabin-Meadow-Quartz`). An explicitly supplied `--separator`
overrides the dash. The digit and symbol are selected securely from `0`–`9` and `!@#$%^&*`. These formatting choices
do not increase the reported entropy, which is based on the selected words.

Redirecting stdout (for example, `d6roll --verbose > passphrase.txt`) captures only the passphrase. Errors go to stderr
and exit nonzero.

## Diceware and security

The bundled [EFF long Diceware wordlist][eff-wordlist], from the [Electronic Frontier Foundation's Dice-Generated
Passphrases guide][eff-dice], has 7,776 words with five-digit roll identifiers. Each uniformly selected word offers
7,776 possibilities, or about 12.9 bits of entropy. Six independent selections give `6 × log2(7,776) ≈ 77.5` bits. This
counts the possible selections, not the words' length or appearance; separators and displayed roll identifiers add no
entropy.

The estimate assumes independent, uniform selection with a secure random source. The CLI's labels are *Low* (<40 bits),
*Moderate* (40–<60), *Strong* (60–<75), and *Very strong* (≥75). They are descriptive, not security guarantees. See the
[library guide][library-guide] for the calculation and assumptions in more detail.

Treat passphrases and displayed rolls as secrets. Output in a terminal may remain in scrollback or captured logs.
*Redirect deliberately*. d6roll does not save or transmit passphrases. Architectural rationale is in [doc/adr/][adr].

## License

I believe in freedom (and free beer). d6roll is licensed under the [GPL-3.0][license].

___

Made with ❤️ by [me][contact]. Leave questions, comments, or suggestions as [GitHub issues][github-issues] or send me an email.

[adr]: doc/adr/
[alois-komenda]: https://unsplash.com/@aloisk?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText
[contact]: mailto:hello@kennethbowen.com
[eff-dice]: https://www.eff.org/dice
[eff-wordlist]: https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt
[github-issues]: https://github.com/kpb/d6roll/issues
[go-badge]: https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go
[go-mod]: go.mod
[homebrew]: https://brew.sh/
[install-script]: install.sh
[kpb-tap]: https://github.com/kpb/homebrew-tap
[latest-release]: https://github.com/kpb/d6roll/releases/latest
[library-guide]: doc/library.md
[license]: LICENSE
[license-badge]: https://img.shields.io/github/license/kpb/d6roll
[release-badge]: https://img.shields.io/github/v/release/kpb/d6roll
[tests-badge]: https://github.com/kpb/d6roll/actions/workflows/test.yml/badge.svg
[tests-workflow]: https://github.com/kpb/d6roll/actions/workflows/test.yml
[unsplash-photo]: https://unsplash.com/photos/3-white-dice-on-black-surface-hAdfTKFXmwM?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText
