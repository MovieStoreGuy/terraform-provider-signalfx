package providerhelper

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestMustNewResourceMap(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		entries []ResourceDefinitionFunc
		panics  bool
	}{
		{
			name:    "no entries",
			entries: []ResourceDefinitionFunc{},
			panics:  false,
		},
		{
			name: "duplicate entries",
			entries: []ResourceDefinitionFunc{
				func() (string, *schema.Resource) {
					return "nop", nil
				},
				func() (string, *schema.Resource) {
					return "nop", nil
				},
			},
			panics: true,
		},
		{
			name: "valid entries",
			entries: []ResourceDefinitionFunc{
				func() (string, *schema.Resource) {
					return "a", nil
				},
				func() (string, *schema.Resource) {
					return "b", nil
				},
			},
			panics: false,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.panics {
				assert.Panics(
					t,
					func() {
						MustNewResourceMap(tc.entries...)
					},
					"Must panic when performing insert",
				)
			} else {
				assert.NotPanics(
					t,
					func() {
						MustNewResourceMap(tc.entries...)
					},
					"Must not panic performing insert",
				)
			}
		})
	}
}
