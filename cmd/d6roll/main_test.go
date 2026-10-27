package main

import (
	"bytes"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/kpb/d6roll"
)

var sample = []d6roll.Selection{
	{Word: "abacus", Roll: "11111"}, {Word: "abdomen", Roll: "11112"},
	{Word: "abdominal", Roll: "11113"}, {Word: "abide", Roll: "11114"},
	{Word: "abiding", Roll: "11115"}, {Word: "ability", Roll: "11116"},
}

func TestMemorableFormatting(t *testing.T) {
	for _, tc := range []struct {
		name   string
		args   []string
		random []byte
		want   string
	}{
		{
			name: "short alias",
			args: []string{"-m"}, random: []byte{0, 0, 0},
			want: "Abacus0!-Abdomen-Abdominal-Abide-Abiding-Ability\n",
		},
		{
			name: "long alias",
			args: []string{"--memorable"}, random: []byte{0, 0, 0},
			want: "Abacus0!-Abdomen-Abdominal-Abide-Abiding-Ability\n",
		},
		{
			name: "interior word",
			args: []string{"-m"}, random: []byte{2, 4, 1},
			want: "Abacus-Abdomen-Abdominal4@-Abide-Abiding-Ability\n",
		},
		{
			name: "last word and maximum characters",
			args: []string{"-m"}, random: []byte{5, 9, 7},
			want: "Abacus-Abdomen-Abdominal-Abide-Abiding-Ability9*\n",
		},
		{
			name: "rejected out-of-range draws",
			// 7 is outside the six-word range; 15 is outside the digit range.
			args: []string{"-m"}, random: []byte{7, 5, 15, 9, 7},
			want: "Abacus-Abdomen-Abdominal-Abide-Abiding-Ability9*\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			status := runWithRandom(tc.args, &stdout, &stderr, fakeGenerate, bytes.NewReader(tc.random))
			if status != 0 || stdout.String() != tc.want || stderr.Len() != 0 {
				t.Fatalf("status=%d stdout=%q stderr=%q", status, stdout.String(), stderr.String())
			}
		})
	}
}

func TestMemorableSeparatorOverrides(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"custom", []string{"--memorable", "--separator", "_"}, "Abacus0!_Abdomen_Abdominal\n"},
		{"short space", []string{"-m", "-s", " "}, "Abacus0! Abdomen Abdominal\n"},
		{"long space before preset", []string{"--separator= ", "--memorable"}, "Abacus0! Abdomen Abdominal\n"},
		{"short empty before preset", []string{"-s", "", "-m"}, "Abacus0!AbdomenAbdominal\n"},
		{"long empty", []string{"--memorable", "--separator="}, "Abacus0!AbdomenAbdominal\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			args := append([]string{"-w", "3"}, tc.args...)
			status := runWithRandom(args, &stdout, &stderr, fakeGenerate, bytes.NewReader([]byte{0, 0, 0}))
			if status != 0 || stdout.String() != tc.want || stderr.Len() != 0 {
				t.Fatalf("status=%d stdout=%q stderr=%q", status, stdout.String(), stderr.String())
			}
		})
	}
}

func TestMemorableRandomFailureDoesNotOutputPassword(t *testing.T) {
	for _, tc := range []struct {
		name   string
		random []byte
	}{
		{"word selection", nil},
		{"digit selection", []byte{2}},
		{"symbol selection", []byte{2, 4}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			status := runWithRandom([]string{"-m", "--verbose", "--rolls"}, &stdout, &stderr, fakeGenerate, bytes.NewReader(tc.random))
			if status != 1 || stdout.Len() != 0 || stderr.String() != "d6roll: memorable formatting failed\n" {
				t.Fatalf("status=%d stdout=%q stderr=%q", status, stdout.String(), stderr.String())
			}
		})
	}
}

func fakeGenerate(n int) (d6roll.Passphrase, error) {
	return d6roll.Passphrase{
		Selections: sample[:n], EntropyBits: float64(n) * math.Log2(7776), Wordlist: "eff-large",
	}, nil
}

func output(args []string, generate func(int) (d6roll.Passphrase, error)) (int, string, string) {
	var stdout, stderr bytes.Buffer
	status := run(args, &stdout, &stderr, generate)
	return status, stdout.String(), stderr.String()
}

func TestDefaultOutput(t *testing.T) {
	status, stdout, stderr := output(nil, fakeGenerate)
	if status != 0 || stdout != "abacus abdomen abdominal abide abiding ability\n" || stderr != "" {
		t.Fatalf("status=%d stdout=%q stderr=%q", status, stdout, stderr)
	}
}

func TestExplicitHyphenSeparator(t *testing.T) {
	for _, args := range [][]string{{"-s", "-"}, {"--separator", "-"}} {
		status, stdout, stderr := output(args, fakeGenerate)
		if status != 0 || stdout != "abacus-abdomen-abdominal-abide-abiding-ability\n" || stderr != "" {
			t.Fatalf("args=%q status=%d stdout=%q stderr=%q", args, status, stdout, stderr)
		}
	}
}

