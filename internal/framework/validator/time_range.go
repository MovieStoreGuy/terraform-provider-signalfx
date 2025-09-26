// Copyright Splunk, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	fwtypes "github.com/splunk-terraform/terraform-provider-signalfx/internal/framework/types"
)

type timeRangeBounds struct {
	min string
	max string
}

func TimeRangeGreaterThan(val string) validator.String {
	return timeRangeBounds{min: val}
}

func TimeRangeLessThan(val string) validator.String {
	return timeRangeBounds{max: val}
}

func TimeRangeBetween(min, max string) validator.String {
	return timeRangeBounds{min: min, max: max}
}

func (v timeRangeBounds) Description(_ context.Context) string {
	return v.MarkdownDescription(context.Background())
}

func (v timeRangeBounds) MarkdownDescription(_ context.Context) string {
	switch {
	case v.min != "" && v.max != "":
		return "The value must be between " + v.min + " and " + v.max + "."
	case v.min != "":
		return "The value must be greater than " + v.min + "."
	case v.max != "":
		return "The value must be less than " + v.max + "."
	default:
		return ""
	}
}

func (v timeRangeBounds) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	tr := &fwtypes.TimeRange{
		StringValue: req.ConfigValue,
	}

	dur, err := tr.ParseDuration()
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Time Range",
			"Could not parse the time range value: "+err.Error(),
		)
		return
	}

	var (
		min = fwtypes.TimeRange{StringValue: types.StringValue(v.min)}.ValueDuration()
		max = fwtypes.TimeRange{StringValue: types.StringValue(v.max)}.ValueDuration()
	)

	switch {
	case v.min != "" && v.max != "":
		if min <= dur && dur <= max {
			return
		}
		resp.Diagnostics.AddError(
			"Time Range Out of Bounds",
			"The time range value must be between "+v.min+" and "+v.max+".",
		)
	case v.min != "":
		if min <= dur {
			return
		}
		resp.Diagnostics.AddError(
			"Time Range Too Small",
			"The time range value must be greater than "+v.min+".",
		)
	case v.max != "":
		if dur <= max {
			return
		}
		resp.Diagnostics.AddError(
			"Time Range Too Large",
			"The time range value must be less than "+v.max+".",
		)
	}
}
