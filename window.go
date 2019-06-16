package menu

import (
	"github.com/juju/errors"
	"github.com/localhots/themenu/fonts"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	paddingX = 8
	paddingY = 8
)

var (
	windowWidth  = 600
	windowHeight = 800
	borderWidth  = 0

	fontName = "InconsolataGo Regular"
	fontSize = 24
	font     *ttf.Font
)

type window struct {
	window             *sdl.Window
	allItems, curItems []menuItem
}

func (w window) render(r *sdl.Renderer) error {
	// Resize window
	windowHeight = borderWidth*2 + // Border
		len(w.curItems)*itemHeight() + // Backgrounds
		(len(w.curItems)+1)*paddingY // Paddings
	w.window.SetSize(int32(windowWidth), int32(windowHeight))

	// Reset background
	if err := drawWindowBackground(r); err != nil {
		return errors.Annotate(err, "draw window background")
	}
	if err := drawWindowBorder(r, int32(borderWidth)); err != nil {
		return errors.Annotate(err, "draw window border")
	}

	// Draw items
	for i, itm := range w.curItems {
		offsetY := int32(i*(itemHeight()+paddingY) + paddingY + borderWidth)
		if err := itm.render(r, offsetY); err != nil {
			return errors.Annotate(err, "render item")
		}
	}

	return nil
}

func drawWindowBackground(r *sdl.Renderer) error {
	if err := setDrawColor(r, theme.Background); err != nil {
		return errors.Annotate(err, "set background color")
	}
	backgroundFrame := &sdl.Rect{
		X: 0,
		Y: 0,
		W: int32(windowWidth),
		H: int32(windowHeight),
	}
	if err := r.FillRect(backgroundFrame); err != nil {
		return errors.Annotate(err, "draw background")
	}
	return nil
}

func drawWindowBorder(r *sdl.Renderer, width int32) error {
	if err := setDrawColor(r, theme.ItemBackground); err != nil {
		return errors.Annotate(err, "set border color")
	}
	for i := int32(1); i <= width; i++ {
		borderFrame := &sdl.Rect{
			X: i,
			Y: i,
			W: int32(windowWidth) - i*2,
			H: int32(windowHeight) - i*2,
		}
		if err := r.DrawRect(borderFrame); err != nil {
			return errors.Annotate(err, "draw border")
		}
	}

	return nil
}

func useFont(name string) error {
	fontPath, err := fonts.Find(name)
	if err != nil {
		return err
	}
	if fontPath == "" {
		return errors.New("font not found")
	}

	font, err = ttf.OpenFont(fontPath, fontSize)
	if err != nil {
		return errors.Annotatef(err, "load font %s", name)
	}

	return nil
}

func itemHeight() int {
	return fontSize + 2*paddingY
}

func setDrawColor(r *sdl.Renderer, c sdl.Color) error {
	return r.SetDrawColor(c.R, c.G, c.B, c.A)
}
