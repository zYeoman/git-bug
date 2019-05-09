package mardown

import (
	"github.com/logrusorgru/aurora"
)

var shades = []aurora.Color{
	aurora.GreenFg | aurora.BoldFm,
	aurora.GreenFg | aurora.BoldFm,
	aurora.BrightFg | aurora.GreenFg,
	aurora.GreenFg,
}

// Return the color function corresponding to the level.
// Beware, level start counting from 1.
func shade(level int) aurora.Color {
	if level < 1 {
		level = 1
	}
	if level > len(shades) {
		level = len(shades)
	}
	return shades[level-1]
}
