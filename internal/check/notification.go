package check

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/splunk-terraform/terraform-provider-signalfx/internal/definition/notification/notificationhelper"
)

func Notification() schema.SchemaValidateDiagFunc {
	return func(i interface{}, p cty.Path) diag.Diagnostics {
		s, ok := i.(string)
		if !ok {
			return diag.Errorf("expected %v to be of type string", i)
		}
		// Using the helper library to avoid repeating code
		_, err := notificationhelper.NewNotificationFromString(s)
		return diag.FromErr(err)
	}
}
