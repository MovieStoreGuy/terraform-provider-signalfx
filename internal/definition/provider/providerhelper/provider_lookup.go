package providerhelper

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/bgentry/go-netrc/netrc"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/go-homedir"
)

// ProviderLookupFunc allows for unified type to preconfigure the
// the state and allow it to be easer to ensure providers are working as expected.
type ProviderLookupFunc func(ctx context.Context, s *State) error

func (fn ProviderLookupFunc) Do(ctx context.Context, s *State) error {
	if fn == nil {
		return errors.New("not implemented")
	}
	return fn(ctx, s)
}

// NewDefaultProviderLookups returns the list of
func NewDefaultProviderLookups() []ProviderLookupFunc {
	return []ProviderLookupFunc{
		FileProviderLookup("/etc/signalfx.conf"),
		UserFileProviderLookup(user.Current),
		NetrcFileProviderLookup(os.Getenv("NETRC")),
	}
}

func FileProviderLookup(path string) ProviderLookupFunc {
	return func(ctx context.Context, s *State) error {
		tflog.Debug(ctx, "Reading provider file", map[string]any{
			"path": path,
		})

		f, err := os.Open(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				err = errors.New("file not found")
			}
			return err
		}

		tflog.Debug(ctx, "Reading file content")

		dec := json.NewDecoder(f)
		dec.DisallowUnknownFields()

		if err = dec.Decode(s); err != nil {
			if errors.Is(err, io.EOF) {
				err = errors.New("no file content")
			}
			return errors.Join(err, f.Close())
		}

		return f.Close()
	}
}

func UserFileProviderLookup(current func() (*user.User, error)) ProviderLookupFunc {
	return func(ctx context.Context, s *State) error {
		u, err := current()
		if err != nil {
			return err
		}
		return FileProviderLookup(path.Join(u.HomeDir, ".signalfx.conf")).Do(ctx, s)
	}
}

func NetrcFileProviderLookup(path string) ProviderLookupFunc {
	return func(ctx context.Context, s *State) error {
		if path == "" {
			var err error
			path, err = homedir.Expand(filepath.Join("~/", NetrcFile))
			if err != nil {
				return err
			}
		}
		st, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				err = nil
			}
			return err
		}

		if st.IsDir() {
			return nil
		}

		tflog.Debug(ctx, "Reading netrc file", map[string]any{
			"path": path,
		})

		m, err := netrc.FindMachine(path, "api.signalfx.com")
		if err != nil {
			return err
		}
		if m != nil && !m.IsDefault() {
			s.AuthToken = m.Password
		}
		return nil
	}
}
