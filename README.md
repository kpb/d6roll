<img src="image/alois-komenda-hAdfTKFXmwM-unsplash.jpg" alt="Three white dice on a dark surface" width="720">

Photo by [Alois Komenda](https://unsplash.com/@aloisk?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText) on [Unsplash](https://unsplash.com/photos/3-white-dice-on-black-surface-hAdfTKFXmwM?utm_source=unsplash&utm_medium=referral&utm_content=creditCopyText).

# d6roll

[![License][license-badge]][license]
[![Latest release][release-badge]][latest-release]
[![Tests][tests-badge]][tests-workflow]
[![Go version][go-badge]][go-mod]

d6roll is a small command-line utility for generating Diceware passphrases from the [EFF long wordlist](https://www.eff.org/dice). It has no third-party dependencies or runtime wordlist file. Using it from Go? See the [library guide](doc/library.md).

## Install and use

On Linux (amd64) or macOS (amd64 or arm64), inspect [install.sh](install.sh)
before running the installer:

```sh
curl -fsSL https://raw.githubusercontent.com/kpb/d6roll/main/install.sh | sh
```

The installer downloads the latest release to `~/.local/bin` by default, verifies
the archive against the release's SHA-256 `checksums.txt`, and requires no root
access. Set `VERSION` and `INSTALL_DIR` on the right of the pipe to override
either default:

```sh
curl -fsSL https://raw.githubusercontent.com/kpb/d6roll/main/install.sh |
  VERSION=v0.1.0 INSTALL_DIR="$HOME/bin" sh
```

Add the installation directory to your `PATH` if needed (use `$HOME/bin` instead
when following the override example):

```sh
export PATH="$HOME/.local/bin:$PATH"
```

Checksums downloaded from the same release detect corruption, but do not
independently authenticate the publisher. Review the script before running it.
In a `curl | sh` pipeline, curl failure can be masked (or a partial script may
execute); download the script separately when a reliable exit status matters.
The installer does not support Windows or Linux arm64.

Prebuilt archives for Linux (amd64), macOS (amd64 and arm64), and Windows (amd64)
are available on the [latest release][latest-release] page. Download the archive
for your platform, verify it against `checksums.txt`, extract it, and put the
binary on your `PATH`. On Windows, use the ZIP archive instead of the shell
installer. Alternatively, install with Go (Go 1.22 or later):

```sh
go install github.com/kpb/d6roll/cmd/d6roll@latest
# Or from a checkout:
go build -o d6roll ./cmd/d6roll

d6roll                       # six words, hyphen-separated; stdout only
d6roll -w 8 -s ' '          # eight space-separated words
d6roll --verbose             # passphrase on stdout, metadata on stderr
d6roll --rolls               # roll table and metadata on stderr
d6roll help                  # usage on stdout
d6roll --help
```

The CLI accepts 3–32 words. `--verbose` reports word count, entropy, a descriptive strength label, and wordlist name. `--rolls` also reports the five-digit roll for each word. Redirecting stdout (for example, `d6roll --verbose > passphrase.txt`) captures only the passphrase. Errors go to stderr and exit nonzero.

## Diceware and security

The bundled [EFF long Diceware wordlist](https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt) has 7,776 words with five-digit roll identifiers. Credit: Electronic Frontier Foundation. Each displayed roll identifies its word; it adds no entropy. Six independently selected words have about 77.5 bits of entropy, assuming uniform selection with a secure random source. The CLI's labels are *Low* (<40 bits), *Moderate* (40–<60), *Strong* (60–<75), and *Very strong* (≥75). These are descriptive buckets, not security guarantees; the numeric estimate is authoritative. See the [library guide](doc/library.md) for the calculation and metadata details.

Treat passphrases and displayed rolls as secrets. Output in a terminal may remain in scrollback or captured logs; redirect deliberately. d6roll does not save or transmit passphrases. Architectural rationale is in [doc/adr/](doc/adr/). The repository license is [GPL-3.0](LICENSE).

Made with ❤️ by [me](mailto:hello@kennethbowen.com). Leave questions, comments, or suggestions as [GitHub issues](https://github.com/kpb/d6roll/issues) or send me an email.

[go-badge]: https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go
[go-mod]: go.mod
[latest-release]: https://github.com/kpb/d6roll/releases/latest
[license-badge]: https://img.shields.io/github/license/kpb/d6roll
[license]: LICENSE
[release-badge]: https://img.shields.io/github/v/release/kpb/d6roll
[tests-badge]: https://github.com/kpb/d6roll/actions/workflows/test.yml/badge.svg
[tests-workflow]: https://github.com/kpb/d6roll/actions/workflows/test.yml
