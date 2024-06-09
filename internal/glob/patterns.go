package glob

import (
	"os"
	"path/filepath"
	"unsafe"
)

// Patterns represents a list of patterns that can be used to match against a
// given filepath.
// Any negative patterns in the list will take precedence for matching.
type Patterns []Pattern

// NewPatterns is a function that will create a new Patterns object from a
// variadic list of string patterns.
func NewPatterns(patterns ...string) Patterns {
	result := make([]Pattern, len(patterns))
	for i, pattern := range patterns {
		result[i] = Pattern(pattern)
	}
	return result
}

func transmutePatterns(p []string) Patterns {
	return *(*Patterns)(unsafe.Pointer(&p))
}

// Match is a function that will match a given pattern against a name. This
// function will return a boolean value indicating whether the pattern matched
// the name and was valid.
func (p Patterns) Match(name string) bool {
	result := statusUnmatched
	for _, pattern := range p {
		r, err := pattern.match(string(pattern), name)
		if err != nil {
			return false
		}
		if r == statusRejected {
			return false
		}
		if r == statusMatched {
			result = statusMatched
		}
	}
	return result == statusMatched
}

// Filter is a function that will filter a list of names against a list of patterns.
// This function will return a list of names that matched any of the patterns.
func (p Patterns) Filter(names ...string) []string {
	var result []string
	for _, name := range names {
		if p.Match(name) {
			result = append(result, name)
		}
	}
	return result
}

// Glob is a function that will walk a given base directory and return a list of
// paths that match any of the patterns.
func (p Patterns) Glob(base string) []string {
	var patterns Patterns
	for _, pattern := range p {
		patterns = append(patterns, pattern.Prepend(base))
	}

	var paths []string
	_ = filepath.Walk(base, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if patterns.Match(path) {
			paths = append(paths, path)
		}
		return nil
	})
	return paths
}

// Abs is a function that will make all patterns absolute. If the pattern is
// already absolute, it will return the pattern as is. Any negations are
// semantically replaced so that it continues to work as expected.
func (p Patterns) Abs() (Patterns, error) {
	var result Patterns
	for _, pattern := range p {
		abs, err := pattern.Abs()
		if err != nil {
			return nil, err
		}
		result = append(result, abs)
	}
	return result, nil
}

// Prepend returns all patterns with the specified base prepended to it.
// If the pattern is already absolute, it will return the pattern as is. Any
// negations are semantically replaced so that it continues to work as expected.
func (p Patterns) Prepend(base string) Patterns {
	result := make(Patterns, 0, len(p))
	for _, pattern := range p {
		result = append(result, pattern.Prepend(base))
	}
	return result
}
