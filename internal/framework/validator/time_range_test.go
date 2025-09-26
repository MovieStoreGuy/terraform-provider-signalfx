// Copyright Splunk, Inc.
// SPDX-License-Identifier: MPL-2.0

package fwvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeStringRequest(val string) validator.StringRequest {
	return validator.StringRequest{
		ConfigValue: types.StringValue(val),
	}
}

func TestTimeRangeGreaterThan(t *testing.T) {
	v := TimeRangeGreaterThan("1h")
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("2h"), resp)
	assert.False(t, resp.Diagnostics.HasError())

	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("1h"), resp)
	assert.False(t, resp.Diagnostics.HasError())

	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("30m"), resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Time Range Too Small")
}

func TestTimeRangeLessThan(t *testing.T) {
	v := TimeRangeLessThan("2h")
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("1h"), resp)
	assert.False(t, resp.Diagnostics.HasError())

	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("2h"), resp)
	assert.False(t, resp.Diagnostics.HasError())

	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("3h"), resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Time Range Too Large")
}

func TestTimeRangeBetween(t *testing.T) {
	v := TimeRangeBetween("1h", "3h")
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("2h"), resp)
	assert.False(t, resp.Diagnostics.HasError())

	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("1h"), resp)
	assert.False(t, resp.Diagnostics.HasError())

	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("3h"), resp)
	assert.False(t, resp.Diagnostics.HasError())

	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("30m"), resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Time Range Out of Bounds")

	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("4h"), resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Time Range Out of Bounds")
}

func TestTimeRangeValidator_NullOrUnknown(t *testing.T) {
	v := TimeRangeGreaterThan("1h")
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), validator.StringRequest{
		ConfigValue: types.StringNull(),
	}, resp)
	assert.False(t, resp.Diagnostics.HasError())

	resp = &validator.StringResponse{}
	v.ValidateString(context.Background(), validator.StringRequest{
		ConfigValue: types.StringUnknown(),
	}, resp)
	assert.False(t, resp.Diagnostics.HasError())
}

func TestTimeRangeValidator_InvalidFormat(t *testing.T) {
	v := TimeRangeGreaterThan("1h")
	resp := &validator.StringResponse{}
	v.ValidateString(context.Background(), makeStringRequest("notaduration"), resp)
	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Invalid Time Range")
}

func TestMarkdownDescription(t *testing.T) {
	assert.Equal(t,
		"The value must be between 1h and 2h.",
		timeRangeBounds{min: "1h", max: "2h"}.MarkdownDescription(context.Background()),
	)
	assert.Equal(t,
		"The value must be greater than 1h.",
		timeRangeBounds{min: "1h"}.MarkdownDescription(context.Background()),
	)
	assert.Equal(t,
		"The value must be less than 2h.",
		timeRangeBounds{max: "2h"}.MarkdownDescription(context.Background()),
	)
	assert.Equal(t,
		"",
		timeRangeBounds{}.MarkdownDescription(context.Background()),
	)
}

func TestDescriptionDelegatesToMarkdownDescription(t *testing.T) {
	bounds := timeRangeBounds{min: "1h", max: "2h"}
	require.Equal(t, bounds.MarkdownDescription(context.Background()), bounds.Description(context.Background()))
}
