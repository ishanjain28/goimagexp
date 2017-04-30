package imagexp

/*
1. Organize code
2. Add filters.
*/

import (
	"bytes"
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
	"sync"
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

type Gray16Transformation func(r, g, b, a uint32) color.Gray16

type RGBATransformation func(r, g, b, a uint32) color.RGBA64

var PARTS int = 50
var wg sync.WaitGroup

func GrayscaleTransform(transformationFunction Gray16Transformation, ipPath string) (*image.Gray16, error) {
	var finalImage *image.Gray16

	img := &Image{}
	img.path = ipPath

	tempDecodeImage, err := img.Decode()
	if err != nil {
		return nil, err
	}

	img.decodedImage = tempDecodeImage
	img.SetDimension(img.decodedImage.Bounds().Max.X, img.decodedImage.Bounds().Max.Y)

	gImage := grayImage{Image{img.path, img.width, img.height, img.decodedImage}}

	finalImage = gImage.Create(transformationFunction)
	return finalImage, nil
}

func ColorTransform(transformationFunction RGBATransformation, ipPath string) (*image.RGBA64, error) {
	var finalImage *image.RGBA64
	img := &Image{}
	//Set Path
	img.path = ipPath
	//Decode
	//There is some problem with jpeg images, Extremely lit areas are converted to dark patches
	tempDecodeImage, err := img.Decode()
	if err != nil {
		return nil, err
	}

	img.decodedImage = tempDecodeImage
	img.SetDimension(img.decodedImage.Bounds().Max.X, img.decodedImage.Bounds().Max.Y)

	cImage := colorImage{Image{img.path, img.width, img.height, img.decodedImage}}
	finalImage = cImage.Create(transformationFunction)

	return finalImage, nil
}

func (gImage *grayImage) Create(transformationFunction Gray16Transformation) *image.Gray16 {
	newGrayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{gImage.width, gImage.height}})

	rowPerPart := gImage.height / PARTS
	remainderRows := gImage.height % PARTS
	//fmt.Printf("Row Per Part: %d \t Remainder Rows: %d", rowPerPart, remainderRows)

	for j := 0; j < PARTS; j++ {
		wg.Add(1)
		startFromRow := rowPerPart * j
		upToRow := rowPerPart * (j + 1)
		if j == PARTS-1 {
			upToRow += remainderRows
		}

		//fmt.Printf("%d-%d\n", startFromRow, upToRow)
		go gImage.applyTransformation(startFromRow, upToRow, newGrayImage, transformationFunction)
	}

	wg.Wait()
	return newGrayImage
}

func (cImage *colorImage) Create(transformationFunction RGBATransformation) *image.RGBA64 {
	newColorImage := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})
	newGrayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})
	finalImage := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})

	rowPerPart := cImage.height / PARTS
	remainderRows := cImage.height % PARTS
	//fmt.Printf("Row Per Part: %d \t Remainder Rows: %d", rowPerPart, remainderRows)

	for j := 0; j < PARTS; j++ {
		wg.Add(1)
		startFromRow := rowPerPart * j
		upToRow := rowPerPart * (j + 1)
		if j == PARTS-1 {
			upToRow += remainderRows
		}

		//fmt.Printf("%d-%d\n", startFromRow, upToRow)
		go cImage.applyTransformation(startFromRow, upToRow, newColorImage, newGrayImage, transformationFunction)
	}

	wg.Wait()
	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}}, newGrayImage, image.Point{0, 0}, draw.Src)
	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}}, newColorImage, image.Point{0, 0}, draw.Over)

	return finalImage
}

func (img *grayImage) applyTransformation(startFromRow, upToRow int, grayImage *image.Gray16, transformationFunction Gray16Transformation) {
	for i := 0; i <= img.width; i++ {
		for j := startFromRow; j < upToRow; j++ {
			point := img.decodedImage.At(i, j)
			r, g, b, a := point.RGBA()
			grayColor := transformationFunction(r, g, b, a)
			grayImage.Set(i, j, grayColor)
		}
	}
	wg.Done()
}

