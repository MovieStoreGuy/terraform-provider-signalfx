package check

import (
	"slices"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SeverityLevel() schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Errorf("expected %v to be of type string", i)
		}

		labels := []string{
			"Critical", "Major", "Minor", "Warning", "Info",
		}

		if slices.Contains(labels, value) {
			return nil
		}

		return diag.Errorf("value %q is not allowed; must be one of: %v", value, labels)
	}
}
