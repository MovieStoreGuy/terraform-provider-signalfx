package check

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/splunk-terraform/terraform-provider-signalfx/internal/visual"
)

func ColorValue(p visual.ColorPallete) schema.SchemaValidateDiagFunc {
	return func(i any, attr cty.Path) diag.Diagnostics {
		s, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity:      diag.Error,
					Summary:       fmt.Sprintf("expected %v to be type string", i),
					AttributePath: attr,
				},
			}
		}

		if _, ok := p.GetColorCode(s); ok {
			return nil
		}
		return diag.Diagnostics{
			{
				Severity:      diag.Error,
				Summary:       fmt.Sprintf("value %q must be one of %v", s, p.Colors()),
				AttributePath: attr,
			},
		}
	}
}
