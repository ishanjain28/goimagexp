package imagexp

/*
1. Performance Boost: Use go routines to access and modify pixels in multiple locations at the same time
2. Organize code
3. Add filters.
*/

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path"
	"strings"
)

type Image struct {
	path         string
	width        int
	height       int
	decodedImage image.Image
}

type grayImage struct {
	Image
}

type colorImage struct {
	Image
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

func TransformImage(transformationName string, path string) image.Image {

	var finalImage interface{}
	//Create a new instance of struct
	img := &Image{}
	//Set Path
	img.path = path

	//Decode
	img.decodedImage = img.Decode()
	//Set Image Dimension
	img.SetDimension(img.decodedImage.Bounds().Max.X, img.decodedImage.Bounds().Max.Y)
	//Print a message about Image Dimension
	fmt.Printf("Image Resolution: %dx%d\n", img.width, img.height)

	var cImage colorImage
	var gImage grayImage

	switch transformationName {

	//Improve this
	case BASIC:
		gImage = grayImage{Image{img.path, img.width, img.height, img.decodedImage}}
		finalImage = gImage.Create(BASIC)
	case BASICIMPROVED:
		gImage = grayImage{Image{img.path, img.width, img.height, img.decodedImage}}
		finalImage = gImage.Create(BASICIMPROVED)
	case DESATURATION:
		gImage = grayImage{Image{img.path, img.width, img.height, img.decodedImage}}
		finalImage = gImage.Create(DESATURATION)
	case DECOMPOSITIONMAX:
		gImage = grayImage{Image{img.path, img.width, img.height, img.decodedImage}}
		finalImage = gImage.Create(DECOMPOSITIONMAX)
	case DECOMPOSITIONMIN:
		gImage = grayImage{Image{img.path, img.width, img.height, img.decodedImage}}
		finalImage = gImage.Create(DECOMPOSITIONMIN)
	case SINGLERED:
		gImage = grayImage{Image{img.path, img.width, img.height, img.decodedImage}}
		finalImage = gImage.Create(SINGLERED)
	case SINGLEGREEN:
		gImage = grayImage{Image{img.path, img.width, img.height, img.decodedImage}}
		finalImage = gImage.Create(SINGLEGREEN)
	case SINGLEBLUE:
		gImage = grayImage{Image{img.path, img.width, img.height, img.decodedImage}}
		finalImage = gImage.Create(SINGLEBLUE)
	case REDONLYFILTER:
		cImage = colorImage{Image{img.path, img.width, img.height, img.decodedImage}}
		finalImage = cImage.Create(REDONLYFILTER)
	case GREENONLYFILTER:
		cImage = colorImage{Image{img.path, img.width, img.height, img.decodedImage}}
		finalImage = cImage.Create(GREENONLYFILTER)
	case BLUEONLYFILTER:
		cImage = colorImage{Image{img.path, img.width, img.height, img.decodedImage}}
		finalImage = cImage.Create(BLUEONLYFILTER)

	}

	return finalImage.(image.Image)
}

func (gImage *grayImage) Create(FilterName string) *image.Gray16 {

	newGrayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{gImage.width, gImage.height}})
	switch FilterName {
	case BASIC:
		gImage.applyTransformation(newGrayImage, basic)
	case BASICIMPROVED:
		gImage.applyTransformation(newGrayImage, basicImproved)
	case DESATURATION:
		gImage.applyTransformation(newGrayImage, desaturation)
	case DECOMPOSITIONMAX:
		gImage.applyTransformation(newGrayImage, decompositionMax)
	case DECOMPOSITIONMIN:
		gImage.applyTransformation(newGrayImage, decompositionMin)
	case SINGLERED:
		gImage.applyTransformation(newGrayImage, singleChannelRed)
	case SINGLEGREEN:
		gImage.applyTransformation(newGrayImage, singleChannelGreen)
	case SINGLEBLUE:
		gImage.applyTransformation(newGrayImage, singleChannelBlue)

	}

	return newGrayImage
}

