// Copyright Splunk, Inc.
// SPDX-License-Identifier: MPL-2.0

package convert

import (
	"testing"

	"github.com/signalfx/signalfx-go/detector"
	"github.com/splunk-terraform/terraform-provider-signalfx/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestToDetectorPublishLabelOptions(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		in     any
		expect *detector.PublishLabelOptions
	}{
		{
			name:   "nil value",
			in:     nil,
			expect: nil,
		},
		{
			name: "default map",
			in: map[string]any{
				"label":        "pub-01",
				"display_name": "",
				"value_unit":   "s",
				"value_prefix": "",
				"value_suffix": "",
				"color":        "",
			},
			expect: &detector.PublishLabelOptions{
				Label:     "pub-01",
				ValueUnit: "s",
			},
		},
		{
			name: "all values defined",
			in: map[string]any{
				"label":        "pub-01",
				"display_name": "P-01",
				"value_unit":   "s",
				"value_prefix": "^",
				"value_suffix": "$",
				"color":        "red",
			},
			expect: &detector.PublishLabelOptions{
				Label:        "pub-01",
				DisplayName:  "P-01",
				ValueUnit:    "s",
				ValuePrefix:  "^",
				ValueSuffix:  "$",
				PaletteIndex: common.AsPointer[int32](0),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				tc.expect,
				ToDetectorPublishLabelOptions(tc.in),
			)
		})
	}
}
