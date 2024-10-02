package providerhelper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadClient(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name  string
		state any
		err   error
	}{
		{
			name:  "state not set",
			state: nil,
			err:   ErrStateNotImplemented,
		},
		{
			name:  "state defined",
			state: &State{},
			err:   nil,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := LoadClient(context.Background(), tc.state)
			require.ErrorIs(t, err, tc.err, "Must match the expected error value")
		})
	}
}

func TestLoadApplicationURL(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		state     any
		fragments []string
		url       string
		errVal    string
	}{
		{
			name:      "no state set",
			state:     nil,
			fragments: []string{},
			url:       "",
			errVal:    "expected to implement type State",
		},
		{
			name: "custom domain set",
			state: &State{
				CustomAppURL: "http://custom.signalfx.com",
			},
			fragments: []string{},
			url:       "http://custom.signalfx.com/",
			errVal:    "",
		},
		{
			name: "custom domain with fragments",
			state: &State{
				CustomAppURL: "http://custom.signalfx.com",
			},
			fragments: []string{
				"detector",
				"aaaa",
				"edit",
			},
			url:    "http://custom.signalfx.com/#detector/aaaa/edit",
			errVal: "",
		},
		{
			name: "invalid domain set",
			state: &State{
				CustomAppURL: "domain",
			},
			fragments: []string{},
			url:       "",
			errVal:    "parse \"domain\": invalid URI for request",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			u, err := LoadApplicationURL(context.Background(), tc.state, tc.fragments...)
			require.Equal(t, tc.url, u, "Must match the expected url")
			if tc.errVal != "" {
				require.EqualError(t, err, tc.errVal, "Must match expected error message")
			} else {
				require.NoError(t, err, "Must not error when loading url")
			}
		})
	}
}

func TestStateValidation(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		state  State
		errVal string
	}{
		{
			name:   "state not set",
			state:  State{},
			errVal: "auth token not set; api url is not set",
		},
		{
			name: "state valid",
			state: State{
				AuthToken: "aaa",
				APIURL:    "http://api.signalfx.com",
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if err := tc.state.Validate(); tc.errVal != "" {
				require.EqualError(t, err, tc.errVal, "Must match the expected error")
			} else {
				require.NoError(t, err, "Must not error when validation")
			}
		})
	}
}
