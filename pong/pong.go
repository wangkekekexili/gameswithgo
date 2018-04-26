package main

import (
	"image/color"
	"log"

	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width  = 800
	height = 600
)

type center struct {
	x, y int
}

type ball struct {
	center
	r      int
	xv, yv int
	c      color.Color
}

func newBall(x, y, r, xv, yv int, c color.Color) *ball {
	return &ball{
		center: center{
			x: x,
			y: y,
		},
		r:  r,
		xv: xv,
		yv: yv,
		c:  c,
	}
}

func (b *ball) draw(pixels []byte) {
	for y := -b.r; y <= b.r; y++ {
		for x := -b.r; x <= b.r; x++ {
			if x*x+y*y <= b.r*b.r {
				fillPixel(pixels, x+b.x, y+b.y, b.c)
			}
		}
	}
}

func (b *ball) leftTouch(p *paddle) bool {
	if b.y+b.r < p.y-p.h/2 {
		return false
	}
	if b.y-b.r > p.y+p.h/2 {
		return false
	}
	return b.x-b.r <= p.x+p.w/2
}

func (b *ball) rightTouch(p *paddle) bool {
	if b.y+b.r < p.y-p.h/2 {
		return false
	}
	if b.y-b.r > p.y+p.h/2 {
		return false
	}
	return b.x+b.r >= p.x-p.w/2
}

func (b *ball) update(left, right *paddle) {
	b.x += b.xv
	b.y += b.yv
	if b.y-b.r <= 0 {
		b.y = b.r
		b.yv = -b.yv
	}
	if b.y+b.r >= height-1 {
		b.y = height - 1 - b.r
		b.yv = -b.yv
	}
	if b.leftTouch(left) {
		b.xv = -b.xv
	}
	if b.rightTouch(right) {
		b.xv = -b.xv
	}
	if b.x-b.r <= 0 || b.x+b.r > width {
		b.x = 400
		b.y = 300
	}
}

type paddle struct {
	center
	w int
	h int
	c color.Color
}

func newPaddle(x, y, w, h int, c color.Color) *paddle {
	return &paddle{
		center: center{
			x: x,
			y: y,
		},
		w: w,
		h: h,
		c: c,
	}
}

func (p *paddle) draw(pixels []byte) {
	startX, endX := p.x-p.w/2, p.x+p.w/2
	startY, endY := p.y-p.h/2, p.y+p.h/2
	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			fillPixel(pixels, x, y, color.White)
		}
	}
}

func (p *paddle) update(keyboardState []uint8) {
	if keyboardState[sdl.SCANCODE_UP] != 0 {
		if p.y > p.h/2+2 {
			p.y -= 3
		}
	}
	if keyboardState[sdl.SCANCODE_DOWN] != 0 {
		if p.y < height-p.h/2-3 {
			p.y += 3
		}
	}
}

func (p *paddle) aiUpdate(b *ball) {
	p.y = b.y
	if p.y < p.h/2 {
		p.y = p.h / 2
	}
	if p.y > height-1-p.h/2 {
		p.y = height - 1 - p.h/2
	}
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Destroy()

	texture, err := r.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		log.Fatal(err)
	}
	defer texture.Destroy()

	ball := newBall(400, 300, 10, 2, 2, color.White)
	player1 := newPaddle(5, 50, 10, 100, color.White)
	player2 := newPaddle(width-6, 50, 10, 100, color.White)
	state := sdl.GetKeyboardState()
	for range time.NewTicker(10 * time.Millisecond).C {
		sdl.PollEvent()

		pixels := make([]byte, width*height*4)
		r.Clear()

		ball.update(player1, player2)
		player1.update(state)
		player2.aiUpdate(ball)

		ball.draw(pixels)
		player1.draw(pixels)
		player2.draw(pixels)

		texture.Update(nil, pixels, width*4)
		r.Copy(texture, nil, nil)
		r.Present()
	}
}

func fillPixel(pixels []byte, x, y int, c color.Color) {
	r, g, b, a := c.RGBA()
	index := 4 * (y*width + x)
	if index <= len(pixels)-4 {
		pixels[index] = byte(r)
		pixels[index+1] = byte(g)
		pixels[index+2] = byte(b)
		pixels[index+3] = byte(a)
	}
}
