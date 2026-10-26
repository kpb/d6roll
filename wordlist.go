package d6roll

import (
	"bufio"
	_ "embed"
	"fmt"
	"strings"
)

// The original EFF large wordlist: https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt
//
//go:embed wordlists/eff_large.txt
var effLargeData string

type wordlist struct {
	name    string
	entries []Selection
}

var effLarge = mustParseEFFWordlist(effLargeData)

func mustParseEFFWordlist(data string) wordlist {
	list, err := parseEFFWordlist(data)
	if err != nil {
		panic(fmt.Sprintf("d6roll: invalid bundled EFF wordlist: %v", err))
	}
	return list
}

func parseEFFWordlist(data string) (wordlist, error) {
	list := wordlist{name: "eff-large"}
	words := make(map[string]bool, 7776)
	rolls := make(map[string]bool, 7776)
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		if len(fields) != 2 || !validRoll(fields[0]) || fields[1] == "" || strings.TrimSpace(fields[1]) != fields[1] {
			return wordlist{}, fmt.Errorf("invalid entry at line %d", len(list.entries)+1)
		}
		roll, word := fields[0], fields[1]
		if words[word] || rolls[roll] {
			return wordlist{}, fmt.Errorf("duplicate word or roll at line %d", len(list.entries)+1)
		}
		words[word], rolls[roll] = true, true
		list.entries = append(list.entries, Selection{Word: word, Roll: roll})
	}
	if err := scanner.Err(); err != nil {
		return wordlist{}, fmt.Errorf("read wordlist: %w", err)
	}
	if len(list.entries) != 7776 {
		return wordlist{}, fmt.Errorf("expected 7776 entries, got %d", len(list.entries))
	}
	return list, nil
}

func validRoll(roll string) bool {
	if len(roll) != 5 {
		return false
	}
	for _, digit := range roll {
		if digit < '1' || digit > '6' {
			return false
		}
	}
	return true
}
