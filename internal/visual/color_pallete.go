package visual

type ColorPallete map[string]int

func NewPalleteColors() ColorPallete {
	return ColorPallete{
		"gray":       0,
		"blue":       1,
		"azure":      2,
		"navy":       3,
		"brown":      4,
		"orange":     5,
		"yellow":     6,
		"magenta":    7,
		"purple":     8,
		"pink":       9,
		"violet":     10,
		"lilac":      11,
		"iris":       12,
		"emerald":    13,
		"green":      14,
		"aquamarine": 15,
	}
}

func NewFullPalleteColors() ColorPallete {
	return ColorPallete{
		"gray":        0,
		"blue":        1,
		"azure":       2,
		"navy":        3,
		"brown":       4,
		"orange":      5,
		"yellow":      6,
		"magenta":     7,
		"purple":      8,
		"pink":        9,
		"violet":      10,
		"lilac":       11,
		"iris":        12,
		"emerald":     13,
		"green":       14,
		"aquamarine":  15,
		"red":         16,
		"gold":        17,
		"greenyellow": 18,
		"chartreuse":  19,
		"jade":        20,
	}
}

func (p ColorPallete) GetColorCode(val string) (int, bool) {
	v, ok := p[val]
	return v, ok
}

func (p ColorPallete) Colors() []string {
	// This will ensure that colors are returned in color code order
	colors := make([]string, len(p))
	for c, i := range p {
		colors[i] = c
	}
	return colors
}
