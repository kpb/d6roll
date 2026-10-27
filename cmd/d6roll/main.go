package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/kpb/d6roll"
)

const version = "0.2.1"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, d6roll.Generate))
}

func run(args []string, stdout, stderr io.Writer, generate func(int) (d6roll.Passphrase, error)) int {
	return runWithRandom(args, stdout, stderr, generate, rand.Reader)
}

func runWithRandom(args []string, stdout, stderr io.Writer, generate func(int) (d6roll.Passphrase, error), random io.Reader) int {
	flags := flag.NewFlagSet("d6roll", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var count int
	var separator string
	var verbose, rolls, help, showVersion, memorable bool
	flags.IntVar(&count, "w", 6, "number of words (3–32)")
	flags.IntVar(&count, "words", 6, "number of words (3–32)")
	flags.StringVar(&separator, "s", " ", "word separator")
	flags.StringVar(&separator, "separator", " ", "word separator")
	flags.BoolVar(&verbose, "v", false, "show metadata on stderr")
	flags.BoolVar(&verbose, "verbose", false, "show metadata on stderr")
	flags.BoolVar(&rolls, "rolls", false, "show Diceware rolls and metadata on stderr")
	flags.BoolVar(&memorable, "m", false, "format words as a memorable password")
	flags.BoolVar(&memorable, "memorable", false, "format words as a memorable password")
	flags.BoolVar(&help, "h", false, "show help")
	flags.BoolVar(&help, "help", false, "show help")
	flags.BoolVar(&showVersion, "V", false, "show version")
	flags.BoolVar(&showVersion, "version", false, "show version")
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: d6roll [flags]")
		fmt.Fprintln(flags.Output(), "       d6roll help")
		fmt.Fprintln(flags.Output(), "")
		fmt.Fprintln(flags.Output(), "Flags:")
		fmt.Fprintln(flags.Output(), "  -h, --help             show help")
		fmt.Fprintln(flags.Output(), "  -V, --version          show version")
		fmt.Fprintln(flags.Output(), "      --rolls            show Diceware rolls and metadata on stderr")
		fmt.Fprintln(flags.Output(), "  -m, --memorable        capitalize words, use dashes, add a digit and symbol")
		fmt.Fprintln(flags.Output(), "  -s, --separator string word separator (default \" \"; memorable: \"-\")")
		fmt.Fprintln(flags.Output(), "  -v, --verbose          show metadata on stderr")
		fmt.Fprintln(flags.Output(), "  -w, --words int        number of words (3–32) (default 6)")
	}
	if len(args) == 1 && args[0] == "help" {
		flags.SetOutput(stdout)
		flags.Usage()
		return 0
	}
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if help {
		flags.Usage()
		return 0
	}
	if showVersion {
		fmt.Fprintf(stdout, "d6roll %s\n", version)
		return 0
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "d6roll: unexpected positional arguments")
		return 2
	}
	if count < 3 || count > 32 {
		fmt.Fprintln(stderr, "d6roll: word count must be between 3 and 32")
		return 2
	}

	phrase, err := generate(count)
	if err != nil {
		fmt.Fprintf(stderr, "d6roll: %v\n", err)
		return 1
	}
	words := phrase.Words()
	if memorable {
		separatorSet := false
		flags.Visit(func(f *flag.Flag) {
			if f.Name == "s" || f.Name == "separator" {
				separatorSet = true
			}
		})
		if !separatorSet {
			separator = "-"
		}
		for i, word := range words {
			if word != "" {
				first, size := utf8.DecodeRuneInString(word)
				words[i] = string(unicode.ToUpper(first)) + word[size:]
			}
		}
		index, err := randomIndex(random, len(words))
		if err != nil {
			fmt.Fprintln(stderr, "d6roll: memorable formatting failed")
			return 1
		}
		digit, err := randomIndex(random, 10)
		if err != nil {
			fmt.Fprintln(stderr, "d6roll: memorable formatting failed")
			return 1
		}
		const symbols = "!@#$%^&*"
		symbol, err := randomIndex(random, len(symbols))
		if err != nil {
			fmt.Fprintln(stderr, "d6roll: memorable formatting failed")
			return 1
		}
		words[index] += fmt.Sprintf("%d%c", digit, symbols[symbol])
	}
	if _, err := fmt.Fprintln(stdout, strings.Join(words, separator)); err != nil {
		fmt.Fprintf(stderr, "d6roll: write output: %v\n", err)
		return 1
	}
	if rolls {
		fmt.Fprintln(stderr, "\nRoll     Word")
		fmt.Fprintln(stderr, "-----    ----")
		for _, selection := range phrase.Selections {
			fmt.Fprintf(stderr, "%-8s %s\n", selection.Roll, selection.Word)
		}
		fmt.Fprintln(stderr)
	}
	if verbose || rolls {
		fmt.Fprintf(stderr, "Words:     %d\nEntropy:   %.1f bits\nStrength:  %s\nWordlist:  %s\n",
			len(phrase.Selections), phrase.EntropyBits, strength(phrase.EntropyBits), phrase.Wordlist)
	}
	return 0
}

func randomIndex(random io.Reader, size int) (int, error) {
	n, err := rand.Int(random, big.NewInt(int64(size)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

// Strength labels are descriptive display buckets, not security guarantees.
func strength(bits float64) string {
	switch {
	case bits < 40:
		return "Low"
	case bits < 60:
		return "Moderate"
	case bits < 75:
		return "Strong"
	default:
		return "Very strong"
	}
}
