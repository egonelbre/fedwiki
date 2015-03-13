package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

type Reference struct {
	Rune rune
	Name string
}

type ByRune []Reference

func (a ByRune) Len() int           { return len(a) }
func (a ByRune) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRune) Less(i, j int) bool { return a[i].Rune < a[j].Rune }

func main() {
	fromrune := make(map[rune]Reference, len(entity))
	for name, r := range entity {
		if !unicode.IsSymbol(r) && !unicode.IsPunct(r) {
			continue
		}

		// use the shorter name
		if cur, exists := fromrune[r]; exists {
			if len(name) < len(cur.Name) {
				fromrune[r] = Reference{r, name}
			}
			continue
		}
		fromrune[r] = Reference{r, name}
	}

	sorted := make([]Reference, 0, len(fromrune))
	for _, t := range fromrune {
		t.Name = strings.ToLower(t.Name)
		t.Name = strings.Trim(t.Name, ";")
		sorted = append(sorted, t)
	}
	sort.Sort(ByRune(sorted))

	file, err := os.Create("runename.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprintln(file, "// generated code")
	fmt.Fprintln(file, "package fedwiki")
	fmt.Fprintln(file)
	fmt.Fprintln(file, "var runename = map[rune]string {")
	for _, tok := range sorted {
		fmt.Fprintf(file, "\t'\\U%08X': \"%s\",\n", tok.Rune, tok.Name)
	}
	fmt.Fprintln(file, "}")
}
