package util

import (
	"sort"
	"strings"
)

// RestoreAbstract converts an OpenAlex inverted index back to plain text.
func RestoreAbstract(idx map[string][]int) string {
	if idx == nil {
		return ""
	}
	type posWord struct {
		pos  int
		word string
	}
	var pairs []posWord
	for word, positions := range idx {
		for _, pos := range positions {
			pairs = append(pairs, posWord{pos, word})
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].pos < pairs[j].pos
	})
	var b strings.Builder
	for _, p := range pairs {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(p.word)
	}
	return b.String()
}
