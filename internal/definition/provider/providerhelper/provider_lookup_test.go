package providerhelper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderLookupFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		fn     ProviderLookupFunc
		expect string
	}{
		{
			name:   "provider is unset",
			fn:     nil,
			expect: "not implemented",
		},
		{
			name: "provider failed",
			fn: func(ctx context.Context, s *State) error {
				return errors.New("failed")
			},
			expect: "failed",
		},
		{
			name: "provider success",
			fn: func(ctx context.Context, s *State) error {
				return nil
			},
			expect: "",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if actual := tc.fn.Do(context.Background(), &State{}); tc.expect != "" {
				require.EqualError(t, actual, tc.expect, "Must match the expected error")
			} else {
				require.NoError(t, actual, "Must not return an error")
			}
		})
	}
}

func TestNewDefaultProviderLookups(t *testing.T) {
	t.Parallel()

	// Not testing the functionality of each of the providers considering
	// these are being tested individually
	assert.Len(t, NewDefaultProviderLookups(), 3, "Must have three default looks returned")
}

func TestFileProviderLookup(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		create  bool
		content string
		expect  State
		errVal  string
	}{
		{
			name:    "file does not exist",
			create:  false,
			content: "",
			expect:  State{},
			errVal:  "file not found",
		},
		{
			name:    "no file contents",
			create:  true,
			content: "",
			expect:  State{},
			errVal:  "no file content",
		},
		{
			name:    "invalid json",
			create:  true,
			content: `{"auth_token":"aaa"`,
			expect:  State{},
			errVal:  "unexpected EOF",
		},
		{
			name:    "valid json, no values defined",
			create:  true,
			content: `{}`,
			expect:  State{},
			errVal:  "",
		},
		{
			name:    "unknown json fields set",
			create:  true,
			content: `{"provider":"signalfx"}`,
			expect:  State{},
			errVal:  "json: unknown field \"provider\"",
		},
		{
			name:    "partial json set",
			create:  true,
			content: `{"auth_token":"aaa"}`,
			expect:  State{AuthToken: "aaa"},
			errVal:  "",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			path := path.Join(t.TempDir(), ".signalfx.conf")
			if tc.create {
				f, err := os.Create(path)
				require.NoError(t, err, "Must not error creating temp file")
				_, _ = io.WriteString(f, tc.content)
				require.NoError(t, f.Close(), "Must not error closing file")
			}

			var actual State
			if err := FileProviderLookup(path).Do(context.Background(), &actual); tc.errVal != "" {
				require.EqualError(t, err, tc.errVal, "Must match the expected error value")
			} else {
				require.NoError(t, err, "Must not ")
			}
			assert.Equal(t, tc.expect, actual, "Must match the expected configuration")
		})
	}
}

func TestUserFileProviderLookup(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		current func(t *testing.T) func() (*user.User, error)
		create  bool
		content string
		expect  State
		errVal  string
	}{
		{
			name: "loaded user",
			current: func(t *testing.T) func() (*user.User, error) {
				dir := t.TempDir()
				return func() (*user.User, error) {
					return &user.User{HomeDir: dir}, nil
				}
			},
			create:  true,
			content: `{}`,
			expect:  State{},
			errVal:  "",
		},
		{
			name: "failed to load user",
			current: func(t *testing.T) func() (*user.User, error) {
				return func() (*user.User, error) {
					return nil, errors.New("user not found")
				}
			},
			create:  false,
			content: ``,
			expect:  State{},
			errVal:  "user not found",
		},
		// Since this extends FileProviderLookup, there is no need to repeat tests
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Since concurrent test can depend on the same directory,
			// parallel is not enabled.
			current := tc.current(t)
			if tc.create {
				usr, _ := current()

				f, err := os.Create(path.Join(usr.HomeDir, ".signalfx.conf"))
				require.NoError(t, err, "Must not error creating file")

				_, _ = io.WriteString(f, tc.content)
				require.NoError(t, f.Close(), "Must not error closing file")
			}

			var actual State
			if err := UserFileProviderLookup(current).Do(context.Background(), &actual); tc.errVal != "" {
				require.EqualError(t, err, tc.errVal, "Must match the expected error")
			} else {
				require.NoError(t, err, "Must not error when reading user provided file")
			}
			assert.Equal(t, tc.expect, actual, "Must match the expected state")
		})
	}
}

func TestNetrcFileProvider(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		newFile func(t *testing.T) (path string)
		expect  State
		errVal  string
	}{
		{
			name: "no path set",
			newFile: func(t *testing.T) (path string) {
				return ""
			},
			expect: State{},
			errVal: "",
		},
		{
			name: "path defined, no file exists",
			newFile: func(t *testing.T) (path string) {
				return filepath.Join(t.TempDir(), NetrcFile)
			},
			expect: State{},
			errVal: "",
		},
		{
			name: "file exist, no value",
			newFile: func(t *testing.T) (path string) {
				p := filepath.Join(t.TempDir(), NetrcFile)
				f, err := os.Create(p)
				require.NoError(t, err, "Must not error when creating temp dir")
				require.NoError(t, f.Close(), "Must not error when closing file")
				return p
			},
			expect: State{},
			errVal: "",
		},
		{
			name: "path exist, defined as directory",
			newFile: func(t *testing.T) (path string) {
				return t.TempDir()
			},
			expect: State{},
			errVal: "",
		},
		{
			name: "file exist, auth defined",
			newFile: func(t *testing.T) (path string) {
				p := filepath.Join(t.TempDir(), NetrcFile)
				f, err := os.Create(p)
				require.NoError(t, err, "Must not error when creating file")
				_, _ = fmt.Fprintln(f, "machine api.signalfx.com login user1 password secret")
				require.NoError(t, f.Close(), "Must not error closing file")
				return p
			},
			expect: State{
				AuthToken: "secret",
			},
			errVal: "",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var state State
			if err := NetrcFileProviderLookup(tc.newFile(t)).Do(context.Background(), &state); tc.errVal != "" {
				require.EqualError(t, err, tc.errVal, "Must match the expected error")
			} else {
				require.NoError(t, err, "Must not return an error")
			}
			assert.Equal(t, tc.expect, state, "Must match the expected state")
		})
	}
}
