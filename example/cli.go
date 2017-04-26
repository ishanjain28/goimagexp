package main

import (
	"github.com/ishanjain28/goimagexp"
	"image/png"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	str := []string{
		imagexp.DESATURATION,
		imagexp.BASIC,
		imagexp.BASICIMPROVED,
		imagexp.SINGLEBLUE,
		imagexp.SINGLEGREEN,
		imagexp.SINGLERED,
		imagexp.REDONLYFILTER,
		imagexp.GREENONLYFILTER,
		imagexp.BLUEONLYFILTER,
		imagexp.DECOMPOSITIONMAX,
		imagexp.DECOMPOSITIONMIN,
	}
	for _, v := range str {
		wg.Add(1)
		create(v)
	}

	wg.Wait()
}

func create(v string) {
	image := imagexp.TransformImage(v, "image.png")

	file, _ := os.Create(v + ".png")

	defer file.Close()

	png.Encode(file, image)
	wg.Done()
}