func (cImage *colorImage) Create(FilterName string) *image.RGBA64 {
	newColorImage := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})
	newGrayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})
	finalImage := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})

	rowPerPart := cImage.height / PARTS
	remainderRows := cImage.height % PARTS
	fmt.Println(rowPerPart, remainderRows)
	switch FilterName {
	case REDONLYFILTER:
		cImage.applyTransformation(newColorImage, newGrayImage, redFilter)
	case GREENONLYFILTER:
		cImage.applyTransformation(newColorImage, newGrayImage, greenFilter)
	case BLUEONLYFILTER:
		//fmt.Println(rowPerPart, remainderRows)
		cImage.applyTransformation(newColorImage, newGrayImage, blueFilter)
		//for j := 0; j < PARTS; j++ {
		//	wg.Add(1)
		//	startFromRow := partLimit * j
		//	upToRow := partLimit * (j + 1)
		//	if j == PARTS-1 {
		//		upToRow += difference
		//	}
		//	go cImage.applyTransformation(startFromRow, upToRow, newColorImageAddress, newGrayImageAddress, redFilter)
		//}
	}
	//wg.Wait()
	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}}, newGrayImage, image.Point{0, 0}, draw.Src)
	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}}, newColorImage, image.Point{0, 0}, draw.Over)

	return finalImage
}

func (img *grayImage) applyTransformation(grayImage *image.Gray16, avgFunction func(r, g, b, a uint32) color.Gray16) {
	for i := 0; i <= img.width; i++ {
		for j := 0; j <= img.height; j++ {
			point := img.decodedImage.At(i, j)
			r, g, b, a := point.RGBA()
			grayColor := avgFunction(r, g, b, a)
			grayImage.Set(i, j, grayColor)
		}
	}
}

func (img *colorImage) applyTransformation(colorImage *image.RGBA64, grayImage *image.Gray16, transformationFunction func(r, g, b, a uint32) color.RGBA64) {
	for i := 0; i <= img.width; i++ {
		for j := 0; j <= img.height; j++ {
			point := img.decodedImage.At(i, j)
			r, g, b, a := point.RGBA()

			pixelColor := transformationFunction(r, g, b, a)
			grayColor := basicImproved(r, g, b, a)
			grayImage.SetGray16(i, j, grayColor)
			colorImage.SetRGBA64(i, j, pixelColor)
		}
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

func basic(r, g, b, _ uint32) color.Gray16 {
	avg := float64((r + g + b) / 3)

	return color.Gray16{uint16(math.Ceil(avg))}
}

func basicImproved(r, g, b, _ uint32) color.Gray16 {
	avg := float64(0.3)*float64(r) + float64(0.59)*float64(g) + float64(0.11)*float64(b)

	return color.Gray16{uint16(math.Ceil(avg))}
}

func desaturation(r, g, b, a uint32) color.Gray16 {
	avg := float64(maxOfThree(r, g, b, a)+minOfThree(r, g, b, a)) / 2
	return color.Gray16{uint16(math.Ceil(avg))}
}

func decompositionMax(r, g, b, a uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(maxOfThree(r, g, b, a))))}
}

func decompositionMin(r, g, b, a uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(minOfThree(r, g, b, a))))}
}

func maxOfThree(r, g, b, _ uint32) uint32 {
	return Max(Max(r, g), b)
}

func minOfThree(r, g, b, _ uint32) uint32 {
	return Min(Min(r, g), b)
}

// This is how, I'll do it, Until I figure out a better way
func singleChannelRed(r, _, _, _ uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(r)))}
}

func singleChannelGreen(r, _, _, _ uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(r)))}
}

func singleChannelBlue(r, _, _, _ uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(r)))}
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

//func (cImage *Image) applyTransformation(startFromRow, upToRow int, colorImage **image.RGBA64, grayImage **image.Gray16, transformationFunction func(r, g, b, a uint32) color.RGBA64) {
//	for i := startFromRow; i <= upToRow; i++ {
//		for j := 0; j <= cImage.width; j++ {
//			point := cImage.decodedImage.At(i, j)
//			r, g, b, a := point.RGBA()
//
//			pixelColor := transformationFunction(r, g, b, a)
//			grayAVG := basicImproved(r, g, b)
//			grayColor := color.Gray16{uint16(math.Ceil(grayAVG))}
//			(*grayImage).SetGray16(i, j, grayColor)
//			(*colorImage).SetRGBA64(i, j, pixelColor)
//		}
//	}
//	wg.Done()
//}
