// Copyright Splunk, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwdashboard

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/signalfx/signalfx-go/chart"
	fwtypes "github.com/splunk-terraform/terraform-provider-signalfx/internal/framework/types"
)

type ResourceDashboardChartType struct {
	ID             types.String                              `tfsdk:"id"`
	Name           types.String                              `tfsdk:"name"`
	Tags           types.List                                `tfsdk:"tags"`
	Description    types.String                              `tfsdk:"description"`
	Position       *ResourceDashboardChartPositionType       `tfsdk:"position"`
	HeatMap        *ResourceDashboardChartHeatMapType        `tfsdk:"heatmap"`
	List           *ResourceDashboardChartListType           `tfsdk:"list"`
	Table          *ResourceDashboardChartTableType          `tfsdk:"table"`
	SingleValue    *ResourceDashboardChartSingleValueType    `tfsdk:"single_value"`
	Text           *ResourceDashboardChartTextType           `tfsdk:"text"`
	TimeSeries     *ResourceDashboardChartTimeSeriesType     `tfsdk:"time_series"`
	SLO            *ResourceDashboardChartSLOType            `tfsdk:"slo"`
	Program        *ResourceDashboardChartProgramOptionType  `tfsdk:"program"`
	DataOptions    *ResourceDashboardChartDataOptionType     `tfsdk:"data_options"`
	PublishOptions *ResourceDashboardChartPublishOptionsType `tfsdk:"publish_options"`
}

type ResourceDashboardChartPositionType struct {
	Width  types.Int32 `tfsdk:"width"`
	Height types.Int32 `tfsdk:"height"`
	Row    types.Int32 `tfsdk:"row"`
	Column types.Int32 `tfsdk:"column"`
}

func (rdcp *ResourceDashboardChartPositionType) CheckCollision(other *ResourceDashboardChartPositionType) bool {
	if rdcp == nil || other == nil {
		return false
	}

	return !(rdcp.Row.ValueInt32()+rdcp.Width.ValueInt32() <= other.Row.ValueInt32() ||
		other.Row.ValueInt32()+other.Width.ValueInt32() <= rdcp.Row.ValueInt32() ||
		rdcp.Column.ValueInt32()+rdcp.Height.ValueInt32() <= other.Column.ValueInt32() ||
		other.Column.ValueInt32()+other.Height.ValueInt32() <= rdcp.Column.ValueInt32())
}

type ResourceDashboardChartProgramOptionType struct {
	Text            types.String      `tfsdk:"text"`
	MinResolution   fwtypes.TimeRange `tfsdk:"min_resolution"`
	MaxDelay        fwtypes.TimeRange `tfsdk:"max_delay"`
	DisableSampling types.Bool        `tfsdk:"disable_sampling"`
	Timezone        types.String      `tfsdk:"timezone"`
}

type ResourceDashboardChartDataOptionType struct {
	NoData            *ResourceDashboardChartNoDataOptionType `tfsdk:"no_data"`
	RefreshInterval   fwtypes.TimeRange                       `tfsdk:"refresh_interval"`
	HideMissingValues types.Bool                              `tfsdk:"hide_missing_values"`
	MaxPrecision      types.Int32                             `tfsdk:"max_precision"`
	UnitPrefix        types.String                            `tfsdk:"unit_prefix"`
}

type ResourceDashboardChartNoDataOptionType struct {
	Message  types.String `tfsdk:"message"`
	LinkText types.String `tfsdk:"link_text"`
	LinkURL  types.String `tfsdk:"link_url"`
}

type ResourceDashboardChartHeatMapType struct {
	ColorBy types.String `tfsdk:"color_by"`
	GroupBy types.List   `tfsdk:"group_by"`
}

type ResourceDashboardChartListType struct {
	HideMissingValues types.Bool `tfsdk:"hide_missing_values"`
}

type ResourceDashboardChartTableType struct {
	GroupBy           types.List `tfsdk:"group_by"`
	HideMissingValues types.Bool `tfsdk:"hide_missing_values"`
}

