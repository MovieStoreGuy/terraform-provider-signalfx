package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderValidation(t *testing.T) {
	t.Parallel()

	require.NoError(t, New().InternalValidate(), "Must not report an error when validating")
}

func TestNewStateConfigureFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		values map[string]any
		diags  diag.Diagnostics
	}{
		{
			name:   "using default values",
			values: map[string]any{},
			diags: diag.Diagnostics{
				{Severity: diag.Error, Summary: "auth token not set"},
			},
		},
		{
			name: "provided all values",
			values: map[string]any{
				"auth_token":        "aaa",
				"api_url":           "https://api.signalfx.com",
				"custom_domain_url": "https://custom.signalfx.com",
			},
			diags: nil,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			diag := New().Configure(
				context.Background(),
				terraform.NewResourceConfigRaw(tc.values),
			)
			assert.Equal(t, tc.diags, diag, "Must match the expected value")
		})
	}
}
