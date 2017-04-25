package main

import (
	"fmt"
	"github.com/ishanjain28/goimagexp"
	"os"
	"image/png"
)

func main() {
	fmt.Println("hello")

	finalImage := imagexp.TransformImage(imagexp.BLUEONLYFILTER, "image.jpg")

	final_file, _ := os.Create("test.png")

	png.Encode(final_file, finalImage)
}