package mardown

import (
	"github.com/MichaelMure/git-bug/util/colors"
)

var shades = []func(a ...interface{}) string{
	colors.GreenBold,
	colors.GreenBold,
	colors.HiGreen,
	colors.Green,
}

// Return the color function corresponding to the level.
// Beware, level start counting from 1.
func shade(level int) func(a ...interface{}) string {
	if level < 1 {
		level = 1
	}
	if level > len(shades) {
		level = len(shades)
	}
	return shades[level-1]
}
