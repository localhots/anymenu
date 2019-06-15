package menu

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var frameLimit = 60

// Main is the main function.
func Main() {
	conf := flag.String("config", "config.json", "Path to config file")
	flag.IntVar(&windowWidth, "width", windowWidth, "Window width")
	flag.IntVar(&windowHeight, "height", windowHeight, "Window height")
	flag.IntVar(&borderWidth, "border", borderWidth, "Border width")
	flag.IntVar(&frameLimit, "fps", frameLimit, "FPS limit")
	flag.StringVar(&fontName, "fontname", fontName, "Font name (must be TrueType)")
	flag.IntVar(&fontSize, "fontsize", fontSize, "Font size")
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.Println("Loading config")
	items, err := getItems(*conf)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	keepConditionsEvaluated(items)

	log.Println("Loading font")
	if err := ttf.Init(); err != nil {
		log.Fatalf("Failed to initialize TrueType package: %v", err)
	}
	defer ttf.Quit()
	if err := useFont(fontName); err != nil {
		log.Fatalf("Failed to open font: %v", err)
	}

	log.Println("Initializing SDL")
	if err := sdl.Init(sdl.INIT_EVENTS); err != nil {
		log.Fatalf("Failed to initialize sdl: %v", err)
	}
	defer sdl.Quit()
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	log.Println("Creating window")
	sdlWindow, err := sdl.CreateWindow("menu",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(windowWidth), int32(windowHeight),
		sdl.WINDOW_SHOWN|sdl.WINDOW_BORDERLESS|sdl.WINDOW_ALLOW_HIGHDPI|sdl.WINDOW_OPENGL,
	)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}
	defer sdlWindow.Destroy()

	log.Println("Creating renderer")
	renderer, err := sdl.CreateRenderer(sdlWindow, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.Fatalf("Failed to create renderer: %v", err)
	}
	defer renderer.Destroy()

	w := window{
		window:   sdlWindow,
		allItems: items,
		curItems: items,
	}
	w.eventLoop(renderer)
}

func (w *window) eventLoop(r *sdl.Renderer) {
	// Throttle rendering
	ticker := time.NewTicker(time.Second / time.Duration(frameLimit))
	defer ticker.Stop()

	for range ticker.C {
		event := sdl.PollEvent()
		switch tevt := event.(type) {
		case *sdl.QuitEvent:
			log.Println("Shutting down")
			return
		case *sdl.KeyboardEvent:
			if err := w.handleKeyEvent(tevt); err != nil {
				log.Printf("Failed to process key event: %v", err)
			}
		}

		if err := w.render(r); err != nil {
			log.Printf("Failed to render window: %v", err)
		}
		// Ta-dah!
		r.Present()
	}
}

func (w *window) handleKeyEvent(e *sdl.KeyboardEvent) error {
	if e.Keysym.Sym == sdl.K_ESCAPE {
		w.curItems = w.allItems
		return nil
	}

	// Try and transform a key code into a string containing a letter
	key := string(rune(e.Keysym.Sym))
	for i, itm := range w.curItems {
		if itm.ActionKey == key {
			if len(itm.Items) > 0 {
				w.curItems = itm.Items
				return nil
			}

			w.curItems[i].active = (e.State == sdl.PRESSED)
			if e.State == sdl.RELEASED {
				if itm.Condition == nil || !itm.Condition.busy {
					go w.curItems[i].call()
				}
			}
		}
	}

	return nil
}

func getItems(path string) ([]listItem, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var items []listItem
	err = json.Unmarshal(b, &items)
	return items, err
}

func keepConditionsEvaluated(items []listItem) {
	for _, item := range items {
		if item.Condition != nil {
			go item.Condition.evaluate()
		}
		keepConditionsEvaluated(item.Items)
	}
}
