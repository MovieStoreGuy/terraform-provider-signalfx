// Copyright Splunk, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwdashboard

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/signalfx/signalfx-go/chart"
	"github.com/signalfx/signalfx-go/dashboard"

	fwembed "github.com/splunk-terraform/terraform-provider-signalfx/internal/framework/embed"
	"github.com/splunk-terraform/terraform-provider-signalfx/internal/framework/fwerr"
	fwshared "github.com/splunk-terraform/terraform-provider-signalfx/internal/framework/shared"
	fwtypes "github.com/splunk-terraform/terraform-provider-signalfx/internal/framework/types"
	fwvalidator "github.com/splunk-terraform/terraform-provider-signalfx/internal/framework/validator"
	pmeta "github.com/splunk-terraform/terraform-provider-signalfx/internal/providermeta"
)

type ResourceDashboardWireframe struct {
	fwembed.DatasourceData
	fwembed.ResourceIDImporter
}

var (
	_ resource.Resource                   = (*ResourceDashboardWireframe)(nil)
	_ resource.ResourceWithImportState    = (*ResourceDashboardWireframe)(nil)
	_ resource.ResourceWithValidateConfig = (*ResourceDashboardWireframe)(nil)
)

func NewResourceDashboardWireframe() resource.Resource {
	return &ResourceDashboardWireframe{}
}

func (rdw *ResourceDashboardWireframe) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dashboard_wireframe"
}