func TestFlagAliasesAndSeparators(t *testing.T) {
	for _, args := range [][]string{
		{"-w", "3", "-s", " "}, {"--words=3", "--separator= "},
	} {
		status, stdout, stderr := output(args, fakeGenerate)
		if status != 0 || stdout != "abacus abdomen abdominal\n" || stderr != "" {
			t.Fatalf("args=%q status=%d stdout=%q stderr=%q", args, status, stdout, stderr)
		}
	}
}

func TestMaximumCLIWordCount(t *testing.T) {
	status, stdout, stderr := output([]string{"-w=32"}, func(n int) (d6roll.Passphrase, error) {
		if n != 32 {
			t.Fatalf("count = %d", n)
		}
		selections := make([]d6roll.Selection, n)
		for i := range selections {
			selections[i].Word = "abacus"
		}
		return d6roll.Passphrase{Selections: selections}, nil
	})
	if status != 0 || len(strings.Split(strings.TrimSpace(stdout), " ")) != 32 || stderr != "" {
		t.Fatalf("status=%d stdout=%q stderr=%q", status, stdout, stderr)
	}
}

func TestVerboseAndRollsAreOnlyOnStderr(t *testing.T) {
	for _, args := range [][]string{{"-v"}, {"--verbose"}, {"--rolls"}} {
		status, stdout, stderr := output(args, fakeGenerate)
		if status != 0 || stdout != "abacus abdomen abdominal abide abiding ability\n" {
			t.Fatalf("args=%q status=%d stdout=%q", args, status, stdout)
		}
		for _, want := range []string{"Words:     6", "Entropy:   77.5 bits", "Strength:  Very strong", "Wordlist:  eff-large"} {
			if !strings.Contains(stderr, want) {
				t.Fatalf("args=%q stderr missing %q: %q", args, want, stderr)
			}
		}
		if args[0] == "--rolls" {
			if !strings.Contains(stderr, "Roll     Word") || !strings.Contains(stderr, "11112    abdomen") {
				t.Fatalf("roll table missing: %q", stderr)
			}
		} else if strings.Contains(stderr, "11112") {
			t.Fatalf("verbose showed rolls: %q", stderr)
		}
	}
}

func TestHelpAndVersion(t *testing.T) {
	status, stdout, stderr := output([]string{"help"}, nil)
	if status != 0 || stderr != "" || !strings.Contains(stdout, "Usage: d6roll") || !strings.Contains(stdout, "d6roll help") {
		t.Fatalf("help command: status=%d stdout=%q stderr=%q", status, stdout, stderr)
	}
	for _, want := range []string{
		"-h, --help", "-V, --version", "--rolls", "-m, --memorable", "-s, --separator string",
		"-v, --verbose", "-w, --words int", `word separator (default " "; memorable: "-")`,
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("help missing %q: %q", want, stdout)
		}
	}
	if strings.Contains(stdout, "\n  -verbose") || strings.Contains(stdout, "\n  -words") {
		t.Fatalf("help showed non-POSIX long option spelling: %q", stdout)
	}
	for _, flag := range []string{"-h", "--help"} {
		status, stdout, stderr := output([]string{flag}, nil)
		if status != 0 || stdout != "" || !strings.Contains(stderr, "Usage: d6roll") {
			t.Fatalf("help: status=%d stdout=%q stderr=%q", status, stdout, stderr)
		}
	}
	for _, flag := range []string{"-V", "--version"} {
		status, stdout, stderr := output([]string{flag}, nil)
		if status != 0 || stdout != "d6roll 0.3.0\n" || stderr != "" {
			t.Fatalf("version: status=%d stdout=%q stderr=%q", status, stdout, stderr)
		}
	}
}

func TestInvalidArgsAndGenerationFailure(t *testing.T) {
	for _, args := range [][]string{
		{"--words=2"}, {"-w", "33"}, {"--words=oops"}, {"extra"}, {"help", "extra"}, {"-w", "3", "help"}, {"-w", "3", "extra"}, {"--unknown"},
	} {
		status, stdout, stderr := output(args, func(int) (d6roll.Passphrase, error) {
			t.Fatal("invalid args invoked generator")
			return d6roll.Passphrase{}, nil
		})
		if status == 0 || stdout != "" || stderr == "" {
			t.Fatalf("args=%q status=%d stdout=%q stderr=%q", args, status, stdout, stderr)
		}
	}
	status, stdout, stderr := output(nil, func(int) (d6roll.Passphrase, error) {
		return d6roll.Passphrase{}, errors.New("source failed")
	})
	if status == 0 || stdout != "" || !strings.Contains(stderr, "source failed") {
		t.Fatalf("failure: status=%d stdout=%q stderr=%q", status, stdout, stderr)
	}
}

func TestStrengthBuckets(t *testing.T) {
	for _, tc := range []struct {
		bits float64
		want string
	}{{39, "Low"}, {40, "Moderate"}, {60, "Strong"}, {75, "Very strong"}} {
		if got := strength(tc.bits); got != tc.want {
			t.Errorf("strength(%v) = %q, want %q", tc.bits, got, tc.want)
		}
	}
}
