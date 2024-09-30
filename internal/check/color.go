package check

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/splunk-terraform/terraform-provider-signalfx/internal/visual"
)

func ColorValue(p visual.ColorPallete) schema.SchemaValidateDiagFunc {
	return func(i any, _ cty.Path) diag.Diagnostics {
		s, ok := i.(string)
		if !ok {
			return diag.Errorf("expected %v to be type string", i)
		}

		if _, ok := p.GetColorCode(s); ok {
			return nil
		}
		return diag.Errorf("value %q must be one of %v", s, p.Colors())
	}
}
