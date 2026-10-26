package d6roll

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestEFFWordlistIntegrity(t *testing.T) {
	if got := fmt.Sprintf("%x", sha256.Sum256([]byte(effLargeData))); got != "addd35536511597a02fa0a9ff1e5284677b8883b83e986e43f15a3db996b903e" {
		t.Fatalf("EFF source checksum changed: %s", got)
	}
	list, err := parseEFFWordlist(effLargeData)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(list.entries); got != 7776 {
		t.Fatalf("entries = %d, want 7776", got)
	}
	for _, tc := range []struct{ roll, word string }{
		{"11111", "abacus"}, {"14352", "bullhorn"},
		{"14565", "canyon"}, {"66666", "zoom"},
	} {
		found := false
		for _, entry := range list.entries {
			if entry.Roll == tc.roll {
				if entry.Word != tc.word {
					t.Errorf("roll %s = %q, want %q", tc.roll, entry.Word, tc.word)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing roll %s", tc.roll)
		}
	}
}

func TestRejectInvalidEFFData(t *testing.T) {
	for _, tc := range []struct{ name, data string }{
		{"short", strings.TrimSuffix(effLargeData, "66666\tzoom\n")},
		{"empty word", strings.Replace(effLargeData, "11111\tabacus", "11111\t", 1)},
		{"duplicate word", strings.Replace(effLargeData, "11112\tabdomen", "11112\tabacus", 1)},
		{"duplicate roll", strings.Replace(effLargeData, "11112\tabdomen", "11111\tabdomen", 1)},
		{"invalid digit", strings.Replace(effLargeData, "11111\tabacus", "11110\tabacus", 1)},
		{"invalid length", strings.Replace(effLargeData, "11111\tabacus", "1111\tabacus", 1)},
		{"blank line", "\n" + effLargeData},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := parseEFFWordlist(tc.data); err == nil {
				t.Fatal("accepted invalid wordlist")
			}
		})
	}
}

func TestDeterministicGeneration(t *testing.T) {
	// crypto/rand.Int reads two bytes per EFF index; indices 0, 1, 0.
	g, err := NewGenerator(bytes.NewReader([]byte{0, 0, 0, 1, 0, 0}))
	if err != nil {
		t.Fatal(err)
	}
	p, err := g.Generate(3)
	if err != nil {
		t.Fatal(err)
	}
	want := []Selection{{"abacus", "11111"}, {"abdomen", "11112"}, {"abacus", "11111"}}
	if !reflect.DeepEqual(p.Selections, want) {
		t.Fatalf("selections = %#v, want %#v", p.Selections, want)
	}
	if p.Wordlist != "eff-large" || math.Abs(p.EntropyBits-3*math.Log2(7776)) > 1e-10 {
		t.Fatalf("metadata = %+v", p)
	}
	words := p.Words()
	if !reflect.DeepEqual(words, []string{"abacus", "abdomen", "abacus"}) {
		t.Fatalf("words = %q", words)
	}
	words[0] = "changed"
	if p.Selections[0].Word != "abacus" {
		t.Fatal("Words exposed the selection slice")
	}
}

func TestGenerateWordsAndValidation(t *testing.T) {
	g, err := NewGenerator(bytes.NewReader([]byte{0, 0, 0, 1, 0, 2}))
	if err != nil {
		t.Fatal(err)
	}
	words, err := g.GenerateWords(3)
	if err != nil || !reflect.DeepEqual(words, []string{"abacus", "abdomen", "abdominal"}) {
		t.Fatalf("GenerateWords = %q, %v", words, err)
	}
	if _, err := NewGenerator(nil); err == nil {
		t.Fatal("accepted nil reader")
	}
	for _, count := range []int{-1, 0, maxWords + 1, math.MaxInt} {
		if _, err := g.Generate(count); err == nil {
			t.Fatalf("accepted count %d", count)
		}
		if _, err := Generate(count); err == nil {
			t.Fatalf("Generate accepted count %d", count)
		}
		if _, err := GenerateWords(count); err == nil {
			t.Fatalf("GenerateWords accepted count %d", count)
		}
		if _, err := g.GenerateWords(count); err == nil {
			t.Fatalf("generator.GenerateWords accepted count %d", count)
		}
	}
	if _, err := (&Generator{}).Generate(1); err == nil {
		t.Fatal("accepted zero-value generator")
	}
	if _, err := (*Generator)(nil).Generate(1); err == nil {
		t.Fatal("accepted nil generator")
	}
	if p, err := Generate(1); err != nil || len(p.Selections) != 1 {
		t.Fatalf("Generate(1) failed: %v", err)
	}
	if words, err := GenerateWords(1); err != nil || len(words) != 1 {
		t.Fatalf("GenerateWords(1) failed: %v", err)
	}
}

func TestLibraryMaximumWordCount(t *testing.T) {
	g, err := NewGenerator(bytes.NewReader(make([]byte, maxWords*2)))
	if err != nil {
		t.Fatal(err)
	}
	p, err := g.Generate(maxWords)
	if err != nil || len(p.Selections) != maxWords {
		t.Fatalf("Generate(%d) returned %d selections, %v", maxWords, len(p.Selections), err)
	}
}

type failedReader struct{ err error }

func (r failedReader) Read([]byte) (int, error) { return 0, r.err }

type panicReader struct{}

func (panicReader) Read([]byte) (int, error) { panic("invalid count consumed randomness") }

func TestOversizedCountDoesNotReadRandomness(t *testing.T) {
	g, err := NewGenerator(panicReader{})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Generate(maxWords + 1); err == nil {
		t.Fatal("accepted oversized count")
	}
}

func TestRandomFailureReturnsNoPhrase(t *testing.T) {
	wantErr := errors.New("source failed")
	random := io.MultiReader(bytes.NewReader([]byte{0, 0}), failedReader{wantErr})
	g, err := NewGenerator(random)
	if err != nil {
		t.Fatal(err)
	}
	p, err := g.Generate(2)
	if !errors.Is(err, wantErr) || len(p.Selections) != 0 {
		t.Fatalf("Generate = %+v, %v", p, err)
	}
}

func TestUniformIndexSelection(t *testing.T) {
	// 3 is outside [0,3), so rejection sampling must consume the next byte.
	index, err := selectIndex(bytes.NewReader([]byte{3, 2}), 3)
	if err != nil || index != 2 {
		t.Fatalf("selectIndex = %d, %v", index, err)
	}
	if _, err := selectIndex(bytes.NewReader(nil), 0); err == nil {
		t.Fatal("accepted empty wordlist")
	}
	if _, err := selectIndex(bytes.NewReader([]byte{3}), 3); !errors.Is(err, io.EOF) {
		t.Fatalf("missing rejection/read failure: %v", err)
	}
}

func TestArbitraryWordlistSizeAndMissingRoll(t *testing.T) {
	g := &Generator{random: bytes.NewReader([]byte{2, 2}), list: &wordlist{
		name: "test", entries: []Selection{{Word: "a"}, {Word: "b"}, {Word: "c"}},
	}}
	p, err := g.Generate(2)
	if err != nil || !reflect.DeepEqual(p.Words(), []string{"c", "c"}) || p.Selections[0].Roll != "" {
		t.Fatalf("Generate = %+v, %v", p, err)
	}
	if math.Abs(p.EntropyBits-2*math.Log2(3)) > 1e-10 {
		t.Fatalf("entropy = %v", p.EntropyBits)
	}
}
