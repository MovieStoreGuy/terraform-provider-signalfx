package providerhelper

import (
	"context"
	"errors"
	"net/url"
	"path"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/signalfx/signalfx-go"
	"go.uber.org/multierr"
)

var (
	ErrStateNotImplemented = errors.New("expected to implement type State")
)

type State struct {
	AuthToken    string           `json:"auth_token"`
	APIURL       string           `json:"api_url"`
	CustomAppURL string           `json:"custom_app_url"`
	Client       *signalfx.Client `json:"-"`
}

// LoadClient returns the configured [signalfx.Client] ready to use.
//
// Note that it is a shared instance so high amounts of parallelism could cause issues.
func LoadClient(ctx context.Context, meta any) (*signalfx.Client, error) {
	if c, ok := meta.(*State); ok {
		return c.Client, nil
	}
	tflog.Error(ctx, "Failed to load state from meta value", map[string]any{
		"meta": meta,
	})
	return nil, ErrStateNotImplemented
}

// LoadApplicationURL will generate the FQDN using the set CustomAppURL from the provider.
func LoadApplicationURL(ctx context.Context, meta any, fragments ...string) (string, error) {
	s, ok := meta.(*State)
	if !ok {
		return "", ErrStateNotImplemented
	}
	u, err := url.ParseRequestURI(s.CustomAppURL)
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}
	u.Fragment = path.Join(fragments...)
	return u.String(), nil
}

func (s *State) Validate() (errs error) {
	if s.AuthToken == "" {
		errs = multierr.Append(errs, errors.New("auth token not set"))
	}
	if s.APIURL == "" {
		errs = multierr.Append(errs, errors.New("api url is not set"))
	}
	return errs
}
