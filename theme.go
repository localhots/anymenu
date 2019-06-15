package menu

import (
	"github.com/veandco/go-sdl2/sdl"
)

type colorScheme struct {
	Background     sdl.Color
	ItemBackground sdl.Color
	ItemText       sdl.Color
	TextBusy       sdl.Color
}

var theme = themeDracula

var themeGrayScale = colorScheme{
	Background:     sdl.Color{R: 20, G: 20, B: 20, A: 255},
	ItemBackground: sdl.Color{R: 60, G: 60, B: 60, A: 255},
	ItemText:       sdl.Color{R: 240, G: 240, B: 240, A: 255},
}

var themeBlueOnBlack = colorScheme{
	Background:     sdl.Color{R: 20, G: 20, B: 20, A: 255},
	ItemBackground: sdl.Color{R: 0, G: 100, B: 200, A: 255},
	ItemText:       sdl.Color{R: 255, G: 255, B: 255, A: 255},
}

var themeDracula = colorScheme{
	Background:     rgb(40, 42, 54),
	ItemBackground: rgb(68, 71, 90),
	ItemText:       rgb(248, 248, 242),
	TextBusy:       rgb(255, 184, 108),
}

func rgb(r, g, b uint8) sdl.Color {
	return sdl.Color{R: r, G: g, B: b, A: 255}
}
