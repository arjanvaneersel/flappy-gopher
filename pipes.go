package main

import (
	"fmt"
	"math/rand"
	"time"

	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipes struct {
	mu sync.RWMutex

	texture *sdl.Texture
	speed   int32

	pipes []*pipe
}

type pipe struct {
	mu sync.RWMutex

	x, width, height int32

	inverted bool
}

func newPipes(r *sdl.Renderer) (*pipes, error) {
	t, err := img.LoadTexture(r, "res/imgs/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("Error while rendering pipe image: %v", err)
	}

	ps := &pipes{
		texture: t,
		speed:   2,
	}

	go func() {
		for {
			ps.mu.Lock()
			ps.pipes = append(ps.pipes, newPipe())
			ps.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return ps, nil
}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(r, ps.texture); err != nil {
			return err
		}
	}
	return nil
}

func (ps *pipes) restart() {
	ps.mu.Lock()
	ps.mu.Unlock()

	ps.pipes = nil
}

func (ps *pipes) update() {
	ps.mu.Lock()
	ps.mu.Unlock()

	var rem []*pipe
	for _, p := range ps.pipes {
		p.update(ps.speed)
		if p.x+p.width > 0 {
			rem = append(rem, p)
		}
	}
	ps.pipes = rem
}

func (ps *pipes) destroy() {
	ps.mu.Lock()
	ps.mu.Unlock()

	ps.texture.Destroy()
}

func (ps *pipes) touch(b *bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		p.touch(b)
	}
}

func newPipe() *pipe {
	return &pipe{
		x:        800,
		height:   100 + int32(rand.Intn(300)),
		width:    50,
		inverted: rand.Float32() > 0.5,
	}
}

func (p *pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rect := &sdl.Rect{W: p.width, H: p.height, X: p.x, Y: 600 - p.height}
	flip := sdl.FLIP_NONE
	if p.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL

	}

	err := r.CopyEx(texture, nil, rect, 0, nil, flip)
	if err != nil {
		return fmt.Errorf("Error while copying background: %v", err)
	}
	return nil
}

func (p *pipe) update(speed int32) {
	p.mu.Lock()
	p.mu.Unlock()

	p.x -= speed
}

func (p *pipe) touch(b *bird) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	b.touch(p)
}
