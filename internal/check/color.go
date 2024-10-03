package check

import (
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tfext "github.com/splunk-terraform/terraform-provider-signalfx/internal/tfextension"
	"github.com/splunk-terraform/terraform-provider-signalfx/internal/visual"
)

func ColorValue(p visual.ColorPallete) schema.SchemaValidateDiagFunc {
	return func(i any, attr cty.Path) diag.Diagnostics {
		s, ok := i.(string)
		if !ok {
			return tfext.AsErrorDiagnostics(
				fmt.Errorf("expected %v to be type string", i), attr,
			)
		}

		if _, ok := p.GetColorCode(s); ok {
			return nil
		}
		return tfext.AsErrorDiagnostics(
			fmt.Errorf("value %q must be one of %v", s, p.Colors()),
			attr,
		)
	}
}
