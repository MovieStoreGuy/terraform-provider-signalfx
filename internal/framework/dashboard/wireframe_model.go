// Copyright Splunk, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwdashboard

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	fwtypes "github.com/splunk-terraform/terraform-provider-signalfx/internal/framework/types"
)

type ResourceDashboardWireframeModel struct {
	ID                types.String                            `tfsdk:"id"`
	GroupID           types.String                            `tfsdk:"group_id"`
	URL               types.String                            `tfsdk:"url"`
	Name              types.String                            `tfsdk:"name"`
	Description       types.String                            `tfsdk:"description"`
	MaxDelayOverride  types.Int32                             `tfsdk:"max_delay_override"`
	Density           types.String                            `tfsdk:"density"`
	Tags              types.List                              `tfsdk:"tags"`
	TimeRange         fwtypes.TimeRange                       `tfsdk:"time_range"`
	AuthorizedWriters *ResourceDashboardAuthorizedWritersType `tfsdk:"authorized_writers"`
	Permissions       *ResourceDashboardAccessControlListType `tfsdk:"permissions"`
	Charts            []*ResourceDashboardChartType           `tfsdk:"charts"`
}

type ResourceDashboardAuthorizedWritersType struct {
	Teams types.List `tfsdk:"teams"`
	Users types.List `tfsdk:"users"`
}

type ResourceDashboardAccessControlListType struct {
	Parent types.String                                        `tfsdk:"parent"`
	ACL    []*ResourceDashboardAccessControlListPermissionType `tfsdk:"acl"`
}

type ResourceDashboardAccessControlListPermissionType struct {
	Actions       types.List   `tfsdk:"actions"`
	PrincipalID   types.String `tfsdk:"principal_id"`
	PrincipalType types.String `tfsdk:"principal_type"`
}