type ResourceDashboardChartSingleValueType struct {
	ShowSparkline types.Bool `tfsdk:"show_sparkline"`
}

type ResourceDashboardChartTextType struct {
	Markdown types.String `tfsdk:"content"`
}

type ResourceDashboardChartTimeSeriesType struct {
	Options *ResourceDashboardChartTimeSeriesOptionsType `tfsdk:"options"`
}

type ResourceDashboardChartTimeSeriesOptionsType struct {
	ColorBy         types.String `tfsdk:"color_by"`
	IncludeZero     types.Bool   `tfsdk:"include_zero"`
	ShowDataMarkers types.Bool   `tfsdk:"show_data_markers"`
	ShowLegend      types.Bool   `tfsdk:"show_legend"`
	ShowEventLints  types.Bool   `tfsdk:"show_event_lines"`
	SortBy          types.String `tfsdk:"sort_by"`
	Stacked         types.Bool   `tfsdk:"stacked"`
}

type ResourceDashboardChartSLOType struct {
	ID types.String `tfsdk:"slo_id"`
}

type ResourceDashboardChartPublishOptionsType struct {
	Label       types.String `tfsdk:"label"`
	DisplayName types.String `tfsdk:"display_name"`
	ValuePrefix types.String `tfsdk:"value_prefix"`
	ValueSuffix types.String `tfsdk:"value_suffix"`
	ValueUnit   types.String `tfsdk:"value_unit"`
	PlotType    types.String `tfsdk:"plot_type"`
}

func (ch *ResourceDashboardChartType) NewHeatmapCreateUpdateChartRequest(ctx context.Context) (*chart.CreateUpdateChartRequest, diag.Diagnostics) {
	if ch.HeatMap == nil || ch.DataOptions == nil || ch.Program == nil {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Missing Heatmap definition",
				"Cannot create heatmap chart request without heatmap, data options, and program defined",
			),
		}
	}

	heatmap := &chart.CreateUpdateChartRequest{
		Name:        ch.Name.ValueString(),
		Description: ch.Description.ValueString(),
		ProgramText: ch.Program.Text.ValueString(),
		Options: &chart.Options{
			Type:       "Heatmap",
			UnitPrefix: ch.DataOptions.UnitPrefix.ValueString(),
			ProgramOptions: &chart.GeneralOptions{
				DisableSampling: ch.Program.DisableSampling.ValueBool(),
				Timezone:        ch.Program.Timezone.ValueString(),
			},
			HideMissingValues: ch.DataOptions.HideMissingValues.ValueBool(),
			ColorBy:           ch.HeatMap.ColorBy.ValueString(),
		},
	}
	diags := ch.Tags.ElementsAs(ctx, &heatmap.Tags, false)
	diags.Append(ch.HeatMap.GroupBy.ElementsAs(ctx, &heatmap.Options.GroupBy, false)...)

	if !ch.DataOptions.RefreshInterval.IsNull() || !ch.DataOptions.RefreshInterval.IsUnknown() {
		dur := int32(ch.DataOptions.RefreshInterval.ValueDuration().Milliseconds())
		heatmap.Options.RefreshInterval = &dur
	}

	if !ch.Program.MinResolution.IsNull() || !ch.Program.MinResolution.IsUnknown() {
		dur := int32(ch.Program.MinResolution.ValueDuration().Milliseconds())
		heatmap.Options.ProgramOptions.MinimumResolution = &dur
	}

	if !ch.Program.MaxDelay.IsNull() || !ch.Program.MaxDelay.IsUnknown() {
		dur := int32(ch.Program.MaxDelay.ValueDuration().Milliseconds())
		heatmap.Options.ProgramOptions.MaxDelay = &dur
	}

	// TODO: Missing color options

	return heatmap, diags
}

