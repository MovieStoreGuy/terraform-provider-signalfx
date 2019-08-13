/*
 * Integrations API
 *
 * APIs for creating, retrieving, updating, and deleting SignalFx integrations to the systems you use.<br> An integration provides SignalFx with information from the external system that you're connecting to. You'll need to retrieve this information from the external system before you use the API. Each external system is different, so to see a summary of its requirements and procedures, view its request body description. # Authentication To create, update, delete, or validate an integration, you need to authenticate your request using a session token associated with a SignalFx administrator. To **retrieve** an integration, your session token doesn't need to be associated with an administrator. You can also retrieve integrations using an org token.<br> In the web UI, session tokens are known as <strong>user access</strong> tokens, and org tokens are known as <strong>access tokens</strong>. <br> To learn more about authentication tokens, see the topic [Authentication Tokens](https://developers.signalfx.com/administration/access_tokens_overview.html) in the Developers Guide. # Supported service types SignalFx offers integrations for the following:<br>   * Data collection from other monitoring systems such as AWS CloudWatch   * Authentication using your existing Single Sign-On (**SSO**) system   * Sending alerts using your preferred messaging, chat, or incident management service <br> To use one of these integrations, you first register it with SignalFx. After that, you configure the integration to communicate between the system you're using and SignalFx. ## Data collection SignalFx integrations APIs support data collection for the following services:<br>   * Amazon Web Services (**AWS**)   * Google Cloud Platform (**GCP**)   * Microsoft Azure   * NewRelic  ## Authentication using SSO SignalFx integration APIs support SAML-based SSO integrations for the following services:<br>   * Microsoft Active Directory Federation Services (**ADFS**)   * Bitium   * Okta   * OneLogin   * PingOne  ## Alerts using message, chat, or incident management services SignalFx integration APIs support alert notifications using the following services: <br>   * BigPanda   * Office 365   * Opsgenie   * PagerDuty   * ServiceNow   * Slack   * VictorOps   * Webhook   * xMatters<br>  **NOTE:** You can't create Office 365 integrations using the API, and your ability to update them in a **PUT** request is limited, but you can retrieve their data or delete them. To create an Office 365 integration, use the the web UI. <br> # Viewing request body documentation The *request* body format for the following operations depends on the type of integration you use:<br>   * POST `/integration`   * PUT `/integration/{id}`<br>  The *response* body format for the following operations also depends on the type of integration you use:<br>   * GET `/integration`   * GET `/integration/{id}`  <br>  To see the request or response body format for an integration: <br>   1. Find the endpoint and method.   2. For a request body, find the section *REQUEST BODY SCHEMA*. For a     response body, find the section *RESPONSE SCHEMA*.   3. Scroll down to the `type` property.   4. At the end of the description for `type`, find the dropdown box that      contains the integration type. By default, it's set to *AWSCloudWatch*.   5. To see a complete list of integrations, click the down arrow. A list      with a vertical scroll bar appears.   6. Select the integration type from the list. The request body properties      for this integration type now appear.
 *
 * API version: 3.3.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package integration

// Specifies the data collection integration between AWS CloudWatch and SignalFx, in the form of a JSON object.
type AwsCloudWatchIntegration struct {
	// The creation date and time for the integration object, in Unix time UTC-relative. The system sets this value, and you can't modify it.
	Created int64 `json:"created,omitempty"`
	// SignalFx-assigned user ID of the user that created the integration object. If the system created the object, the value is \"AAAAAAAAAA\". The system sets this value, and you can't modify it.
	Creator string `json:"creator,omitempty"`
	// Flag that indicates the state of the integration object. If  `true`, the integration is enabled. If `false`, the integration is disabled, and you must enable it by setting \"enabled\" to `true` in a **PUT** request that updates the object. <br> **NOTE:** SignalFx always sets the flag to `true` when you call  **POST** `/integration` to create an integration.
	Enabled bool `json:"enabled"`
	// SignalFx-assigned ID of an integration you create in the web UI or API. Use this property to retrieve an integration using the **GET**, **PUT**, or **DELETE** `/integration/{id}` endpoints or the **GET** `/integration/validate{id}/` endpoint, as described in this topic.
	Id string `json:"id,omitempty"`
	// The last time the integration was updated, in Unix time UTC-relative. This value is \"read-only\".
	LastUpdated int64 `json:"lastUpdated,omitempty"`
	// SignalFx-assigned ID of the last user who updated the integration. If the last update was by the system, the value is \"AAAAAAAAAA\". This value is \"read-only\".
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	// A human-readable label for the integration. This property helps you identify a specific integration when you're using multiple integrations for the same service.
	Name       string        `json:"name,omitempty"`
	Type       Type          `json:"type"`
	AuthMethod AwsAuthMethod `json:"authMethod,omitempty"`
	// Comma-separated list of custom AWS CloudWatch namespaces to monitor. Custom namespaces contain custom metrics that you define in AWS; SignalFx imports the metrics so you can monitor them. See the \"Publishing Metrics\" AWS documentation for more information. If you specify the `customNamespaceSyncRules` property, SignalFx ignores the value of `customCloudWatchNamespaces`.
	CustomCloudWatchNamespaces string `json:"customCloudWatchNamespaces,omitempty"`
	// Array of objects that specify custom AWS namespaces and filters. Each element controls the data collected by SignalFx for the specified namespace.  If you specify this property, SignalFx ignores values in the  \"customCloudWatchNamespaces\" property.  To learn more about this property, see the topic  [Integrating With AWS](https://developers.signalfx.com/integrating/aws_integration_overview.html).
	CustomNamespaceSyncRules []*AwsCustomNameSpaceSyncRule `json:"customNamespaceSyncRules,omitempty"`
	// Flag that controls how SignalFx imports usage metrics from AWS to use with AWS Optimizer. If `true`, SignalFx imports the metrics.
	EnableAwsUsage bool `json:"enableAwsUsage,omitempty"`
	// If you specify `\"authMethod\": \"ExternalId\"` in your request to create an AWS integration object, the response object contains a value for \"externalId\". Use this value and the ARN value you get from AWS to update the integration object. SignalFx can then connect to AWS using the integration object.<br> **NOTE:** SignalFx sets this value, and you can't change it.
	ExternalId string `json:"externalId,omitempty"`
	// Flag that controls how SignalFx imports Cloud Watch metrics. If `true`, SignalFx imports Cloud Watch metrics from AWS.
	ImportCloudWatch bool `json:"importCloudWatch,omitempty"`
	// If you specify `\"authMethod\": \"SecurityToken\"` in your request to create an AWS integration object, use this property to specify the key.
	Key string `json:"key,omitempty"`
	// Array of namespace sync rules. Each element in the array is an object that contains an AWS namespace name and a filter that controls the data that SignalFx collects for the namespace. If you specify this property, SignalFx ignores the values in the AWS CloudWatch Integration Model \"services\" property. If you don't specify either property, SignalFx syncs all data in all AWS namespaces. To learn more, see the topic [Integrating With AWS](https://developers.signalfx.com/integrating/aws_integration_overview.html).
	NamespaceSyncRules []*AwsNameSpaceSyncRule `json:"namespaceSyncRules,omitempty"`
	PollRate           *PollRate               `json:"pollRate,omitempty"`
	// Array of AWS regions that SignalFx should monitor. The API supports the following AWS regions: <br><br> **Regular AWS regions**    * ap-northeast-1   * ap-northeast-2   * ap-south-1   * ap-southeast-1   * ap-southeast-2   * ca-central-1   * eu-central-1   * eu-north-1   * eu-west-1   * eu-west-2   * eu-west-3   * sa-east-1   * us-east-1   * us-east-2   * us-west-1   * us-west-2  **GovCloud AWS regions**    * us-gov-east-1   * us-gov-west-1  **China AWS regions**    * cn-north-1   * cn-northwest-1  If you don't specify the \"regions\" property, or if you specify `\"regions\": []`, the API adds all of the regular AWS regions to your integration. **Note:** You can't mix regions from different sets. For example, you can't specify `\"regions\": [ \"eu-west-1\", \"cn-north-1\"]`.
	Regions []string `json:"regions,omitempty"`
	// Role ARN that you add to an existing AWS integration object.<br> When you create an AWS integration object and specify \"ExternalId\" a he authentication method, SignalFx responds with an external ID. You provide this ID to AWS, which responds with a role ARN.<br> To finish the connection between SignalFx and AWS, you update the AWS integration object using a PUT request. In this request, you specify the \"roleArn\" property using the value you obtained from AWS. <br> **NOTE:** To ensure security, SignalFx doesn't return this property in response objects.
	RoleArn string `json:"roleArn,omitempty"`
	// Array of AWS services that you want SignalFx to monitor. Each element is a string designating an AWS service
	Services []AwsService `json:"services,omitempty"`
	// If you specify `\"authMethod\": \"SecurityToken\"` in your request to create an AWS integration object, use this property to specify the token.
	Token string `json:"token,omitempty"`
	// Flag that controls how SignalFx checks for large amounts of data for this AWS integration. If `true`, SignalFx checks to see if the integration is returning a large amount of data.
	EnableCheckLargeVolume bool `json:"enableCheckLargeVolume,omitempty"`
	// If `true`, this property indicates that SignalFx is receiving a large volume of data and tags from AWS.
	IsLargeVolume bool `json:"isLargeVolume,omitempty"`
}