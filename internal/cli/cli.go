package cli

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"text/template"

	"github.com/bitwizeshift/protobuild/internal/ansi"
	"github.com/spf13/cobra"
)

var (
	appName atomic.Pointer[string]
)

func init() {
	name := new(string)
	*name = filepath.Base(defaultAppName())
	appName.Store(name)
}

func defaultAppName() string {
	if len(os.Args[0]) == 0 {
		return "cli"
	}
	return filepath.Base(os.Args[0])
}

func AppName() string {
	name := appName.Load()
	return *name
}

func SetAppName(name string) {
	appName.Store(&name)
}

func ErrorPrefix() string {
	return ansi.Sprintf("%serror: %s%s:%s",
		ansi.FGRed,
		ansi.FGWhite,
		AppName(),
		ansi.Reset,
	)
}

func WarningPrefix() string {
	return ansi.Sprintf("%swarning:%s",
		ansi.FGYellow,
		ansi.Reset,
	)
}

func NoticePrefix() string {
	return ansi.Sprintf("%snotice:%s",
		ansi.FGCyan,
		ansi.Reset,
	)
}

var (
	//go:embed usage.template
	usageTemplate string

	//go:embed help.template
	helpTemplate string

	//go:embed version.template
	versionTemplate string

	//go:embed panic.template
	panicTemplate string
)

func SetDefaults(cmd *cobra.Command) {
	cmd.SetErrPrefix(ErrorPrefix())
	cmd.SetHelpTemplate(helpTemplate)
	cmd.SetUsageTemplate(usageTemplate)
	cmd.SetVersionTemplate(versionTemplate)
}

type payload struct {
	Error      string
	StackTrace []string
}

// HandlePanic is a function that can be deferred to provide a cleaner
// panic handler. This makes use of the `panic.template` file.
func HandlePanic() {
	if r := recover(); r != nil {
		tracePanic(os.Stderr, r)
		os.Exit(2)
	}
}

func tracePanic(wr io.Writer, r any) {
	template := template.Must(template.New("panic").Funcs(funcs).Parse(panicTemplate))

	stack := string(debug.Stack())

	payload := &payload{
		Error:      fmt.Sprintf("%v", r),
		StackTrace: strings.Split(stack, "\n"),
	}
	if err := template.Execute(wr, payload); err != nil {
		panic(err)
	}
	os.Exit(1)
}