func (ch *ResourceDashboardChartType) NewListCreateUpdateChartRequest(ctx context.Context) (*chart.CreateUpdateChartRequest, diag.Diagnostics) {
	if ch.List == nil || ch.DataOptions == nil || ch.Program == nil {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Missing List definition",
				"Cannot create list chart request without list, data options, and program defined",
			),
		}
	}

	list := &chart.CreateUpdateChartRequest{
		Name:        ch.Name.ValueString(),
		Description: ch.Description.ValueString(),
		ProgramText: ch.Program.Text.ValueString(),
		Options: &chart.Options{
			Type:       "List",
			UnitPrefix: ch.DataOptions.UnitPrefix.ValueString(),
		},
	}

	return list, nil
}

func (ch *ResourceDashboardChartType) NewTableCreateUpdateChartRequest(ctx context.Context) (*chart.CreateUpdateChartRequest, diag.Diagnostics) {
	if ch.Table == nil || ch.DataOptions == nil || ch.Program == nil {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Missing Table definition",
				"Cannot create table chart request without table, data options, and program defined",
			),
		}
	}

	table := &chart.CreateUpdateChartRequest{
		Name:        ch.Name.ValueString(),
		Description: ch.Description.ValueString(),
		ProgramText: ch.Program.Text.ValueString(),
		Options: &chart.Options{
			Type:       "Table",
			UnitPrefix: ch.DataOptions.UnitPrefix.ValueString(),
			ProgramOptions: &chart.GeneralOptions{
				DisableSampling: ch.Program.DisableSampling.ValueBool(),
				Timezone:        ch.Program.Timezone.ValueString(),
			},
			HideMissingValues: ch.DataOptions.HideMissingValues.ValueBool(),
		},
	}
	diags := ch.Tags.ElementsAs(ctx, &table.Tags, false)
	diags.Append(ch.Table.GroupBy.ElementsAs(ctx, &table.Options.GroupBy, false)...)

	if !ch.DataOptions.RefreshInterval.IsNull() || !ch.DataOptions.RefreshInterval.IsUnknown() {
		dur := int32(ch.DataOptions.RefreshInterval.ValueDuration().Milliseconds())
		table.Options.RefreshInterval = &dur
	}

	if !ch.Program.MinResolution.IsNull() || !ch.Program.MinResolution.IsUnknown() {
		dur := int32(ch.Program.MinResolution.ValueDuration().Milliseconds())
		table.Options.ProgramOptions.MinimumResolution = &dur
	}

	if !ch.Program.MaxDelay.IsNull() || !ch.Program.MaxDelay.IsUnknown() {
		dur := int32(ch.Program.MaxDelay.ValueDuration().Milliseconds())
		table.Options.ProgramOptions.MaxDelay = &dur
	}

	return table, diags
}

func (ch *ResourceDashboardChartType) NewSingleValueCreateUpdateChartRequest(ctx context.Context) (*chart.CreateUpdateChartRequest, diag.Diagnostics) {
	if ch.SingleValue == nil || ch.DataOptions == nil || ch.Program == nil {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Missing Single Value definition",
				"Cannot create single value chart request without single value, data options, and program defined",
			),
		}
	}

	singleValue := &chart.CreateUpdateChartRequest{
		Name:        ch.Name.ValueString(),
		Description: ch.Description.ValueString(),
		ProgramText: ch.Program.Text.ValueString(),
		Options: &chart.Options{
			Type:       "SingleValue",
			UnitPrefix: ch.DataOptions.UnitPrefix.ValueString(),
			ProgramOptions: &chart.GeneralOptions{
				DisableSampling: ch.Program.DisableSampling.ValueBool(),
				Timezone:        ch.Program.Timezone.ValueString(),
			},
			HideMissingValues: ch.DataOptions.HideMissingValues.ValueBool(),
		},
	}
	diags := ch.Tags.ElementsAs(ctx, &singleValue.Tags, false)

	if !ch.DataOptions.RefreshInterval.IsNull() || !ch.DataOptions.RefreshInterval.IsUnknown() {
		dur := int32(ch.DataOptions.RefreshInterval.ValueDuration().Milliseconds())
		singleValue.Options.RefreshInterval = &dur
	}

	if !ch.Program.MinResolution.IsNull() || !ch.Program.MinResolution.IsUnknown() {
		dur := int32(ch.Program.MinResolution.ValueDuration().Milliseconds())
		singleValue.Options.ProgramOptions.MinimumResolution = &dur
	}

	if !ch.Program.MaxDelay.IsNull() || !ch.Program.MaxDelay.IsUnknown() {
		dur := int32(ch.Program.MaxDelay.ValueDuration().Milliseconds())
		singleValue.Options.ProgramOptions.MaxDelay = &dur
	}

	if ch.SingleValue.ShowSparkline.ValueBool() {
		singleValue.Options.ShowSparkLine = ch.SingleValue.ShowSparkline.ValueBool()
	}

	return singleValue, diags
}

