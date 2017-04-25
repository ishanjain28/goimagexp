package main

import (
	"image/png"
	"path"
	"io"
	"image/jpeg"
	"log"
	"bytes"
	"fmt"
	"image"
	"strings"
	"os"
)


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

func (img *Image) Save(FilterName string, finalImage image.Image) {
	srcFilename := path.Base(img.srcPath)
	srcFilename = strings.Replace(srcFilename, path.Ext(img.srcPath), "", -1)
	//fileExtension := path.Ext(img.srcPath)
	FilterName = strings.Replace(FilterName, " ", "", -1)

	os.Mkdir("final_Image", 0777)
	os.Mkdir(srcFilename, 0777)

	outPath := path.Join(srcFilename, FilterName+".png")
	outfile, err := os.Create(outPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer outfile.Close()

	png.Encode(outfile, finalImage)
}
