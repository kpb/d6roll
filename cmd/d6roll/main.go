package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kpb/d6roll"
)

const version = "0.2.0"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, d6roll.Generate))
}

func run(args []string, stdout, stderr io.Writer, generate func(int) (d6roll.Passphrase, error)) int {
	flags := flag.NewFlagSet("d6roll", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var count int
	var separator string
	var verbose, rolls, help, showVersion bool
	flags.IntVar(&count, "w", 6, "number of words (3–32)")
	flags.IntVar(&count, "words", 6, "number of words (3–32)")
	flags.StringVar(&separator, "s", " ", "word separator")
	flags.StringVar(&separator, "separator", " ", "word separator")
	flags.BoolVar(&verbose, "v", false, "show metadata on stderr")
	flags.BoolVar(&verbose, "verbose", false, "show metadata on stderr")
	flags.BoolVar(&rolls, "rolls", false, "show Diceware rolls and metadata on stderr")
	flags.BoolVar(&help, "h", false, "show help")
	flags.BoolVar(&help, "help", false, "show help")
	flags.BoolVar(&showVersion, "V", false, "show version")
	flags.BoolVar(&showVersion, "version", false, "show version")
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: d6roll [flags]")
		fmt.Fprintln(flags.Output(), "       d6roll help")
		flags.PrintDefaults()
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
	if _, err := fmt.Fprintln(stdout, strings.Join(phrase.Words(), separator)); err != nil {
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
