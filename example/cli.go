package main

import (
	"github.com/ishanjain28/goimagexp"
	"os"
	"image/png"
	"log"
)

func main() {

	//Enter the image path and name of transformation to apply on that image
	finalImage := imagexp.TransformImage(imagexp.DECOMPOSITIONMAX, "image.jpg")

	// finalImage is of type image.Image
	//You can use it wherever you want.
	//Here I am storing it in a file
	final_file, err := os.Create("final_image.png")
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer final_file.Close()
	//Encode image.Image into proper png format and store it in final_file
	png.Encode(final_file, finalImage)
}