func (rdw *ResourceDashboardWireframe) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a resource that allows creating a dashboard with the required charts",
		Attributes: map[string]schema.Attribute{
			"id": fwshared.ResourceIDAttribute(),
			"url": schema.StringAttribute{
				Description: "The URL of the dashboard.",
				Computed:    true,
			},
			"group_id": schema.StringAttribute{
				Description: "The ID of the group to contain the dashboard.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the dashboard.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the dashboard.",
				Optional:    true,
			},
			"max_delay_override": schema.StringAttribute{
				Description: "The maximum delay override in seconds.",
				Optional:    true,
				Computed:    true,
			},
			"density": schema.StringAttribute{
				Description: "The density of the dashboard.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("DEFAULT"),
				Validators: []validator.String{
					stringvalidator.OneOf("DEFAULT", "LOW", "HIGH", "HIGHEST"),
				},
			},
			"tags": schema.ListAttribute{
				Description: "A list of tags to associate with the dashboard.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"time_range": schema.StringAttribute{
				Description: "This sets how far back the dashboard will load data from for each of the charts.",
				Optional:    true,
				Computed:    true,
				CustomType:  fwtypes.TimeRangeType{},
				Default:     stringdefault.StaticString("-3h"),
			},
			"authorized_writers": schema.SingleNestedAttribute{
				Description: "Defines who can write to the dashboard.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"teams": schema.ListAttribute{
						Description: "A list of team IDs that can write to the dashboard.",
						ElementType: types.StringType,
						Optional:    true,
					},
					"users": schema.ListAttribute{
						Description: "A list of user IDs that can write to the dashboard.",
						ElementType: types.StringType,
						Optional:    true,
					},
				},
			},
			"permissions": schema.SingleNestedAttribute{
				Description: "Defines the access control list for the dashboard.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"parent": schema.StringAttribute{
						Optional: true,
					},
					"acl": schema.ListNestedAttribute{
						Description: "A list of access control list permissions.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"actions": schema.ListAttribute{
									Description: "A list of actions that the principal can perform. ",
									ElementType: types.StringType,
									Required:    true,
									Validators: []validator.List{
										listvalidator.ValueStringsAre(
											stringvalidator.OneOf("READ", "WRITE"),
										),
									},
								},
								"principal_id": schema.StringAttribute{
									Description: "The ID of the principal (UserID, TeamID, or OrganizationID).",
									Required:    true,
								},
								"principal_type": schema.StringAttribute{
									Description: "The type of the principal.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf("USER", "TEAM", "ORG"),
									},
								},
							},
						},
					},
				},
			},
			"chart": schema.SetNestedAttribute{
				Description: "A set of charts to include in the dashboard.",
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Set{
					setvalidator.ExactlyOneOf(
						path.MatchRelative().AtName("heatmap"),
						path.MatchRelative().AtName("list"),
						path.MatchRelative().AtName("table"),
						path.MatchRelative().AtName("single_value"),
						path.MatchRelative().AtName("text"),
						path.MatchRelative().AtName("time_series"),
						path.MatchRelative().AtName("slo"),
					),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the chart.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the chart.",
							Required:    true,
						},
						"description": schema.StringAttribute{
							Description: "The description of the chart.",
							Optional:    true,
						},
						"tags": schema.ListAttribute{
							Description: "A list of tags to associate with the chart.",
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
						},
						"position": schema.SingleNestedAttribute{
							Description: "The position of the chart on the dashboard grid.",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								"width": schema.Int32Attribute{
									Description: "The width of the chart in grid units.",
									Optional:    true,
									Computed:    true,
									Default:     int32default.StaticInt32(3),
									Validators: []validator.Int32{
										int32validator.Between(1, 12),
									},
								},
								"height": schema.Int32Attribute{
									Description: "The height of the chart in grid units.",
									Optional:    true,
									Computed:    true,
									Default:     int32default.StaticInt32(1),
									Validators: []validator.Int32{
										int32validator.Between(1, 3),
									},
								},
								"row": schema.Int32Attribute{
									Description: "The row position of the chart on the grid.",
									Required:    true,
									Validators: []validator.Int32{
										int32validator.Between(0, 99),
									},
								},
								"column": schema.Int32Attribute{
									Description: "The column position of the chart on the grid.",
									Required:    true,
									Validators: []validator.Int32{
										int32validator.Between(0, 11),
									},
								},
							},
						},
						"program_options": schema.SingleNestedAttribute{
							Description: "Program Options is used to detail what data is being queried and additional options that can improve query performance.",
							Required:    true,
							Attributes: map[string]schema.Attribute{
								"text": schema.StringAttribute{
									Required:    true,
									Description: "This is the program text used to query data stored in Splunk Observability Cloud.",
								},
								"min_resolution": schema.StringAttribute{
									Optional:   true,
									CustomType: fwtypes.TimeRangeType{},
									Validators: []validator.String{
										fwvalidator.TimeRangeGreaterThan("0s"),
									},
								},
								"max_delay": schema.StringAttribute{
									Optional:   true,
									CustomType: fwtypes.TimeRangeType{},
									Validators: []validator.String{
										fwvalidator.TimeRangeGreaterThan("0s"),
									},
								},
								"disable_sampling": schema.BoolAttribute{
									Optional: true,
									Computed: true,
									Default:  booldefault.StaticBool(false),
								},
								"timezone": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Default:  stringdefault.StaticString("UTC"),
								},
							},
						},
						"heatmap": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"color_by": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Validators: []validator.String{
										stringvalidator.OneOf("Range", "Scale"),
									},
									Default: stringdefault.StaticString("Range"),
								},
								"group_by": schema.ListAttribute{
									Optional:    true,
									ElementType: types.StringType,
								},
							},
						},
						"list": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"hide_missing_values": schema.BoolAttribute{
									Optional: true,
									Computed: true,
									Default:  booldefault.StaticBool(false),
								},
							},
						},
						"table": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"group_by": schema.ListAttribute{
									Optional:    true,
									ElementType: types.StringType,
								},
								"hide_missing_values": schema.BoolAttribute{
									Optional: true,
									Computed: true,
									Default:  booldefault.StaticBool(false),
								},
							},
						},
						"single_value": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"show_sparkline": schema.BoolAttribute{
									Optional: true,
									Computed: true,
									Default:  booldefault.StaticBool(false),
								},
							},
						},
						"text": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"markdown": schema.StringAttribute{
									Required: true,
								},
							},
						},
						"time_series": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"color_by": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Default:  stringdefault.StaticString("Dimension"),
									Validators: []validator.String{
										stringvalidator.OneOf("Metric", "Dimension"),
									},
								},
								"include_zero": schema.BoolAttribute{
									Optional: true,
									Computed: true,
									Default:  booldefault.StaticBool(false),
								},
								"show_data_markers": schema.BoolAttribute{
									Optional: true,
									Computed: true,
									Default:  booldefault.StaticBool(false),
								},
								"show_legend": schema.BoolAttribute{
									Optional: true,
									Computed: true,
									Default:  booldefault.StaticBool(true),
								},
								"show_event_lines": schema.BoolAttribute{
									Optional: true,
									Computed: true,
									Default:  booldefault.StaticBool(true),
								},
								"sort_by": schema.StringAttribute{
									Optional: true,
								},
								"stacked": schema.BoolAttribute{
									Optional: true,
									Computed: true,
									Default:  booldefault.StaticBool(false),
								},
							},
						},
						"slo": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"slo_id": schema.StringAttribute{
									Description: "The ID of the SLO to display.",
									Required:    true,
								},
							},
						},
						"publish_options": schema.ListNestedAttribute{
							Description: "Allows a user to configure how each published stream on the chart is rendered",
							Optional:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"label": schema.StringAttribute{
										Required: true,
									},
									"display_name": schema.StringAttribute{
										Optional: true,
									},
									"value_prefix": schema.StringAttribute{
										Optional: true,
									},
									"value_suffix": schema.StringAttribute{
										Optional: true,
									},
									"value_unit": schema.StringAttribute{
										Optional: true,
									},
									"plot_type": schema.StringAttribute{
										Optional: true,
									},
								},
							},
						},
						"data_options": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"no_data": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"message": schema.StringAttribute{
											Required: true,
										},
										"link_text": schema.StringAttribute{
											Optional: true,
										},
										"link_url": schema.StringAttribute{
											Optional: true,
										},
									},
								},
								"refresh_interval": schema.StringAttribute{
									Optional:   true,
									CustomType: fwtypes.TimeRangeType{},
								},
								"hide_missing_values": schema.BoolAttribute{
									Optional: true,
									Computed: true,
									Default:  booldefault.StaticBool(false),
								},
								"max_precision": schema.Int32Attribute{
									Optional: true,
									Computed: true,
									Default:  int32default.StaticInt32(2),
									Validators: []validator.Int32{
										int32validator.Between(0, 10),
									},
								},
								"time_range": schema.StringAttribute{
									Optional:   true,
									CustomType: fwtypes.TimeRangeType{},
								},
								"unit_prefix": schema.StringAttribute{
									Optional: true,
									Computed: true,
									Default:  stringdefault.StaticString("Metric"),
									Validators: []validator.String{
										stringvalidator.OneOf("Metric", "Binary"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (rdw *ResourceDashboardWireframe) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model ResourceDashboardWireframeModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	charts := make(map[string]*ResourceDashboardChartPositionType)
	for _, ch := range model.Charts {
		switch {
		case ch.HeatMap != nil:
			heatmap, diag := ch.NewHeatmapCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
				return
			}
			created, err := rdw.Details().Client.CreateChart(ctx, heatmap)
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("chart").AtSetValue(types.StringValue(ch.Name.ValueString())),
					"Error Creating Chart",
					fmt.Sprintf("An error was encountered creating the heatmap chart %q: %s", ch.Name.ValueString(), err.Error()),
				)
				continue
			}
			ch.ID = types.StringValue(created.Id)
			charts[created.Id] = ch.Position
		case ch.List != nil:
			list, diag := ch.NewListCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
				return
			}
			created, err := rdw.Details().Client.CreateChart(ctx, list)
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("chart").AtSetValue(types.StringValue(ch.Name.ValueString())),
					"Error Creating Chart",
					fmt.Sprintf("An error was encountered creating the list chart %q: %s", ch.Name.ValueString(), err.Error()),
				)
				continue
			}
			ch.ID = types.StringValue(created.Id)
			charts[created.Id] = ch.Position
		case ch.Table != nil:
			table, diags := ch.NewTableCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
			created, err := rdw.Details().Client.CreateChart(ctx, table)
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("chart").AtSetValue(types.StringValue(ch.Name.ValueString())),
					"Error Creating Chart",
					fmt.Sprintf("An error was encountered creating the table chart %q: %s", ch.Name.ValueString(), err.Error()),
				)
				continue
			}
			ch.ID = types.StringValue(created.Id)
			charts[created.Id] = ch.Position
		case ch.SingleValue != nil:
			singleValue, diags := ch.NewSingleValueCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
			created, err := rdw.Details().Client.CreateChart(ctx, singleValue)
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("chart").AtSetValue(types.StringValue(ch.Name.ValueString())),
					"Error Creating Chart",
					fmt.Sprintf("An error was encountered creating the single value chart %q: %s", ch.Name.ValueString(), err.Error()),
				)
				continue
			}
			ch.ID = types.StringValue(created.Id)
			charts[created.Id] = ch.Position
		case ch.Text != nil:
			text, diags := ch.NewTextCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
			created, err := rdw.Details().Client.CreateChart(ctx, text)
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("chart").AtSetValue(types.StringValue(ch.Name.ValueString())),
					"Error Creating Chart",
					fmt.Sprintf("An error was encountered creating the text chart %q: %s", ch.Name.ValueString(), err.Error()),
				)
				continue
			}
			ch.ID = types.StringValue(created.Id)
			charts[created.Id] = ch.Position
		case ch.TimeSeries != nil:
			timeSeries, diags := ch.NewTimeSeriesCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
			created, err := rdw.Details().Client.CreateChart(ctx, timeSeries)
			if err != nil {
				resp.Diagnostics.AddAttributeError(
					path.Root("chart").AtSetValue(types.StringValue(ch.Name.ValueString())),
					"Error Creating Chart",
					fmt.Sprintf("An error was encountered creating the time series chart %q: %s", ch.Name.ValueString(), err.Error()),
				)
				continue
			}
			ch.ID = types.StringValue(created.Id)
			charts[created.Id] = ch.Position
		case ch.SLO != nil:
			created, err := rdw.Details().Client.CreateSloChart(ctx, &chart.CreateUpdateSloChartRequest{
				SloId: ch.SLO.ID.ValueString(),
			})
			if resp.Diagnostics.Append(fwerr.ErrorHandler(ctx, resp.State, err)...); resp.Diagnostics.HasError() {
				return
			}
			ch.ID = types.StringValue(created.Id)
			charts[created.Id] = ch.Position
		}
	}

	db := &dashboard.CreateUpdateDashboardRequest{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		GroupId:     model.GroupID.ValueString(),
	}

	for id, pos := range charts {
		db.Charts = append(db.Charts, &dashboard.DashboardChart{
			ChartId: id,
			Width:   pos.Width.ValueInt32(),
			Height:  pos.Height.ValueInt32(),
			Row:     pos.Row.ValueInt32(),
			Column:  pos.Column.ValueInt32(),
		})
	}

	created, err := rdw.Details().Client.CreateDashboard(ctx, db)
	if err != nil {
		resp.Diagnostics.Append(fwerr.ErrorHandler(ctx, resp.State, err)...)
		return
	}

	model.ID = types.StringValue(created.Id)
	model.URL = types.StringValue(pmeta.LoadApplicationURL(ctx, rdw.Details(), "dashboard", created.Id))

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (rdw *ResourceDashboardWireframe) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model ResourceDashboardWireframeModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}
	dashboard, err := rdw.Details().Client.GetDashboard(ctx, model.ID.ValueString())
	if err != nil {
		resp.Diagnostics.Append(fwerr.ErrorHandler(ctx, req.State, err)...)
		return
	}

	model.Name = types.StringValue(dashboard.Name)
	model.Description = types.StringValue(dashboard.Description)
	model.GroupID = types.StringValue(dashboard.GroupId)
	if dashboard.ChartDensity != nil {
		model.Density = types.StringValue(string(*dashboard.ChartDensity))
	}

	for _, dchart := range dashboard.Charts {
		details, err := rdw.Details().Client.GetChart(ctx, dchart.ChartId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Chart",
				fmt.Sprintf("An error was encountered reading the chart %q: %s", dchart.ChartId, err.Error()),
			)
			continue
		}
		chart := &ResourceDashboardChartType{
			ID:          types.StringValue(details.Id),
			Name:        types.StringValue(details.Name),
			Description: types.StringValue(details.Description),
			Position: &ResourceDashboardChartPositionType{
				Width:  types.Int32Value(dchart.Width),
				Height: types.Int32Value(dchart.Height),
				Row:    types.Int32Value(dchart.Row),
				Column: types.Int32Value(dchart.Column),
			},
		}
		// TODO(MovieStoreGuy): Determine chart type
		model.Charts = append(model.Charts, chart)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (rdw *ResourceDashboardWireframe) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model ResourceDashboardWireframeModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}
}

