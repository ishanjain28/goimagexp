package imagexp

import (
	"image/color"
	"image"
	"io"
	"path"
	"image/png"
	"log"
	"bytes"
	"image/jpeg"
	"strings"
	"os"
	"math"
	"image/draw"
	"fmt"
	"bufio"
)

type Image struct {
	path         string
	width        int
	height       int
	decodedImage image.Image
}

type grayImage struct {
	Image
	Options GrayOptions
}

type colorImage struct {
	Image
	Filter  string
	Options ColorOptions
}

type GrayOptions interface {
}
type ColorOptions interface {
}

var PARTS int = 50

const (
	BASIC            = "basic"
	BASICIMPROVED    = "basic.improved"
	DESATURATION     = "desaturation"
	DECOMPOSITIONMAX = "decomposition.max"
	DECOMPOSITIONMIN = "decomposition.min"
	SINGLERED        = "single.channel.red"
	SINGLEGREEN      = "single.channel.green"
	SINGLEBLUE       = "single.channel.blue"
	REDONLYFILTER    = "red.only"
	GREENONLYFILTER  = "green.only"
	BLUEONLYFILTER   = "blue.only"
)

func TransformImage(transformationName string, path string) {
	//Create a new instance of struct
	img := &Image{}
	//Set Path
	img.path = path

	//Decode
	img.decodedImage = img.Decode()
	//Set Image Dimension
	img.SetDimension(img.decodedImage.Bounds().Max.X, img.decodedImage.Bounds().Max.Y)

	var cImage colorImage
	//var gImage grayImage

	switch transformationName {
	case BASIC:
	case BASICIMPROVED:
	case DESATURATION:
	case DECOMPOSITIONMAX:
	case DECOMPOSITIONMIN:
	case SINGLERED:
	case SINGLEGREEN:
	case SINGLEBLUE:
	case REDONLYFILTER:
		imga := cImage.Create(REDONLYFILTER)
		fmt.Println(imga)
	case GREENONLYFILTER:
	case BLUEONLYFILTER:
	}
}

func (cImage colorImage) Create(FilterName string) *image.RGBA64 {
	fmt.Println(cImage)
	newColorImage := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})
	newGrayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})

	finalImage := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})

	partLimit := cImage.height / PARTS
	difference := cImage.height % PARTS

	switch FilterName {
	case REDONLYFILTER:

		for j := 0; j < PARTS; j++ {
			startFromRow := partLimit * j
			upToRow := partLimit * (j + 1)
			if j == PARTS-1 {
				upToRow += difference
			}

			go cImage.applyTransformation(startFromRow, upToRow, newColorImage, newGrayImage, redFilter)
		}

	case GREENONLYFILTER:

		for j := 0; j < PARTS; j++ {
			startFromRow := partLimit * j
			upToRow := partLimit * (j + 1)
			if j == PARTS-1 {
				upToRow += difference
			}
			go cImage.applyTransformation(startFromRow, upToRow, newColorImage, newGrayImage, redFilter)
		}

	case BLUEONLYFILTER:
		for j := 0; j < PARTS; j++ {
			startFromRow := partLimit * j
			upToRow := partLimit * (j + 1)
			if j == PARTS-1 {
				upToRow += difference
			}
			go cImage.applyTransformation(startFromRow, upToRow, newColorImage, newGrayImage, redFilter)
		}
	}

	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}}, newGrayImage, image.Point{0, 0}, draw.Src)
	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}}, newColorImage, image.Point{0, 0}, draw.Over)

	return finalImage
}

func (cImage *colorImage) applyTransformation(startFromRow, upToRow int, colorImage *image.RGBA64, grayImage *image.Gray16, transformationFunction func(r, g, b, a uint32) color.RGBA64) {
	for i := startFromRow; i <= upToRow; i++ {
		for j := 0; j <= cImage.width; j++ {
			point := cImage.decodedImage.At(i, j)
			r, g, b, a := point.RGBA()

			pixelColor := transformationFunction(r, g, b, a)
			grayAVG := basicImproved(r, g, b)
			grayColor := color.Gray16{uint16(math.Ceil(grayAVG))}
			grayImage.SetGray16(i, j, grayColor)
			colorImage.SetRGBA64(i, j, pixelColor)
		}
	}
}

func (gImage *grayImage) applyTransformation(x, y int) {

}

func (gImage grayImage) Create(FilterName string) {
	switch FilterName {
	case "basic":
	case "basic.improved":
	case "desaturation":
	case "decomposition.max":
	case "decomposition.min":
	case "single.channel.red":
	case "single.channel.green":
	case "single.channel.blue":

	}
}

func (img *Image) SetDimension(width int, height int) {
	img.width = width
	img.height = height
}

func (img *Image) Decode() image.Image {
	imageFile, err := os.Open(img.path)
	if err != nil {
		log.Fatalf("Error Occurred in opening file: %s", err)
	}
	defer imageFile.Close()

	fileExtension := path.Ext(img.path)
	var decodedImage image.Image

	switch fileExtension {
	case ".png":
		decodedImage, err = png.Decode(imageFile)
		if err != nil {
			log.Fatalf("Error in decoding png: %s", err)
		}
	case ".jpg", ".jpeg":
		var jpegBuffer bytes.Buffer

		decodedJPEG, err := jpeg.Decode(imageFile)
		if err != nil {
			log.Fatalf("Error in decoding jpeg: %s", err)
		}

		png.Encode(&jpegBuffer, decodedJPEG)

		decodedImage, err = png.Decode(&jpegBuffer)

		if err != nil {
			log.Fatalf("Error in encoding jpeg to png: %s", err)
		}
	}

	return decodedImage
}

// Use file Extension to idenitfy which decoder to use.
func (img *Image) Save(SaveDir string, finalImage image.Image, shouldCreateDir bool) {
	destFileName := path.Base(img.path)
	destFileName = strings.Replace(destFileName, path.Ext(img.path), "", -1)
	//fileExtension := path.Ext(img.srcPath)
	SaveDir = strings.Replace(SaveDir, " ", "", -1)

	if shouldCreateDir {
		os.Mkdir(SaveDir, 0777)
	}

	outPath := path.Join(SaveDir, destFileName+".png")
	outfile, err := os.Create(outPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer outfile.Close()

	png.Encode(outfile, finalImage)
}

func basic(r, g, b uint32) float64 {
	return float64((r + g + b) / 3)
}

func basicImproved(r, g, b uint32) float64 {
	return float64(0.3)*float64(r) + float64(0.59)*float64(g) + float64(0.11)*float64(b)
}

func desaturation(r, g, b uint32) float64 {
	return float64(maxOfThree(r, g, b)+minOfThree(r, g, b)) / 2
}

func maxOfThree(r, g, b uint32) uint32 {
	return Max(Max(r, g), b)
}

func minOfThree(r, g, b uint32) uint32 {
	return Min(Min(r, g), b)
}

func redFilter(r, g, b, a uint32) color.RGBA64 {

	if !(r > b) || !(r > g) {
		return color.RGBA64{uint16(255), uint16(255), uint16(255), uint16(0)}
	}

	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func greenFilter(r, g, b, a uint32) color.RGBA64 {
	if !(g > r) || !(g > b) {
		return color.RGBA64{uint16(255), uint16(255), uint16(255), uint16(0)}
	}

	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func blueFilter(r, g, b, a uint32) color.RGBA64 {
	if !(b > g) || !(b > r) {
		return color.RGBA64{uint16(255), uint16(255), uint16(255), uint16(0)}
	}

	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}