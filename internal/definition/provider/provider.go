package provider

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/signalfx/signalfx-go"

	"github.com/splunk-terraform/terraform-provider-signalfx/internal/definition/provider/providerhelper"
	"github.com/splunk-terraform/terraform-provider-signalfx/internal/definition/team"
	"github.com/splunk-terraform/terraform-provider-signalfx/version"
)

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SFX_AUTH_TOKEN", ""),
				Description: "Splunk Observability Cloud auth token",
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SFX_API_URL", "https://api.signalfx.com"),
				Description: "API URL for your Splunk Observability Cloud org, may include a realm",
			},
			"custom_app_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SFX_CUSTOM_APP_URL", "https://app.signalfx.com"),
				Description: "Application URL for your Splunk Observability Cloud org, often customized for organizations using SSO",
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     120,
				Description: "Timeout duration for a single HTTP call in seconds. Defaults to 120",
			},
			"retry_max_attempts": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     4,
				Description: "Max retries for a single HTTP call. Defaults to 4",
			},
			"retry_wait_min_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Minimum retry wait for a single HTTP call in seconds. Defaults to 1",
			},
			"retry_wait_max_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Maximum retry wait for a single HTTP call in seconds. Defaults to 30",
			},
		},
		ResourcesMap: providerhelper.MustNewResourceMap(
			team.NewResource(),
		),
		DataSourcesMap:       providerhelper.MustNewResourceMap(),
		ConfigureContextFunc: newStateConfigureFunc(),
	}
}

func newStateConfigureFunc() schema.ConfigureContextFunc {
	return func(ctx context.Context, rd *schema.ResourceData) (any, diag.Diagnostics) {
		var s providerhelper.State
		for _, fn := range providerhelper.NewDefaultProviderLookups() {
			if err := fn.Do(ctx, &s); err != nil && !errors.Is(err, os.ErrNotExist) {
				tflog.Debug(ctx, "Failed to update provider state from lookup provider", map[string]any{
					"error": err,
				})
			}
		}
		if v, ok := rd.GetOk("auth_token"); ok {
			s.AuthToken = v.(string)
		}
		if v, ok := rd.GetOk("api_url"); ok {
			s.APIURL = v.(string)
		}
		if v, ok := rd.GetOk("custom_app_url"); ok {
			s.CustomAppURL = v.(string)
		}
		if err := s.Validate(); err != nil {
			return nil, diag.FromErr(err)
		}

		tp := logging.NewLoggingHTTPTransport(&http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			DialContext:         (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		})

		retryer := retryablehttp.NewClient()
		retryer.RetryMax = rd.Get("retry_max_attempts").(int)
		retryer.RetryWaitMin = time.Duration(rd.Get("retry_wait_min_seconds").(int)) * time.Second
		retryer.RetryWaitMax = time.Duration(rd.Get("retry_wait_max_seconds").(int)) * time.Second
		retryer.HTTPClient.Timeout = time.Duration(rd.Get("timeout_seconds").(int)) * time.Second
		retryer.HTTPClient.Transport = tp

		tflog.Debug(ctx, "Configured retry http client", map[string]any{
			"max-attempts": retryer.RetryMax,
			"wait-min":     retryer.RetryWaitMin.String(),
			"wait-max":     retryer.RetryWaitMax.String(),
			"timeout":      retryer.HTTPClient.Timeout.String(),
		})

		var err error
		s.Client, err = signalfx.NewClient(
			s.AuthToken,
			signalfx.APIUrl(s.APIURL),
			signalfx.HTTPClient(retryer.StandardClient()),
			signalfx.UserAgent(fmt.Sprintf("terraform-provider-signalfx/%s", version.ProviderVersion)),
		)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return &s, nil
	}
}
