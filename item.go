package menu

import (
	"log"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/veandco/go-sdl2/sdl"
)

type listItem struct {
	Label     string       `json:"label"`
	ActionKey string       `json:"key"`
	Keep      bool         `json:"keep"`
	Command   string       `json:"command"`
	Condition *ifCondition `json:"if"`
	Items     []listItem   `json:"items"`

	active bool

	subLabel string
}

type ifCondition struct {
	Command string              `json:"command"`
	Output  map[string]listItem `json:"output"`

	cachedOutput *string
	busy         bool
}

func (li *listItem) render(r *sdl.Renderer, offsetY int32) error {
	if err := drawItemBackground(r, offsetY, li.active); err != nil {
		return errors.Annotate(err, "draw item background")
	}
	if err := drawItemLabel(r, offsetY, *li); err != nil {
		return errors.Annotate(err, "draw item label")
	}
	if err := drawActionKeyLabel(r, offsetY, *li); err != nil {
		return errors.Annotate(err, "draw item label")
	}

	return nil
}

func (li *listItem) call() error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Action failed: %v", err)
		}
	}()

	switch {
	case li.Condition != nil:
		if li.Condition.cachedOutput == nil {
			return nil
		}

		out := *li.Condition.cachedOutput
		if match, ok := li.Condition.Output[out]; ok {
			_, err := execCommand(match.Command)
			return err
		}

		log.Println("no match", out)
		return nil
	case li.Command != "":
		_, err := execCommand(li.Command)
		return err
	default:
		return nil
	}
}

func (li *listItem) label() string {
	if li.subLabel != "" {
		return li.Label + " " + li.subLabel
	}
	return li.Label
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

func drawItemLabel(r *sdl.Renderer, offsetY int32, li listItem) error {
	label := li.Label
	if li.Condition != nil {
		if li.Condition.cachedOutput != nil {
			out := *li.Condition.cachedOutput
			if match, ok := li.Condition.Output[out]; ok {
				label = match.Label
			}
		}
		if li.Condition.busy {
			label += " [busy]"
		}
	}
	color := theme.ItemText
	if li.Condition != nil && li.Condition.busy {
		color = theme.TextBusy
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

func drawActionKeyLabel(r *sdl.Renderer, offsetY int32, li listItem) error {
	label := strings.ToUpper(li.ActionKey)
	color := theme.ItemText
	if li.Condition != nil && li.Condition.busy {
		color = theme.TextBusy
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

func (c *ifCondition) evaluate() {
	fn := func() {
		c.busy = true
		out, err := execCommand(c.Command)
		if err != nil {
			log.Printf("Error evaluating condition: %v", err)
		} else {
			out = strings.TrimSpace(out)
			c.cachedOutput = &out
		}
		c.busy = false
	}
	fn()

	t := time.NewTicker(2 * time.Second)
	defer t.Stop()
	for range t.C {
		fn()
	}
}
