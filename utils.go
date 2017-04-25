package main

import (
	"image/png"
	"path"
	"io"
	"image/jpeg"
	"log"
	"bytes"
	"fmt"
)

func Basic(r, g, b uint32) float64 {
	return float64((r + g + b) / 3)
}

func BasicImproved(r, g, b uint32) float64 {
	return float64(0.3)*float64(r) + float64(0.59)*float64(g) + float64(0.11)*float64(b)
}

func Desaturation(r, g, b uint32) float64 {
	return float64(MaxOfThree(r, g, b)+MinOfThree(r, g, b)) / 2
}

func MaxOfThree(r, g, b uint32) uint32 {
	return Max(Max(r, g), b)
}

func MinOfThree(r, g, b uint32) uint32 {
	return Min(Min(r, g), b)
}

func Max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func Min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// Use file Extension to idenitfy which decoder to use.
func (img *Image) IdentifyDecoder(imageFile io.Reader) {

	fileExtension := path.Ext(img.srcPath)
	var err error

	if fileExtension == ".png" {
		img.srcImage, err = png.Decode(imageFile)
		if err != nil {
			log.Fatalf("%s", err)
		}
	} else if fileExtension == ".jpg" || fileExtension == ".jpeg" {

		var jpegBugger bytes.Buffer

		decodedJPEG, err := jpeg.Decode(imageFile)

		if err != nil {
			log.Fatalf("%s", err)
		}

		png.Encode(&jpegBugger, decodedJPEG)

		img.srcImage, err = png.Decode(&jpegBugger)

		fmt.Println(img.srcImage.At(100, 100))
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

}
