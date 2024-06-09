package glob

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Pattern represents a globbing pattern that can be used to match against a
// given filepath.
//
// Unlike the builtin [filepath.Match] function, this supports both `**` for
// matching against arbitrary numbers of directories, as well as `!` for
// negating the match.
type Pattern string

// Match is a function that will match a given pattern against a path. This
// function will return a boolean value indicating whether the pattern matched
// the path, as well as an error if one occurred.
func (p Pattern) Match(path string) bool {
	status, err := p.match(string(p), path)
	return err == nil && status == statusMatched
}

// Filter returns a list of paths that match the pattern.
func (p Pattern) Filter(paths ...string) []string {
	var filtered []string
	for _, path := range paths {
		if p.Match(path) {
			filtered = append(filtered, path)
		}
	}
	return filtered
}

// Glob returns a list of paths that match the pattern by recursively searching
// elements in the base directory.
func (p Pattern) Glob(base string) []string {
	var paths []string
	pattern := p.Prepend(base)
	_ = filepath.Walk(base, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if pattern.Match(path) {
			paths = append(paths, path)
		}
		return nil
	})
	return paths
}

// Abs returns this pattern that was made absolute. If the pattern is already
// absolute, it will return the pattern as is. Any negations are semantically
// replaced so that it continues to work as expected.
func (p Pattern) Abs() (Pattern, error) {
	trimmed := strings.TrimLeft(string(p), "!")
	if filepath.IsAbs(trimmed) {
		return p, nil
	}
	abs, err := filepath.Abs(trimmed)
	if err != nil {
		return "", err
	}
	offset := len(p) - len(trimmed)
	prefix := string(p[:offset])
	return Pattern(prefix + abs), nil
}

// Prepend returns this pattern that prepends the given base, if the pattern
// is not from an absolute path. If the pattern is already absolute, it will
// return the pattern as is. Any negations are semantically replaced so that it
// continues to work as expected.
func (p Pattern) Prepend(base string) Pattern {
	trimmed := strings.TrimLeft(string(p), "!")
	if filepath.IsAbs(trimmed) {
		return p
	}
	offset := len(p) - len(trimmed)
	prefix := string(p[:offset])
	return Pattern(prefix + filepath.Join(base, trimmed))
}

// String converts this pattern to a string.
func (p Pattern) String() string {
	return string(p)
}

var _ fmt.Stringer = (*Pattern)(nil)

type status int

const (
	statusMatched status = iota
	statusRejected
	statusUnmatched
)

// match performs the internal matching, and allows for explicit rejection of
// matching paths.
func (p Pattern) match(pattern, path string) (status, error) {
	if strings.HasPrefix(pattern, "!") {
		return p.matchInverse(pattern[1:], path)
	}
	matched, err := p.matchParts(pattern, path)
	if err != nil {
		return statusUnmatched, err
	}
	if matched {
		return statusMatched, nil
	}
	return statusUnmatched, nil
}

func (p Pattern) matchParts(pattern, path string) (bool, error) {
	patternParts := strings.Split(pattern, string(filepath.Separator))
	pathParts := strings.Split(path, string(filepath.Separator))

	return p.matchRecursive(patternParts, pathParts)
}

func (p Pattern) matchRecursive(patternParts, pathParts []string) (bool, error) {
	pathIdx := 0
	for patternIdx, patternPart := range patternParts {
		if patternPart == "**" {
			if patternIdx < len(patternParts) {
				matched, err := p.matchRecursive(patternParts[patternIdx+1:], pathParts[pathIdx:])
				if matched || err != nil {
					return matched, err
				}
			}
			if pathIdx < len(pathParts) {
				matched, err := p.matchRecursive(patternParts[patternIdx:], pathParts[pathIdx+1:])
				if matched || err != nil {
					return matched, err
				}
			}
			return false, nil
		}
		if pathIdx >= len(pathParts) {
			return false, nil
		}
		matched, err := filepath.Match(patternPart, pathParts[pathIdx])
		if !matched || err != nil {
			return false, err
		}
		pathIdx++
	}
	return pathIdx == len(pathParts), nil
}

func (p Pattern) matchInverse(pattern string, path string) (status, error) {
	status, err := p.match(pattern, path)
	if err != nil {
		return statusUnmatched, err
	}
	if status == statusMatched {
		return statusRejected, nil
	}
	return status, nil
}
