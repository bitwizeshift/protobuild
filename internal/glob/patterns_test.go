package glob_test

import (
	"path/filepath"
	"testing"

	"github.com/bitwizeshift/protobuild/internal/glob"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPatternsMatch(t *testing.T) {
	testCases := []struct {
		name    string
		pattern glob.Patterns
		path    string
		want    bool
	}{
		{
			name:    "simple match",
			pattern: glob.NewPatterns("foo"),
			path:    "foo",
			want:    true,
		}, {
			name:    "simple mismatch",
			pattern: glob.NewPatterns("foo"),
			path:    "bar",
			want:    false,
		}, {
			name:    "simple match with directory",
			pattern: glob.NewPatterns(filepath.Join("foo", "bar")),
			path:    filepath.Join("foo", "bar"),
			want:    true,
		}, {
			name:    "simple mismatch with directory",
			pattern: glob.NewPatterns(filepath.Join("foo", "bar")),
			path:    filepath.Join("foo", "baz"),
			want:    false,
		}, {
			name:    "simple match with wildcard",
			pattern: glob.NewPatterns(filepath.Join("foo", "*.foo")),
			path:    filepath.Join("foo", "bar.foo"),
			want:    true,
		}, {
			name:    "simple mismatch with wildcard",
			pattern: glob.NewPatterns(filepath.Join("foo", "*.foo")),
			path:    filepath.Join("bar", "foo.foo"),
			want:    false,
		}, {
			name:    "simple mismatch with wildcard from nesting",
			pattern: glob.NewPatterns(filepath.Join("foo", "*")),
			path:    filepath.Join("foo", "bar", "baz"),
			want:    false,
		}, {
			name:    "simple match with recursive wildcard",
			pattern: glob.NewPatterns(filepath.Join("foo", "**")),
			path:    filepath.Join("foo", "bar", "baz"),
			want:    true,
		}, {
			name:    "recursive wildcard with suffix",
			pattern: glob.NewPatterns(filepath.Join("foo", "**", "baz")),
			path:    filepath.Join("foo", "bar-1", "bar-2", "baz"),
			want:    true,
		}, {
			name:    "recursive wildcard matches base",
			pattern: glob.NewPatterns(filepath.Join("foo", "**", "baz")),
			path:    filepath.Join("foo", "baz"),
			want:    true,
		}, {
			name:    "recursive wildcard does not match base",
			pattern: glob.NewPatterns(filepath.Join("foo", "**", "baz")),
			path:    filepath.Join("foo", "bar"),
			want:    false,
		}, {
			name:    "simple negation",
			pattern: glob.NewPatterns(filepath.Join("!foo")),
			path:    filepath.Join("foo"),
			want:    false,
		}, {
			name:    "simple negation with directory",
			pattern: glob.NewPatterns(filepath.Join("!foo", "bar")),
			path:    filepath.Join("foo", "bar"),
			want:    false,
		}, {
			name:    "simple error",
			pattern: glob.NewPatterns("["),
			path:    "foo",
			want:    false,
		}, {
			name:    "recursive wildcard with error",
			pattern: glob.NewPatterns(filepath.Join("foo", "**", "[")),
			path:    filepath.Join("foo", "bar", "baz"),
			want:    false,
		}, {
			name:    "recursive wildcard with error in base",
			pattern: glob.NewPatterns(filepath.Join("foo", "**", "[")),
			path:    filepath.Join("foo", "baz"),
			want:    false,
		}, {
			name:    "negation error",
			pattern: glob.NewPatterns(filepath.Join("![")),
			path:    filepath.Join("foo"),
			want:    false,
		}, {
			name:    "negated and pattern matched returns false",
			pattern: glob.NewPatterns("*", "!foo"),
			path:    "foo",
			want:    false,
		}, {
			name:    "matches pattern but not negation",
			pattern: glob.NewPatterns("*", "!foo"),
			path:    "bar",
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.pattern.Match(tc.path)

			if got != tc.want {
				t.Errorf("Pattern.TryMatch(%s): want %v, got %v", tc.name, tc.want, got)
			}
		})
	}
}

func TestPatternsFilter(t *testing.T) {
	testCases := []struct {
		name     string
		patterns glob.Patterns
		paths    []string
		want     []string
	}{
		{
			name:     "simple match",
			patterns: glob.NewPatterns("foo"),
			paths:    []string{"foo", "bar", "baz"},
			want:     []string{"foo"},
		}, {
			name:     "simple mismatch",
			patterns: glob.NewPatterns("foo"),
			paths:    []string{"bar", "baz"},
			want:     nil,
		}, {
			name:     "match with glob",
			patterns: glob.NewPatterns("foo*"),
			paths:    []string{"foo", "foobar", "bar", "baz"},
			want:     []string{"foo", "foobar"},
		}, {
			name:     "match with directory",
			patterns: glob.NewPatterns(filepath.Join("foo", "bar")),
			paths:    []string{"foo", filepath.Join("foo", "bar"), filepath.Join("foo", "baz"), "bar"},
			want:     []string{filepath.Join("foo", "bar")},
		}, {
			name:     "match with recursive wildcard and negation",
			patterns: glob.NewPatterns(filepath.Join("foo", "**"), filepath.Join("!foo", "bar")),
			paths:    []string{"foo", filepath.Join("foo", "bar"), filepath.Join("foo", "bar", "baz"), "bar"},
			want:     []string{"foo", filepath.Join("foo", "bar", "baz")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.patterns.Filter(tc.paths...)

			if !cmp.Equal(got, tc.want, cmpopts.SortSlices(strless)) {
				t.Errorf("Pattern.Filter(%s): want %v, got %v", tc.name, tc.want, got)
			}
		})
	}
}

func TestPatternsGlob(t *testing.T) {
	want := []string{
		filepath.Join("testdata", "foo", "bar", "baz"),
	}
	pattern := glob.NewPatterns(
		filepath.Join("foo", "bar", "**"),
		filepath.Join("!foo", "bar"),
	)

	got := pattern.Glob("testdata")

	if !cmp.Equal(got, want, cmpopts.SortSlices(strless)) {
		t.Errorf("Pattern.Glob: want %v, got %v", want, got)
	}
}
