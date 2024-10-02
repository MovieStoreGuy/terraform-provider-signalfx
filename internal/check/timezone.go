package check

import (
	"fmt"
	"time"
	_ "time/tzdata" // Importing time zone database to ensure there is failover option

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tfext "github.com/splunk-terraform/terraform-provider-signalfx/internal/tfextension"
)

func TimeZoneLocation() schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		tz, ok := i.(string)
		if !ok {
<<<<<<< HEAD
			return tfext.AsErrorDiagnostics(
				fmt.Errorf("expected %v as string", i),
				p,
			)
		}
		_, err := time.LoadLocation(tz)
		return tfext.AsErrorDiagnostics(err, p)

=======
			return diag.Diagnostics{
				{
					Severity:      diag.Error,
					Summary:       fmt.Sprintf("expected %v as string", i),
					AttributePath: p,
				},
			}
		}
		_, err := time.LoadLocation(tz)
		if err == nil {
			return nil
		}
		return diag.Diagnostics{
			{Severity: diag.Error, Summary: err.Error(), AttributePath: p},
		}
>>>>>>> 485475d (Adding attribute path as part of the diagnostics output)
	}
}
