/*
Package env provides a way to get environment variables and falling back to
default values if they are not set.
*/
package env

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"github.com/bitwizeshift/protobuild/internal/sys"
)

// ErrNoProtobuildPath is returned when the protobuild path could not be
// determined.
var ErrNoProtobuildPath = fmt.Errorf("protobuild path")

// ConfigPath returns the path to the protobuild directory.
func ConfigPath() (string, error) {
	if path := os.Getenv("PROTOBUILD_PATH"); path != "" {
		return path, nil
	}
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".protobuild"), nil
	}
	if build.Default.GOOS == "windows" {
		if userprofile := os.Getenv("USERPROFILE"); userprofile != "" {
			return filepath.Join(userprofile, ".protobuild"), nil
		}
	}
	dirname, err := sys.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrNoProtobuildPath, err)
	}
	return filepath.Join(dirname, ".protobuild"), nil
}

// BinPath the path to where protobuf binaries are stored.
func BinPath() (string, error) {
	return configPath("PROTOBUILD_BIN", "bin")
}

// RegistryPath the path to where the protobuild registry is stored.
func RegistryPath() (string, error) {
	return configPath("PROTOBUILD_REGISTRY", "registry")
}

// CachePath the path to where the protobuild cache is stored.
func CachePath() (string, error) {
	return configPath("PROTOBUILD_CACHE", "cache")
}

func configPath(env, subpath string) (string, error) {
	if path := os.Getenv(env); path != "" {
		return path, nil
	}
	path, err := ConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(path, subpath), nil
}

// Must returns value if err is nil, otherwise panics with the error.
func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
