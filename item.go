package menu

import (
	"strings"

	"github.com/juju/errors"
	"github.com/veandco/go-sdl2/sdl"
)

type menuItem struct {
	ID            string     `json:"id"`
	Label         *string    `json:"label"`
	LabelCommand  *command   `json:"label_cmd"`
	ActionKey     string     `json:"key"`
	ActionCommand *command   `json:"action_cmd"`
	Toggle        *toggle    `json:"switch"`
	Items         []menuItem `json:"items"`
	Invalidates   []string   `json:"invalidates"`

	active bool
}

func (mi *menuItem) prepare() {
	if mi.LabelCommand != nil {
		mi.LabelCommand.keepUpdated()
	}
	if mi.Toggle != nil {
		mi.Toggle.StateCommand.keepUpdated()
	}
	for i := range mi.Items {
		mi.Items[i].prepare()
	}
}

func (mi *menuItem) trigger() {
	if mi.ActionCommand != nil && !mi.ActionCommand.busy {
		mi.ActionCommand.exec()
	}
	if mi.Toggle != nil {
		mi.Toggle.trigger()
	}
}

func (mi *menuItem) label() string {
	if cmd := mi.LabelCommand; cmd != nil {
		// pretty.Println(cmd)
		if cmd.error != nil {
			return cmd.error.Error()
		}
		if cmd.out != nil {
			return *cmd.out
		}
	}
	if cmd := mi.ActionCommand; cmd != nil {
		if cmd.error != nil {
			return cmd.error.Error()
		}
	}
	if mi.Toggle != nil {
		return mi.Toggle.label()
	}
	if mi.Label != nil {
		return *mi.Label
	}
	return "No name"
}

func (mi *menuItem) busy() bool {
	if mi.LabelCommand != nil && mi.LabelCommand.busy {
		return true
	}
	if mi.ActionCommand != nil && mi.ActionCommand.busy {
		return true
	}
	if mi.Toggle != nil && mi.Toggle.StateCommand.busy {
		return true
	}
	return false
}

func (mi *menuItem) render(r *sdl.Renderer, offsetY int32) error {
	if mi.ActionKey != "" {
		if err := drawItemBackground(r, offsetY, mi.active); err != nil {
			return errors.Annotate(err, "draw item background")
		}
	}
	if err := drawItemLabel(r, offsetY, *mi); err != nil {
		return errors.Annotate(err, "draw item label")
	}
	if err := drawActionKeyLabel(r, offsetY, *mi); err != nil {
		return errors.Annotate(err, "draw item label")
	}

	return nil
}

type toggle struct {
	StateCommand command             `json:"state_cmd"`
	States       map[string]menuItem `json:"states"`
}

func (t *toggle) label() string {
	if t.StateCommand.error != nil {
		return t.StateCommand.error.Error()
	}
	if t.StateCommand.out != nil {
		if opt, ok := t.States[*t.StateCommand.out]; ok {
			return opt.label()
		}
	}
	return "No name"
}

func (t *toggle) trigger() {
	if t.StateCommand.busy {
		return
	}

	if t.StateCommand.out != nil {
		if opt, ok := t.States[*t.StateCommand.out]; ok {
			opt.trigger()
			t.StateCommand.resetTimer()
		}
	}
}

func drawItemBackground(r *sdl.Renderer, offsetY int32, active bool) error {
	if err := setDrawColor(r, theme.ItemBackground); err != nil {
		return errors.Annotate(err, "set item background color")
	}

	actionKeyFrame := &sdl.Rect{
		X: int32(paddingX + borderWidth),
		Y: offsetY,
		W: int32(itemHeight()),
		H: int32(itemHeight()),
	}
	var err error
	if active {
		err = r.DrawRect(actionKeyFrame)
	} else {
		err = r.FillRect(actionKeyFrame)
	}
	if err != nil {
		return errors.Annotate(err, "draw item")
	}

	itemBackground := &sdl.Rect{
		X: int32(borderWidth + paddingX + itemHeight() + paddingX),
		Y: offsetY,
		W: int32(windowWidth - 3*paddingX - itemHeight() - 2*borderWidth),
		H: int32(itemHeight()),
	}

	if active {
		err = r.DrawRect(itemBackground)
	} else {
		err = r.FillRect(itemBackground)
	}
	if err != nil {
		return errors.Annotate(err, "draw item")
	}

	return nil
}

func drawItemLabel(r *sdl.Renderer, offsetY int32, mi menuItem) error {
	label := mi.label()
	color := theme.ItemText
	if mi.busy() && mi.ActionKey != "" {
		label += " [busy]"
		color = theme.TextBusy
	}
	if label == "" {
		return nil
	}

	labelTexture, err := renderText(r, label, color)
	if err != nil {
		return errors.Annotate(err, "draw item label")
	}
	defer labelTexture.Destroy()
	lw, lh, err := font.SizeUTF8(label)
	if err != nil {
		return err
	}

	const magicNumber = 3
	itemLabel := &sdl.Rect{
		X: int32(paddingX*2+itemHeight()+(itemHeight()-lh)/2) + magicNumber,
		Y: offsetY + int32(itemHeight()-lh)/2,
		W: int32(lw),
		H: int32(lh),
	}
	if err := r.Copy(labelTexture, nil, itemLabel); err != nil {
		return errors.Annotate(err, "render item label")
	}
	return nil
}

func drawActionKeyLabel(r *sdl.Renderer, offsetY int32, mi menuItem) error {
	label := strings.ToUpper(mi.ActionKey)
	color := theme.ItemText
	if mi.busy() && mi.ActionKey != "" {
		color = theme.TextBusy
	}
	if label == "" {
		return nil
	}

	labelTexture, err := renderText(r, label, color)
	if err != nil {
		return errors.Annotate(err, "draw item label")
	}
	defer labelTexture.Destroy()
	lw, lh, err := font.SizeUTF8(label)
	if err != nil {
		return err
	}

	itemLabel := &sdl.Rect{
		X: int32(paddingX + (itemHeight()-lw)/2),
		Y: offsetY + int32(itemHeight()-lh)/2,
		W: int32(lw),
		H: int32(lh),
	}
	if err := r.Copy(labelTexture, nil, itemLabel); err != nil {
		return errors.Annotate(err, "render item label")
	}
	return nil
}

func renderText(r *sdl.Renderer, text string, c sdl.Color) (t *sdl.Texture, err error) {
	fs, err := font.RenderUTF8Blended(text, c)
	if err != nil {
		return nil, err
	}
	defer fs.Free()

	t, err = r.CreateTextureFromSurface(fs)
	if err != nil {
		return nil, err
	}

	return
}
