// Copyright Splunk, Inc.
// SPDX-License-Identifier: MPL-2.0

package convert

import (
	"github.com/signalfx/signalfx-go/detector"
	"github.com/splunk-terraform/terraform-provider-signalfx/internal/common"
	"github.com/splunk-terraform/terraform-provider-signalfx/internal/visual"
)

func ToDetectorPublishLabelOptions(in any) *detector.PublishLabelOptions {
	viz, ok := in.(map[string]any)
	if !ok {
		return nil
	}

	palette := visual.NewColorPalette()
	opt := &detector.PublishLabelOptions{
		Label:       viz["label"].(string),
		DisplayName: viz["display_name"].(string),
		ValueUnit:   viz["value_unit"].(string),
		ValuePrefix: viz["value_prefix"].(string),
		ValueSuffix: viz["value_suffix"].(string),
	}

	if idx, ok := palette.ColorIndex(viz["color"].(string)); ok {
		opt.PaletteIndex = common.AsPointer(idx)
	}

	return opt
}