func (ch *ResourceDashboardChartType) NewTextCreateUpdateChartRequest(ctx context.Context) (*chart.CreateUpdateChartRequest, diag.Diagnostics) {
	if ch.Text == nil || ch.DataOptions == nil {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Missing Text definition",
				"Cannot create text chart request without text and data options defined",
			),
		}
	}

	text := &chart.CreateUpdateChartRequest{
		Name:        ch.Name.ValueString(),
		Description: ch.Description.ValueString(),
		Options: &chart.Options{
			Type:     "Text",
			Markdown: ch.Text.Markdown.ValueString(),
		},
	}
	diags := ch.Tags.ElementsAs(ctx, &text.Tags, false)

	return text, diags
}

func (ch *ResourceDashboardChartType) NewTimeSeriesCreateUpdateChartRequest(ctx context.Context) (*chart.CreateUpdateChartRequest, diag.Diagnostics) {
	if ch.TimeSeries == nil || ch.DataOptions == nil || ch.Program == nil || ch.TimeSeries.Options == nil {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Missing Time Series definition",
				"Cannot create time series chart request without time series, time series options, data options, and program defined",
			),
		}
	}

	timeSeries := &chart.CreateUpdateChartRequest{
		Name:        ch.Name.ValueString(),
		Description: ch.Description.ValueString(),
		ProgramText: ch.Program.Text.ValueString(),
		Options: &chart.Options{
			Type:       "TimeSeries",
			UnitPrefix: ch.DataOptions.UnitPrefix.ValueString(),
			ProgramOptions: &chart.GeneralOptions{
				DisableSampling: ch.Program.DisableSampling.ValueBool(),
				Timezone:        ch.Program.Timezone.ValueString(),
			},
			HideMissingValues: ch.DataOptions.HideMissingValues.ValueBool(),
			ColorBy:           ch.TimeSeries.Options.ColorBy.ValueString(),
			IncludeZero:       ch.TimeSeries.Options.IncludeZero.ValueBool(),
			ShowEventLines:    ch.TimeSeries.Options.ShowEventLints.ValueBool(),
			SortBy:            ch.TimeSeries.Options.SortBy.ValueString(),
			Stacked:           ch.TimeSeries.Options.Stacked.ValueBool(),
		},
	}
	diags := ch.Tags.ElementsAs(ctx, &timeSeries.Tags, false)

	if !ch.DataOptions.RefreshInterval.IsNull() || !ch.DataOptions.RefreshInterval.IsUnknown() {
		dur := int32(ch.DataOptions.RefreshInterval.ValueDuration().Milliseconds())
		timeSeries.Options.RefreshInterval = &dur
	}

	if !ch.DataOptions.MaxPrecision.IsNull() || !ch.DataOptions.MaxPrecision.IsUnknown() {
		prec := int32(ch.DataOptions.MaxPrecision.ValueInt32())
		timeSeries.Options.MaximumPrecision = &prec
	}

	if !ch.Program.MinResolution.IsNull() || !ch.Program.MinResolution.IsUnknown() {
		dur := int32(ch.Program.MinResolution.ValueDuration().Milliseconds())
		timeSeries.Options.ProgramOptions.MinimumResolution = &dur
	}

	if !ch.Program.MaxDelay.IsNull() || !ch.Program.MaxDelay.IsUnknown() {
		dur := int32(ch.Program.MaxDelay.ValueDuration().Milliseconds())
		timeSeries.Options.ProgramOptions.MaxDelay = &dur
	}

	return timeSeries, diags
}
