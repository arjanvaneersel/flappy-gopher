package main

import (
	"fmt"
	"os"
	"time"

	"runtime"

	img "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("Error while initializing SDL: %v", err)
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		return fmt.Errorf("Could not initialize TTF: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("Error while creating window: %v", err)
	}
	defer w.Destroy()

	err = drawTitle(r, "Flappy Gopher")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	time.Sleep(5 * time.Second)

	s, err := newScene(r)
	if err != nil {
		return fmt.Errorf("Error while getting new scene: %v", err)
	}
	defer s.destroy()

	events := make(chan sdl.Event)
	errc := s.run(events, r)

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}

func drawTitle(r *sdl.Renderer, text string) error {
	r.Clear()

	f, err := ttf.OpenFont("res/fonts/flappy.ttf", 20)
	if err != nil {
		return fmt.Errorf("Error while opening font: %v", err)
	}
	defer f.Close()

	s, err := f.RenderUTF8Solid(text, sdl.Color{R: 255, G: 100, B: 0, A: 255})
	if err != nil {
		return fmt.Errorf("Error while rendering title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("Error while creating texture: %v", err)
	}
	defer t.Destroy()

	err = r.Copy(t, nil, nil)
	if err != nil {
		return fmt.Errorf("Error while copying texture: %v", err)
	}
	defer r.Present()

	return nil
}

func drawBackground(r *sdl.Renderer) error {
	r.Clear()
	t, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return fmt.Errorf("Error while rendering background: %v", err)
	}
	defer t.Destroy()

	err = r.Copy(t, nil, nil)
	if err != nil {
		return fmt.Errorf("Error while copying background: %v", err)
	}

	r.Present()
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
