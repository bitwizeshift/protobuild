package env_test

import (
	"fmt"
	"go/build"
	"path/filepath"
	"testing"

	"github.com/bitwizeshift/protobuild/internal/env"
	"github.com/bitwizeshift/protobuild/internal/sys"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMust_NilError_ReturnsValue(t *testing.T) {
	want := 42

	got := env.Must(want, nil)

	if got != want {
		t.Errorf("Must: got %d, want %d", got, want)
	}
}

func TestMust_Error_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Must: did not panic")
		}
	}()
	var zero int

	env.Must(zero, fmt.Errorf("test error"))
}

func TestConfigPath_EnvSet_ReturnsValue(t *testing.T) {
	want := "/opt/local/.protobuild"
	t.Setenv("PROTOBUILD_PATH", want)

	got, err := env.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath: unexpected error: %v", err)
	}

	if got != want {
		t.Errorf("ConfigPath: got %q, want %q", got, want)
	}
}

func defaultHome() string {
	if build.Default.GOOS == "windows" {
		return "C:\\Users\\test"
	}
	return "/home/test"
}

func TestConfigPath_HomeSet_ReturnsValue(t *testing.T) {
	path := defaultHome()
	t.Setenv("HOME", path)
	want := filepath.Join(path, ".protobuild")

	got, err := env.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath: unexpected error: %v", err)
	}

	if got != want {
		t.Errorf("ConfigPath: got %q, want %q", got, want)
	}
}

func TestConfigPath_UserProfileSet_ReturnsValue(t *testing.T) {
	if build.Default.GOOS != "windows" {
		t.Skip("skipping test on non-windows")
	}

	path := defaultHome()
	t.Setenv("USERPROFILE", path)
	want := filepath.Join(path, ".protobuild")

	got, err := env.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath: unexpected error: %v", err)
	}

	if got != want {
		t.Errorf("ConfigPath: got %q, want %q", got, want)
	}
}

func TestConfigPath_HomeReturnsValue_ReturnsValue(t *testing.T) {
	want := filepath.Join(defaultHome(), ".protobuild")
	sys.SetUserHomeDirFunc(func() (string, error) {
		return defaultHome(), nil
	})
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	got, err := env.ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath: unexpected error: %v", err)
	}

	if got != want {
		t.Errorf("ConfigPath: got %q, want %q", got, want)
	}
}

func TestConfigPath_HomeReturnsError_ReturnsError(t *testing.T) {
	want := fmt.Errorf("test error")
	sys.SetUserHomeDirFunc(func() (string, error) {
		return "", want
	})
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	_, err := env.ConfigPath()

	if got := err; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
		t.Fatalf("ConfigPath: got err %v, want %v", got, want)
	}
}

func TestPath_EnvSpecified_ReturnsValue(t *testing.T) {
	home := defaultHome()
	protobuildPath := filepath.Join(home, ".protobuild")
	testCases := []struct {
		name  string
		value string
		fn    func() (string, error)
		want  string
	}{
		{
			name:  "BinPath",
			value: "PROTOBUILD_BIN",
			fn:    env.BinPath,
			want:  filepath.Join(protobuildPath, "bin"),
		}, {
			name:  "RegistryPath",
			value: "PROTOBUILD_REGISTRY",
			fn:    env.RegistryPath,
			want:  filepath.Join(protobuildPath, "registry"),
		}, {
			name:  "CachePath",
			value: "PROTOBUILD_CACHE",
			fn:    env.CachePath,
			want:  filepath.Join(protobuildPath, "cache"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(tc.value, tc.want)

			got, err := tc.fn()
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", tc.name, err)
			}

			if got != tc.want {
				t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
			}
		})
	}
}

func TestPath_ConfigPathKnown_ReturnsValue(t *testing.T) {
	home := defaultHome()
	protobuildPath := filepath.Join(home, ".protobuild")
	testCases := []struct {
		name string
		fn   func() (string, error)
		want string
	}{
		{
			name: "BinPath",
			fn:   env.BinPath,
			want: filepath.Join(protobuildPath, "bin"),
		}, {
			name: "RegistryPath",
			fn:   env.RegistryPath,
			want: filepath.Join(protobuildPath, "registry"),
		}, {
			name: "CachePath",
			fn:   env.CachePath,
			want: filepath.Join(protobuildPath, "cache"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("HOME", home)
			t.Setenv("USERPROFILE", home)

			got, err := tc.fn()
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", tc.name, err)
			}

			if got != tc.want {
				t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
			}
		})
	}
}

func TestPath_ConfigPathErrors_ReturnsError(t *testing.T) {
	want := fmt.Errorf("test error")
	testCases := []struct {
		name string
		fn   func() (string, error)
	}{
		{
			name: "BinPath",
			fn:   env.BinPath,
		}, {
			name: "RegistryPath",
			fn:   env.RegistryPath,
		}, {
			name: "CachePath",
			fn:   env.CachePath,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("HOME", "")
			t.Setenv("USERPROFILE", "")
			sys.SetUserHomeDirFunc(func() (string, error) {
				return "", want

			})

			_, err := tc.fn()

			if got := err; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Fatalf("%s: got err %v, want %v", tc.name, got, want)
			}
		})
	}
}
