package main

import (
	"image/color"
	"log"
	"math/rand"
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
	radius int
	xv, yv int
	c      color.Color
}

func newBall(x, y, r, xv, yv int, c color.Color) *ball {
	return &ball{
		center: center{
			x: x,
			y: y,
		},
		radius: r,
		xv:     xv,
		yv:     yv,
		c:      c,
	}
}

func (b *ball) draw(pixels []byte) {
	for y := -b.radius; y <= b.radius; y++ {
		for x := -b.radius; x <= b.radius; x++ {
			if x*x+y*y <= b.radius*b.radius {
				fillPixel(pixels, x+b.x, y+b.y, b.c)
			}
		}
	}
}

func (b *ball) update() {
	b.x += b.xv
	b.y += b.yv
	if b.x-b.radius == 0 || b.x+b.radius == width-1 {
		b.xv = -b.xv
	}
	if b.y-b.radius == 0 || b.y+b.radius == height-1 {
		b.yv = -b.yv
	}
}

func init() {
	log.SetFlags(log.Lshortfile)
	rand.Seed(time.Now().Unix())
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

	ball := newBall(rand.Intn(740)+30, rand.Intn(540)+30, 30, 1, 1, color.White)

	for range time.NewTicker(10 * time.Millisecond).C {
		pixels := make([]byte, width*height*4)
		r.Clear()
		ball.update()
		ball.draw(pixels)
		texture.Update(nil, pixels, width*4)
		r.Copy(texture, nil, nil)
		r.Present()
	}
}

func fillPixel(pixels []byte, x, y int, c color.Color) {
	r, g, b, a := c.RGBA()
	index := 4 * (y*width + x)
	pixels[index] = byte(r)
	pixels[index+1] = byte(g)
	pixels[index+2] = byte(b)
	pixels[index+3] = byte(a)
}
