package visual

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPalletes(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		p    ColorPallete
		keys []string
	}{
		{
			name: "default color pallete",
			p:    NewPalleteColors(),
			keys: []string{
				"gray",
				"blue",
				"azure",
				"navy",
				"brown",
				"orange",
				"yellow",
				"magenta",
				"purple",
				"pink",
				"violet",
				"lilac",
				"iris",
				"emerald",
				"green",
				"aquamarine",
			},
		},
		{
			name: "full color pallete",
			p:    NewFullPalleteColors(),
			keys: []string{
				"gray",
				"blue",
				"azure",
				"navy",
				"brown",
				"orange",
				"yellow",
				"magenta",
				"purple",
				"pink",
				"violet",
				"lilac",
				"iris",
				"emerald",
				"green",
				"aquamarine",
				"red",
				"gold",
				"greenyellow",
				"chartreuse",
				"jade",
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.NotEmpty(t, tc.p, "Must have values defined")
			assert.Equal(t, tc.keys, tc.p.Colors(), "Must match the expected keys definition")
		})
	}
}

func TestPalleteGetColorCode(t *testing.T) {
	t.Parallel()

	for code, color := range NewFullPalleteColors().Colors() {
		code, color := code, color
		t.Run(color, func(t *testing.T) {
			t.Parallel()

			actual, exist := NewFullPalleteColors().GetColorCode(color)
			assert.Equal(t, code, actual, "Must match the expected code")
			assert.True(t, exist, "Must match the expected contains value")
		})
	}
}
