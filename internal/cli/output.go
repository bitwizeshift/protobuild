package cli

import (
	"os"

	"github.com/bitwizeshift/protobuild/internal/ansi"
)

func Error(args ...any) {
	prefix := ErrorPrefix() + " "
	args = append([]any{prefix}, args...)
	args = append(args, "\n")
	ansi.Fprint(os.Stderr, args...)
}

func Errorf(format string, args ...any) {
	Error(ansi.Sprintf(format, args...))
}

func Warning(args ...any) {
	prefix := WarningPrefix() + " "
	args = append([]any{prefix}, args...)
	args = append(args, "\n")
	ansi.Fprint(os.Stderr, args...)
}

func Warningf(format string, args ...any) {
	Warning(ansi.Sprintf(format, args...))
}

func Notice(args ...any) {
	prefix := NoticePrefix() + " "
	args = append([]any{prefix}, args...)
	args = append(args, "\n")
	ansi.Fprint(os.Stderr, args...)
}

func Noticef(format string, args ...any) {
	Notice(ansi.Sprintf(format, args...))
}

func Fatal(args ...any) {
	Error(args...)
	os.Exit(1)
}

func Fatalf(format string, args ...any) {
	Errorf(format, args...)
	os.Exit(1)
}
