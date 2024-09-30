package check

import (
	"time"
	_ "time/tzdata" // Importing time zone database to ensure there is failover option

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TimeZoneLocation() schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		tz, ok := i.(string)
		if !ok {
			return diag.Errorf("expected %v as string", i)
		}
		_, err := time.LoadLocation(tz)
		return diag.FromErr(err)
	}
}
