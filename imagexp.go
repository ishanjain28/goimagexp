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

var PARTS int = 200
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

func ColorFilter(transformationFunction RGBATransformation, ipPath string) (*image.RGBA64, error) {
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

func Blur(blurLevel int, ipPath string) (*image.RGBA64, error) {

	img := &Image{}
	img.path = ipPath

	tempDecode, err := img.Decode()
	if err != nil {
		return nil, err
	}
	img.decodedImage = tempDecode
	img.SetDimension(img.decodedImage.Bounds().Max.X, img.decodedImage.Bounds().Max.Y)

	var finalImage *image.RGBA64 = image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})

	//for i := blurLevel; i <= img.width-blurLevel-blurLevel; i += (blurLevel*2 + 1) {
	//	for j := blurLevel; j < img.height-blurLevel-blurLevel; j += (blurLevel*2 + 1) {
	//		r00, g00, b00, a00 := img.decodedImage.At(i-1, j-1).RGBA()
	//		r10, g10, b10, a10 := img.decodedImage.At(i, j-1).RGBA()
	//		r20, g20, b20, a20 := img.decodedImage.At(i+1, j-1).RGBA()
	//
	//		r01, g01, b01, a01 := img.decodedImage.At(i-1, j).RGBA()
	//		r11, g11, b11, a11 := img.decodedImage.At(i, j).RGBA()
	//		r21, g21, b21, a21 := img.decodedImage.At(i+1, j).RGBA()
	//
	//		r02, g02, b02, a02 := img.decodedImage.At(i-1, j+1).RGBA()
	//		r12, g12, b12, a12 := img.decodedImage.At(i, j+1).RGBA()
	//		r22, g22, b22, a22 := img.decodedImage.At(i+1, j+1).RGBA()
	//
	//		rAvg := math.Ceil(float64(r00+r10+r20+r01+r11+r21+r02+r12+r22) / 9)
	//		gAvg := math.Ceil(float64(g00+g10+g20+g01+g11+g21+g02+g12+g22) / 9)
	//		bAvg := math.Ceil(float64(b00+b10+b20+b01+b11+b21+b02+b12+b22) / 9)
	//		aAvg := math.Ceil(float64(a00+a10+a20+a01+a11+a21+a02+a12+a22) / 9)
	//
	//		rgbaColor := color.RGBA64{uint16(rAvg), uint16(gAvg), uint16(bAvg), uint16(aAvg)}
	//
	//		fmt.Println(i, j)
	//		fmt.Println(rgbaColor)
	//		fmt.Println(img.decodedImage.At(i-1, j-1).RGBA())
	//		finalImage.Set(i-1, j-1, rgbaColor)
	//		finalImage.Set(i, j-1, rgbaColor)
	//		finalImage.Set(i+1, j-1, rgbaColor)
	//
	//		finalImage.SetRGBA64(i-1, j, rgbaColor)
	//		finalImage.SetRGBA64(i, j, rgbaColor)
	//		finalImage.SetRGBA64(i+1, j, rgbaColor)
	//
	//		finalImage.SetRGBA64(i-1, j+1, rgbaColor)
	//		finalImage.SetRGBA64(i, j+1, rgbaColor)
	//		finalImage.SetRGBA64(i+1, j+1, rgbaColor)
	//
	//	}
	//}

	return finalImage, nil
}

func (gImage *grayImage) Create(transformationFunction Gray16Transformation) *image.Gray16 {
	newGrayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{gImage.width, gImage.height}})

	rowPerPart := gImage.height / PARTS
	remainderRows := gImage.height % PARTS

	for j := 0; j < PARTS; j++ {
		wg.Add(1)
		startFromRow := rowPerPart * j
		upToRow := rowPerPart * (j + 1)
		if j == PARTS-1 {
			upToRow += remainderRows
		}

		go gImage.applyTransformation(startFromRow, upToRow, newGrayImage, transformationFunction)
	}

	wg.Wait()
	//Wait Until All routines are done
	return newGrayImage
}

func (cImage *colorImage) Create(transformationFunction RGBATransformation) *image.RGBA64 {
	newColorImage := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})
	newGrayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})
	finalImage := image.NewRGBA64(image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}})

	rowPerPart := cImage.height / PARTS
	remainderRows := cImage.height % PARTS

	// Divide workload to multiple go routines, Each one handling some columns of pixels in an image
	for j := 0; j < PARTS; j++ {
		wg.Add(1)
		startFromRow := rowPerPart * j
		upToRow := rowPerPart * (j + 1)

		if j == PARTS-1 {
			upToRow += remainderRows
		}

		go cImage.applyTransformation(startFromRow, upToRow, newColorImage, newGrayImage, transformationFunction)
	}

	wg.Wait()
	//Wait Until All Routines are done.
	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}}, newGrayImage, image.Point{0, 0}, draw.Src)
	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{cImage.width, cImage.height}}, newColorImage, image.Point{0, 0}, draw.Over)
	return finalImage
}

func (img *grayImage) applyTransformation(startFromRow, upToRow int, grayImage *image.Gray16, transformationFunction Gray16Transformation) {
	//fmt.Printf("%d-%d\n", startFromRow, upToRow)
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
	//fmt.Printf("%d-%d\n", startFromRow, uptoRow)
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
		return nil, err
	}
	defer imageFile.Close()

	fileExtension := path.Ext(img.path)
	var decodedImage image.Image

	switch fileExtension {
	case ".png", ".PNG":
		decodedImage, err = png.Decode(imageFile)
		if err != nil {
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
			return nil, err
		}
	}
	return decodedImage, nil
}

// Use file Extension to idenitfy which decoder to use.
func (img *Image) Save(SaveDir string, finalImage image.Image, shouldCreateDir bool) {
	destFileName := path.Base(img.path)
	destFileName = strings.Replace(destFileName, path.Ext(img.path), "", -1)
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
