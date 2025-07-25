---
page_title: "Splunk Observability Cloud: signalfx_azure_integration"
description: |-
  Allows Terraform to create and manage Azure Integrations for Splunk Observability Cloud
---
# Resource: signalfx_azure_integration

Splunk Observability Cloud Azure integrations. For help with this integration see [Monitoring Microsoft Azure](https://docs.splunk.com/observability/en/gdi/get-data-in/connect/azure/azure.html).

~> **NOTE** When managing integrations, use a session token of an administrator to authenticate the Splunk Observability Cloud provider. See [Operations that require a session token for an administrator](https://dev.splunk.com/observability/docs/administration/authtokens#Operations-that-require-a-session-token-for-an-administrator). Otherwise you'll receive a 4xx error.

## Example

```terraform
resource "signalfx_azure_integration" "azure_myteam" {
  name    = "Azure Foo"
  enabled = true

  environment = "azure"

  poll_rate = 300

  secret_key = "XXX"

  app_id = "YYY"

  tenant_id = "ZZZ"

  services = ["microsoft.sql/servers/elasticpools"]

  subscriptions = ["sub-guid-here"]

  # Optional
  additional_services = ["some/service", "another/service"]

  # Optional
  custom_namespaces_per_service {
    service = "Microsoft.Compute/virtualMachines"
    namespaces = [ "monitoringAgent", "customNamespace" ]
  }

  # Optional
  resource_filter_rules {
    filter_source = "filter('azure_tag_service', 'payment') and (filter('azure_tag_env', 'prod-us') or filter('azure_tag_env', 'prod-eu'))"
  }
  resource_filter_rules {
    filter_source = "filter('azure_tag_service', 'notification') and (filter('azure_tag_env', 'prod-us') or filter('azure_tag_env', 'prod-eu'))"
  }
}
```

## Arguments

* `app_id` - (Required) Azure application ID for the Splunk Observability Cloud app. To learn how to get this ID, see the topic [Connect to Microsoft Azure](https://docs.splunk.com/observability/en/gdi/get-data-in/connect/azure/azure.html) in the product documentation.
* `enabled` - (Required) Whether the integration is enabled.
* `custom_namespaces_per_service` - (Optional) Allows for more fine-grained control of syncing of custom namespaces, should the boolean convenience parameter `sync_guest_os_namespaces` be not enough. The customer may specify a map of services to custom namespaces. If they do so, for each service which is a key in this map, we will attempt to sync metrics from namespaces in the value list in addition to the default namespaces.
  * `namespaces` - (Required) The additional namespaces.
  * `service` - (Required) The name of the service.
* `environment` (Optional) What type of Azure integration this is. The allowed values are `\"azure_us_government\"` and `\"azure\"`. Defaults to `\"azure\"`.
* `name` - (Required) Name of the integration.
* `named_token` - (Optional) Name of the org token to be used for data ingestion. If not specified then default access token is used.
* `poll_rate` - (Optional) Azure poll rate (in seconds). Value between `60` and `600`. Default: `300`.
* `resource_filter_rules` - (Optional) List of rules for filtering Azure resources by their tags.
  * `filter_source` - (Required) Expression that selects the data that Splunk Observability Cloud should sync for the resource associated with this sync rule. The expression uses the syntax defined for the SignalFlow `filter()` function. The source of each filter rule must be in the form filter('key', 'value'). You can join multiple filter statements using the and and or operators. Referenced keys are limited to tags and must start with the azure_tag_ prefix.
* `secret_key` - (Required) Azure secret key that associates the Splunk Observability Cloud app in Azure with the Azure tenant ID. To learn how to get this ID, see the topic [Connect to Microsoft Azure](https://docs.splunk.com/observability/en/gdi/get-data-in/connect/azure/azure.html) in the product documentation.
* `services` - (Required) List of Microsoft Azure service names for the Azure services you want Splunk Observability Cloud to monitor. Can be an empty list to import data for all supported services. See [Microsoft Azure services](https://docs.splunk.com/Observability/gdi/get-data-in/integrations.html#azure-integrations) for a list of valid values.
* `subscriptions` - (Required) List of Azure subscriptions that Splunk Observability Cloud should monitor.
* `sync_guest_os_namespaces` - (Optional) If enabled, Splunk Observability Cloud will try to sync additional namespaces for VMs (including VMs in scale sets): telegraf/mem, telegraf/cpu, azure.vm.windows.guest (these are namespaces recommended by Azure when enabling their Diagnostic Extension). If there are no metrics there, no new datapoints will be ingested. Defaults to false.
* `import_azure_monitor` - (Optional) If enabled, Splunk Observability Cloud will sync also Azure Monitor data. If disabled, Splunk Observability Cloud will import only metadata. Defaults to true.
* `tenant_id` (Required) Azure ID of the Azure tenant. To learn how to get this ID, see the topic [Connect to Microsoft Azure](https://docs.splunk.com/observability/en/gdi/get-data-in/connect/azure/azure.html) in the product documentation.
* `use_batch_api` - (Optional) If enabled, Splunk Observability Cloud will collect datapoints using Azure Metrics Batch API. Consider this option if you are synchronizing high loads of data and you want to avoid throttling issues. Contrary to the default Metrics List API, Metrics Batch API is paid. Refer to [Azure documentation](https://azure.microsoft.com/en-us/pricing/details/api-management/) for pricing info.

## Attributes

In a addition to all arguments above, the following attributes are exported:

* `id` - The ID of the integration.
