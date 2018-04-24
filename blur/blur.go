package blur

import (
	"log"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width  = 800
	height = 600
)

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

	pixels := make([]byte, width*height*4)

	for range time.NewTicker(100 * time.Millisecond).C {
		r.Clear()
		fillPixels(pixels)
		texture.Update(nil, pixels, width*4)
		r.Copy(texture, nil, nil)
		r.Present()
	}
}

func fillPixels(pixels []byte) {
	for i := 0; i != len(pixels); i++ {
		if i%4 == 3 {
			pixels[i] = 255
		} else {
			pixels[i] = byte(rand.Intn(256))
		}
	}
}
