/*
 * Organizations API
 *
 * API for adding, retrieving, updating, and deleting members in a SignalFx  organization. <br> ## Overview Data coming in to SignalFx is associated with a single entity called an **organization** To access this data within an organization, users must be  members of the organization. Most SignalFx users belong to a single  organization, which they think of as their \"SignalFx account\". <br> All SignalFx users have user IDs that SignalFx assigns, which identifies the  them across organizations at SignalFx. The SignalFx users for an  organization also have a **member ID* that's specific to the organization. <br> To join an organization, SignalFx users need an invitation from an member that has administrative access to the organization. After users  join an organization, they can do the following for the organization:<br>   * Submit datapoints to their organization   * Use the SignalFx web UI to look at their organization's data   * Make requests with the SignalFx API to work with organization data  ## Authentication To authenticate with SignalFx, the following operations require a session token associated with a SignalFx user that has administrative privileges:<br>   * Retrieve metadata for the organization - **GET** `/organization`   * Create, update, or delete custom categories for an organization - **PATCH** `/organization/custom-categories`   * Invite a user to the organization - **POST** `/organization/member`   * Retrieve information about the organization's users - **GET** `/organization/member`   * Invite one or more users to the organization - **POST** `/organization/members`   * Update the administrative privileges of a user - **PUT** `/organization/member/{id}`   * Delete a user - **DELETE** `/organization/member/{id}`  The following operations can authenticate with either an org token or a session token:   * Get information for an organization member - **GET** `/organization/member/{id}`   * Get all custom categories for an organization - **GET** `/organization/custom-categories`
 *
 * API version: 3.2.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package organization

// Properties containing data that SignalFx maintains for each member of an organization.
type Member struct {
	// SignalFx-assigned ID of the user that created this member
	Creator string `json:"creator,omitempty"`
	// SignalFx-assigned ID of the user that last updated this member
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	// The member creation date and time, in Unix time UTC-relative. The  system sets this value, and you can't change it.
	Created int64 `json:"created,omitempty"`
	// The date and time that the member was last updated, in Unix time  UTC-relative. The system sets this value, and you can't change it.
	LastUpdated int64 `json:"lastUpdated,omitempty"`
	// SignalFx-assigned ID of the member
	Id string `json:"id,omitempty"`
	// SignalFx-assigned user ID for this member
	UserId string `json:"userId,omitempty"`
	// SignalFx-assigned organization ID for the organization that the member belongs to
	OrganizationId string `json:"organizationId,omitempty"`
	// Email address for the SignalFx user associated with this member record
	Email string `json:"email,omitempty"`
	// Full name of the user associated with this member record
	FullName string `json:"fullName,omitempty"`
	// Phone number of the user associated with this member record
	Phone string `json:"phone,omitempty"`
	// Job title of the user associated with this member record
	Title string `json:"title,omitempty"`
	// Administrator status flag. If `true`, this member has administrative authorization for the organization.
	Admin bool `json:"admin,omitempty"`
}