func (rdw *ResourceDashboardWireframe) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model ResourceDashboardWireframeModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	for _, chart := range model.Charts {
		if err := rdw.Details().Client.DeleteChart(ctx, chart.ID.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting Chart",
				fmt.Sprintf("An error was encountered deleting the chart %q: %s", chart.Name.ValueString(), err.Error()),
			)
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if err := rdw.Details().Client.DeleteDashboard(ctx, model.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Dashboard",
			fmt.Sprintf("An error was encountered deleting the dashboard %q: %s", model.Name.ValueString(), err.Error()),
		)
	}
}

func (rdw *ResourceDashboardWireframe) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var model ResourceDashboardWireframeModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &model)...); resp.Diagnostics.HasError() {
		return
	}

	safe := make(map[string]*ResourceDashboardChartPositionType)
	for _, chart := range model.Charts {
		switch {
		case chart.HeatMap != nil:
			req, diags := chart.NewHeatmapCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
			if err := rdw.Details().Client.ValidateChart(ctx, req); err != nil {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("chart").AtSetValue(types.StringValue(chart.Name.ValueString())),
					"Invalid Heatmap Chart",
					fmt.Sprintf("The heatmap chart %q is invalid: %s", chart.Name.ValueString(), err.Error()),
				)
			}
		case chart.List != nil:
			req, diags := chart.NewListCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
			if err := rdw.Details().Client.ValidateChart(ctx, req); err != nil {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("chart").AtSetValue(types.StringValue(chart.Name.ValueString())),
					"Invalid List Chart",
					fmt.Sprintf("The list chart %q is invalid: %s", chart.Name.ValueString(), err.Error()),
				)
			}
		case chart.Table != nil:
			req, diags := chart.NewTableCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
			if err := rdw.Details().Client.ValidateChart(ctx, req); err != nil {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("chart").AtSetValue(types.StringValue(chart.Name.ValueString())),
					"Invalid Table Chart",
					fmt.Sprintf("The table chart %q is invalid: %s", chart.Name.ValueString(), err.Error()),
				)
			}
		case chart.SingleValue != nil:
			req, diags := chart.NewSingleValueCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
			if err := rdw.Details().Client.ValidateChart(ctx, req); err != nil {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("chart").AtSetValue(types.StringValue(chart.Name.ValueString())),
					"Invalid Single Value Chart",
					fmt.Sprintf("The single value chart %q is invalid: %s", chart.Name.ValueString(), err.Error()),
				)
			}
		case chart.Text != nil:
			req, diags := chart.NewTextCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
			if err := rdw.Details().Client.ValidateChart(ctx, req); err != nil {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("chart").AtSetValue(types.StringValue(chart.Name.ValueString())),
					"Invalid Text Chart",
					fmt.Sprintf("The text chart %q is invalid: %s", chart.Name.ValueString(), err.Error()),
				)
			}
		case chart.TimeSeries != nil:
			req, diags := chart.NewTimeSeriesCreateUpdateChartRequest(ctx)
			if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
				return
			}
			if err := rdw.Details().Client.ValidateChart(ctx, req); err != nil {
				resp.Diagnostics.AddAttributeWarning(
					path.Root("chart").AtSetValue(types.StringValue(chart.Name.ValueString())),
					"Invalid Time Series Chart",
					fmt.Sprintf("The time series chart %q is invalid: %s", chart.Name.ValueString(), err.Error()),
				)
			}
		}
		for named, pos := range safe {
			if chart.Position.CheckCollision(pos) {
				resp.Diagnostics.AddError(
					"Chart Position Collision",
					fmt.Sprintf("Chart position collision detected between charts: %q and %q", named, chart.Name.ValueString()),
				)
			}
		}
		safe[chart.Name.ValueString()] = chart.Position
	}
}