func (img *colorImage) applyTransformation(startFromRow, uptoRow int, colorImage *image.RGBA64, grayImage *image.Gray16, transformationFunction RGBATransformation) {
	for i := 0; i <= img.width; i++ {
		for j := startFromRow; j < uptoRow; j++ {
			point := img.decodedImage.At(i, j)
			r, g, b, a := point.RGBA()

			pixelColor := transformationFunction(r, g, b, a)
			grayColor := ImprovedGrayscale(r, g, b, a)
			grayImage.SetGray16(i, j, grayColor)
			colorImage.SetRGBA64(i, j, pixelColor)
		}
	}
	wg.Done()
}

func (img *Image) SetDimension(width int, height int) {
	img.width = width
	img.height = height
}

func (img *Image) Decode() (image.Image, error) {
	imageFile, err := os.Open(img.path)
	if err != nil {
		//log.Printf("Error Occurred in opening file: %s", err)
		return nil, err
	}
	defer imageFile.Close()

	fileExtension := path.Ext(img.path)
	var decodedImage image.Image

	switch fileExtension {
	case ".png", ".PNG":
		decodedImage, err = png.Decode(imageFile)
		if err != nil {
			//log.Printf("Error in decoding png: %s", err)
			return nil, err
		}
	case ".jpg", ".jpeg", ".JPG", ".JPEG":
		var jpegBuffer bytes.Buffer

		decodedJPEG, err := jpeg.Decode(imageFile)
		if err != nil {
			//log.Printf("Error in decoding jpeg: %s", err)
			return nil, err
		}

		png.Encode(&jpegBuffer, decodedJPEG)

		decodedImage, err = png.Decode(&jpegBuffer)

		if err != nil {
			//log.Printf("Error in encoding jpeg to png: %s", err)
			return nil, err
		}
	}

	return decodedImage, nil
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

func BasicGrayscale(r, g, b, _ uint32) color.Gray16 {
	avg := float64((r + g + b) / 3)

	return color.Gray16{uint16(math.Ceil(avg))}
}

func ImprovedGrayscale(r, g, b, _ uint32) color.Gray16 {
	avg := float64(0.3)*float64(r) + float64(0.59)*float64(g) + float64(0.11)*float64(b)

	return color.Gray16{uint16(math.Ceil(avg))}
}

func Desaturation(r, g, b, a uint32) color.Gray16 {
	avg := float64(maxOfThree(r, g, b, a)+minOfThree(r, g, b, a)) / 2
	return color.Gray16{uint16(math.Ceil(avg))}
}

func DecompositionMax(r, g, b, a uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(maxOfThree(r, g, b, a))))}
}

func DecompositionMin(r, g, b, a uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(minOfThree(r, g, b, a))))}
}

func maxOfThree(r, g, b, _ uint32) uint32 {
	return Max(Max(r, g), b)
}

func minOfThree(r, g, b, _ uint32) uint32 {
	return Min(Min(r, g), b)
}

// This is how, I'll do it, Until I figure out a better way
func SingleChannelRed(r, _, _, _ uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(r)))}
}

func SingleChannelGreen(r, _, _, _ uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(r)))}
}

func SingleChannelBlue(r, _, _, _ uint32) color.Gray16 {
	return color.Gray16{uint16(math.Ceil(float64(r)))}
}

func RedFilter(r, g, b, a uint32) color.RGBA64 {

	if !(r > b) || !(r > g) {
		return color.RGBA64{uint16(255), uint16(255), uint16(255), uint16(0)}
	}

	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func GreenFilter(r, g, b, a uint32) color.RGBA64 {
	if !(g > r) || !(g > b) {
		return color.RGBA64{uint16(255), uint16(255), uint16(255), uint16(0)}
	}

	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func BlueFilter(r, g, b, a uint32) color.RGBA64 {
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
