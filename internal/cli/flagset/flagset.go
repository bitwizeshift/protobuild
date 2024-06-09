package flagset

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bitwizeshift/protobuild/internal/ansi"
	"github.com/spf13/pflag"
)

// FlagSet is a wrapper around pflag.FlagSet that provides additional
// functionality for formatting flag usage, and encoding the proper name in a way
// that can be retrieved later.
type FlagSet struct {
	name string
	*pflag.FlagSet
}

// New creates a new FlagSet with the given name.
func New(name string) *FlagSet {
	return &FlagSet{
		name:    name,
		FlagSet: pflag.NewFlagSet(name, pflag.ExitOnError),
	}
}

// Name returns the name of the FlagSet.
func (f *FlagSet) Name() string {
	return f.name
}

// FlagUsages returns a string containing the usage information for the flags
// in the FlagSet.
func (f *FlagSet) FlagUsages() string {
	return f.FormattedFlagUsages(nil)
}

// FormatOptions is used to configure the formatting of the flag usage.
type FormatOptions struct {
	ArgFormat        ansi.Display
	FlagFormat       ansi.Display
	DeprecatedFormat ansi.Display
}

func (fo *FormatOptions) argFormat() ansi.Display {
	if fo == nil || fo.ArgFormat == nil {
		return ansi.None
	}
	return fo.ArgFormat
}

func (fo *FormatOptions) flagFormat() ansi.Display {
	if fo == nil || fo.FlagFormat == nil {
		return ansi.None
	}
	return fo.FlagFormat
}

func (fo *FormatOptions) deprecatedFormat() ansi.Display {
	if fo == nil || fo.FlagFormat == nil {
		return ansi.None
	}
	return fo.DeprecatedFormat
}

// FormattedFlagUsages returns a string containing the usage information for
// the flags in the FlagSet, formatted according to the provided options.
func (f *FlagSet) FormattedFlagUsages(opts *FormatOptions) string {
	buf := new(bytes.Buffer)

	var lines []string

	maxlen := 0
	f.VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}

		line := ""
		if flag.Shorthand != "" && flag.ShorthandDeprecated == "" {
			line = opts.flagFormat().Format("  -%s, --%s", flag.Shorthand, flag.Name)
		} else {
			line = opts.flagFormat().Format("      --%s", flag.Name)
		}

		varname, usage := pflag.UnquoteUsage(flag)
		if varname != "" {
			line += " " + opts.argFormat().Format(varname)
		}

		// This special character will be replaced with spacing once the
		// correct alignment is calculated
		line += "\x00"
		maxlen = max(maxlen, len(ansi.StripFormat(line)))

		line += usage
		if len(flag.Deprecated) != 0 {
			line += opts.deprecatedFormat().Format(" (DEPRECATED: %s)", flag.Deprecated)
		}

		lines = append(lines, line)
	})

	for _, line := range lines {
		sidx := strings.Index(ansi.StripFormat(line), "\x00")
		sidx2 := strings.Index(line, "\x00")

		spacing := strings.Repeat(" ", maxlen-sidx)
		// maxlen + 2 comes from + 1 for the \x00 and + 1 for the (deliberate) off-by-one in maxlen-sidx
		fmt.Fprintln(buf, line[:sidx2], spacing, wrap(maxlen+2, 0, line[sidx2+1:]))
	}

	return buf.String()
}

// Wraps the string `s` to a maximum width `w` with leading indent
// `i`. The first line is not indented (this is assumed to be done by
// caller). Pass `w` == 0 to do no wrapping
func wrap(i, w int, s string) string {
	if w == 0 {
		return strings.Replace(s, "\n", "\n"+strings.Repeat(" ", i), -1)
	}

	// space between indent i and end of line width w into which
	// we should wrap the text.
	wrap := w - i

	var r, l string

	// Not enough space for sensible wrapping. Wrap as a block on
	// the next line instead.
	if wrap < 24 {
		i = 16
		wrap = w - i
		r += "\n" + strings.Repeat(" ", i)
	}
	// If still not enough space then don't even try to wrap.
	if wrap < 24 {
		return strings.Replace(s, "\n", r, -1)
	}

	// Try to avoid short orphan words on the final line, by
	// allowing wrapN to go a bit over if that would fit in the
	// remainder of the line.
	slop := 5
	wrap = wrap - slop

	// Handle first line, which is indented by the caller (or the
	// special case above)
	l, s = wrapN(wrap, slop, s)
	r = r + strings.Replace(l, "\n", "\n"+strings.Repeat(" ", i), -1)

	// Now wrap the rest
	for s != "" {
		var t string

		t, s = wrapN(wrap, slop, s)
		r = r + "\n" + strings.Repeat(" ", i) + strings.Replace(t, "\n", "\n"+strings.Repeat(" ", i), -1)
	}
	return r
}

func wrapN(i, slop int, s string) (string, string) {
	if i+slop > len(s) {
		return s, ""
	}

	w := strings.LastIndexAny(s[:i], " \t\n")
	if w <= 0 {
		return s, ""
	}
	nlPos := strings.LastIndex(s[:i], "\n")
	if nlPos > 0 && nlPos < w {
		return s[:nlPos], s[nlPos+1:]
	}
	return s[:w], s[w+1:]
}
