package check

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/splunk-terraform/terraform-provider-signalfx/internal/visual"
	"github.com/stretchr/testify/assert"
)

func TestColorValue(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		val    any
		expect diag.Diagnostics
	}{
		{
			name: "no value provided",
			val:  nil,
			expect: diag.Diagnostics{
				{Severity: diag.Error, Summary: "expected <nil> to be type string"},
			},
		},
		{
			name:   "valid color",
			val:    "gray",
			expect: nil,
		},
		{
			name: "invalid color",
			val:  "spacegrey",
			expect: diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "value \"spacegrey\" must be one of [gray blue azure navy brown orange yellow magenta purple pink violet lilac iris emerald green aquamarine red gold greenyellow chartreuse jade]",
				},
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual := ColorValue(visual.NewFullPalleteColors())(tc.val, cty.Path{})
			assert.Equal(t, tc.expect, actual, "Must match the expected value")
		})
	}
}
