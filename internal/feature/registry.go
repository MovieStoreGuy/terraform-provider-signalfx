// Copyright Splunk, Inc.
// SPDX-License-Identifier: MPL-2.0

package feature

import (
	"context"
	"fmt"
	"iter"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/mod/semver"

	tfext "github.com/splunk-terraform/terraform-provider-signalfx/internal/tfextension"
)

type Registry struct {
	features sync.Map
}

func NewRegistry() *Registry {
	return &Registry{}
}

// All returns a iterator that is ordered by semver version, and the feature name
func (r *Registry) All() iter.Seq2[string, *Preview] {
	var previews []*Preview
	r.features.Range(func(_, value any) bool {
		previews = append(previews, value.(*Preview))
		return true
	})

	slices.SortFunc(previews, func(a, b *Preview) int {
		if cmp := semver.Compare(a.Introduced(), b.Introduced()); cmp != 0 {
			return cmp
		}
		return strings.Compare(a.Name(), b.Name())
	})

	return func(yield func(string, *Preview) bool) {
		for _, p := range previews {
			if !yield(p.Name(), p) {
				return
			}
		}
	}
}

func (reg *Registry) Configure(ctx context.Context, feature string, enabled bool) error {
	p, ok := reg.Get(feature)
	if !ok {
		return fmt.Errorf("no preview with id %q found", feature)
	}

	if p.GlobalAvailable() {
		tflog.Warn(
			ctx,
			"Preview has been marked as Global Available,"+
				" it is no longer required to be configured. "+
				"If you're experiencing issues with a new feature, "+
				"please reach out to customer support so the issue can be addressed",
			tfext.NewLogFields().
				Field("feature", feature).
				Field("enabled", p.Enabled()),
		)
	}

	p.SetEnabled(enabled)

	tflog.Debug(ctx, "Configured feature preview", tfext.NewLogFields().
		Field("feature", feature).
		Field("enabled", p.Enabled()).
		Field("added_in", p.Introduced()).
		Field("description", p.Description()),
	)

	return nil
}

func (reg *Registry) Get(id string) (*Preview, bool) {
	if v, ok := reg.features.Load(id); ok {
		return v.(*Preview), true
	}
	return nil, false
}

func (reg *Registry) Register(feature string, opts ...PreviewOption) (*Preview, error) {
	matched, err := regexp.MatchString(`^[a-z0-9\._\-]+$`, feature)
	if err != nil {
		// Error here should technically never happen.
		return nil, err
	}

	if !matched {
		return nil, fmt.Errorf("feature %q does not match expected naming format", feature)
	}

	g, err := NewPreview(feature, opts...)
	if err != nil {
		return nil, err
	}

	if _, ok := reg.Get(feature); ok {
		return nil, fmt.Errorf("feature %q already exists", feature)
	}

	reg.features.Store(feature, g)

	return g, nil
}

func (reg *Registry) MustRegister(feature string, opts ...PreviewOption) *Preview {
	g, err := reg.Register(feature, opts...)
	if err != nil {
		panic(err)
	}
	return g
}
