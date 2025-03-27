// Copyright Splunk, Inc.
// SPDX-License-Identifier: MPL-2.0

package convert

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestSchemaListAll(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		schema *schema.Set
		expect []any
	}{
		{
			name:   "nil set",
			schema: nil,
			expect: nil,
		},
		{
			name:   "no values set",
			schema: schema.NewSet(schema.HashInt, nil),
			expect: nil,
		},
		{
			name:   "int set",
			schema: schema.NewSet(schema.HashInt, []any{1, 2}),
			expect: []any{1, 2},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expect, SchemaListAll(tc.schema, ToAny))
		})
	}
}
