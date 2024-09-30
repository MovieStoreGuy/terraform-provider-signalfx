package check

import (
	"slices"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ValueUnit() schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		s, ok := i.(string)
		if !ok {
			return diag.Errorf("expected %v to be type string", i)
		}

		units := []string{
			"Bit",
			"Kilobit",
			"Megabit",
			"Gigabit",
			"Terabit",
			"Petabit",
			"Exabit",
			"Zettabit",
			"Yottabit",
			"Byte",
			"Kibibyte",
			"Mebibyte",
			"Gibibyte",
			"Tebibyte",
			"Pebibyte",
			"Exbibyte",
			"Zebibyte",
			"Yobibyte",
			"Nanosecond",
			"Microsecond",
			"Millisecond",
			"Second",
			"Minute",
			"Hour",
			"Day",
			"Week",
		}

		if slices.Contains(units, s) {
			return nil
		}

		return diag.Errorf("expected %q to be one of %v", s, units)
	}
}
