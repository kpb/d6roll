// Package d6roll generates passphrases from the embedded EFF long Diceware
// wordlist. Entropy assumes independent uniform selections from a secure source.
package d6roll

import (
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"math/big"
)

// Selection records a chosen word and its Diceware roll. Roll is empty if
// the wordlist does not provide roll notation; a roll adds no entropy.
type Selection struct {
	Word string
	Roll string
}

// Passphrase contains selected words and generation metadata. EntropyBits
// assumes a secure random source and uniform independent selection.
type Passphrase struct {
	Selections  []Selection
	EntropyBits float64
	Wordlist    string
}

// Words returns a new slice containing the selected words in order.
func (p Passphrase) Words() []string {
	words := make([]string, len(p.Selections))
	for i, selection := range p.Selections {
		words[i] = selection.Word
	}
	return words
}

// Generator generates passphrases with a specified random source. Its zero
// value is not usable; use NewGenerator to construct one.
type Generator struct {
	random io.Reader
	list   *wordlist
}

// NewGenerator uses random for word selection. In production, random must be
// cryptographically secure; nil is rejected. Use Generate for secure defaults.
func NewGenerator(random io.Reader) (*Generator, error) {
	if random == nil {
		return nil, fmt.Errorf("d6roll: nil random source")
	}
	return &Generator{random: random, list: &effLarge}, nil
}

// Generate selects n words using crypto/rand.Reader. n must be positive.
func Generate(n int) (Passphrase, error) {
	return (&Generator{random: rand.Reader, list: &effLarge}).Generate(n)
}

// GenerateWords selects n words using crypto/rand.Reader.
func GenerateWords(n int) ([]string, error) {
	phrase, err := Generate(n)
	if err != nil {
		return nil, err
	}
	return phrase.Words(), nil
}

// Generate selects n words independently, with replacement. It returns no
// partial passphrase on failure. n must be positive.
func (g *Generator) Generate(n int) (Passphrase, error) {
	if n <= 0 {
		return Passphrase{}, fmt.Errorf("d6roll: word count must be positive")
	}
	if g == nil || g.random == nil || g.list == nil {
		return Passphrase{}, fmt.Errorf("d6roll: uninitialized generator")
	}
	if len(g.list.entries) == 0 {
		return Passphrase{}, fmt.Errorf("d6roll: empty wordlist")
	}
	selections := make([]Selection, 0, n)
	for range n {
		index, err := selectIndex(g.random, len(g.list.entries))
		if err != nil {
			return Passphrase{}, fmt.Errorf("d6roll: select word: %w", err)
		}
		selections = append(selections, g.list.entries[index])
	}
	return Passphrase{
		Selections:  selections,
		EntropyBits: float64(n) * math.Log2(float64(len(g.list.entries))),
		Wordlist:    g.list.name,
	}, nil
}

// GenerateWords selects n words with the generator's random source.
func (g *Generator) GenerateWords(n int) ([]string, error) {
	phrase, err := g.Generate(n)
	if err != nil {
		return nil, err
	}
	return phrase.Words(), nil
}

func selectIndex(random io.Reader, size int) (int, error) {
	if size <= 0 {
		return 0, fmt.Errorf("d6roll: empty wordlist")
	}
	index, err := rand.Int(random, big.NewInt(int64(size)))
	if err != nil {
		return 0, err
	}
	return int(index.Int64()), nil
}
