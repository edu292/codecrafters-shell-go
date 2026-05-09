package trie

import (
	"reflect"
	"sort"
	"testing"
)

func TestTrie(t *testing.T) {
	tr := NewTrie()

	// data setup
	words := []string{"apple", "app", "apply", "bat", "batch"}
	for _, w := range words {
		tr.Insert(w)
	}

	// 1. check IsWord
	t.Run("IsWord", func(t *testing.T) {
		tests := []struct {
			word string
			want bool
		}{
			{"apple", true},
			{"app", true},
			{"appl", false}, // prefix only
			{"bat", true},
			{"cat", false}, // missing
			{"", false},    // empty
		}

		for _, tc := range tests {
			if got := tr.IsWord(tc.word); got != tc.want {
				t.Errorf("IsWord(%q) = %v; want %v", tc.word, got, tc.want)
			}
		}
	})

	// 2. check WithPrefix
	t.Run("WithPrefix", func(t *testing.T) {
		tests := []struct {
			prefix string
			want   []string
		}{
			{"app", []string{"app", "apple", "apply"}},
			{"ba", []string{"bat", "batch"}},
			{"cat", []string{}}, // or []string{} depending on impl
			{"", []string{"app", "apple", "apply", "bat", "batch"}},
		}

		for _, tc := range tests {
			got := tr.WithPrefix(tc.prefix)

			// sort both -> order-independent compare
			sort.Strings(got)
			sort.Strings(tc.want)

			if len(got) == 0 && len(tc.want) == 0 {
				continue
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("WithPrefix(%q) = %v; want %v", tc.prefix, got, tc.want)
			}
		}
	})
}
