package pattern

import (
	"fmt"
)

// Parse extracts parameter names from a URL path pattern with braces.
//
// It identifies all parameters within braces like {name} or {name:regexp}.
func Parse(pattern string) ([]string, error) {
	pairs, err := bracePairs(pattern)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		for colon := pairs[i] + 1; colon <= pairs[i+1]; colon++ {
			if pattern[colon] == ':' || colon == pairs[i+1] {
				names = append(names, pattern[pairs[i]+1:colon])
				break
			}
		}
	}

	return names, nil
}

// bracePairs finds all first matching pairs of braces in a string.
//
// It returns an error in case of unbalanced braces.
func bracePairs(s string) ([]int, error) {
	pairs, level := make([]int, 0, 16), 0
	for i, left := 0, 0; i < len(s) && level >= 0; i++ {
		switch s[i] {
		case '{':
			if level++; level == 1 {
				left = i
			}
		case '}':
			if level--; level == 0 {
				pairs = append(pairs, left, i)
			}
		}
	}
	if level != 0 {
		return nil, fmt.Errorf("alchemy: unbalanced braces in %q", s)
	}

	return pairs, nil
}
