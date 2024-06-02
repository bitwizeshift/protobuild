/*
Package sys provides replaceable abstractions for
*/
package sys

import (
	"os"
	"sync/atomic"
)

var (
	home atomic.Pointer[func() (string, error)]
)

func init() {
	fn := os.UserHomeDir
	home.Store(&fn)
}

// SetUserHomeDirFunc sets the function to be used to determine the home directory.
func SetUserHomeDirFunc(f func() (string, error)) {
	home.Store(&f)
}

// UserHomeDir returns the home directory for the current user.
func UserHomeDir() (string, error) {
	get := home.Load()
	return (*get)()
}
