package main

import (
	"os"
	"log"
	"image/png"
	"fmt"
	"image"
	"image/color"
	"sync"
	"path"
	"image/jpeg"
	"math"
	"image/draw"
)

var wg sync.WaitGroup

var GrayscalingAlgos = []string{
	"basic",
	"basic.improved.for.human.eye",
	"desaturated",
	"decomposition.max",
	"decomposition.min",
	"single.channel.red",
	"single.channel.green",
	"single.channel.blue",
	"inverse.channel.red",
}

type Image struct {
	srcPath  string
	srcImage image.Image
	width    int
	height   int
}

func main() {

	ipImage := Image{}
	fmt.Printf("Enter File Path/Name: ")
	fmt.Scanf("%s", &ipImage.srcPath)

	// Open File
	img, err := os.Open(ipImage.srcPath)
	defer img.Close()
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Use file Extension to idenitfy which decoder to use.
	fileExtension := path.Ext(ipImage.srcPath)
	switch fileExtension {
	case ".png":
		ipImage.srcImage, err = png.Decode(img)
	case ".jpg":
		ipImage.srcImage, err = jpeg.Decode(img)
	case ".jpeg":
		ipImage.srcImage, err = jpeg.Decode(img)
	}

	//Update width and height in ipImage(inputImage)
	ipImage.width, ipImage.height = ipImage.srcImage.Bounds().Max.X, ipImage.srcImage.Bounds().Max.Y
	fmt.Printf("Image Resolution is %dx%d\n", ipImage.width, ipImage.height)

	// Start Creating Images
	for index, name := range GrayscalingAlgos {
		wg.Add(1)
		go ipImage.CreateImages(index, name)
	}
	//Wait until goroutines finish
	wg.Wait()
}

func (img *Image) CreateImages(index int, name string) {
	var grayImage *image.Gray16
	var RGBAImage *image.RGBA
	switch index {
	case 0:
		grayImage = img.CreateGrayImage(0)
		img.Save(name, grayImage)
	case 1:
		grayImage = img.CreateGrayImage(1)
		img.Save(name, grayImage)
	case 2:
		grayImage = img.CreateGrayImage(2)
		img.Save(name, grayImage)
	case 3:
		grayImage = img.CreateGrayImage(3)
		img.Save(name, grayImage)
	case 4:
		grayImage = img.CreateGrayImage(4)
		img.Save(name, grayImage)
	case 5:
		grayImage = img.CreateGrayImage(5)
		img.Save(name, grayImage)
	case 6:
		grayImage = img.CreateGrayImage(6)
		img.Save(name, grayImage)
	case 7:
		grayImage = img.CreateGrayImage(7)
		img.Save(name, grayImage)
	case 8:
		RGBAImage = img.CreateRGBAImage()

		redOnlyFile, _ := os.Create("test.png")
		defer redOnlyFile.Close()
		png.Encode(redOnlyFile, RGBAImage)

	}

	wg.Done()
}

//Creates a GrayScaled Image
func (img *Image) CreateGrayImage(index int) *image.Gray16 {

	grayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})
	for x := 0; x <= img.width; x++ {
		for y := 0; y <= img.height; y++ {
			point := img.srcImage.At(x, y)
			r, g, b, _ := point.RGBA()
			var avg float64
			// This switch uses different algorithms to create values that'll be used in image
			switch index {
			case 0:
				avg = Basic(r, g, b)
			case 1:
				avg = BasicImproved(r, g, b)
			case 2:
				avg = Desaturation(r, g, b)
			case 3:
				avg = float64(MaxOfThree(r, g, b))
			case 4:
				avg = float64(MinOfThree(r, g, b))
			case 5:
				avg = float64(r)
			case 6:
				avg = float64(g)
			case 7:
				avg = float64(b)
			}

			grayColor := color.Gray16{uint16(math.Ceil(avg))}
			//Set the color of pixel
			grayImage.Set(x, y, grayColor)
		}
	}
	return grayImage
}

//Creates a RGBA Image
//This is also used to create an image,
//in which only red color might be visible and all other colors will be replaced with gray
func (img *Image) CreateRGBAImage() *image.RGBA {
	colorImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})
	grayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})
	finalImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})

	for x := 0; x <= img.width; x++ {
		for y := 0; y <= img.height; y++ {

			point := img.srcImage.At(x, y)
			r, g, b, a := point.RGBA()
			avg := BasicImproved(r, g, b)

			pixelRGBColor := KeepRedOnly(r, g, b, a)
			grayColor := color.Gray16{uint16(math.Ceil(avg))}
			grayImage.SetGray16(x, y, grayColor)
			colorImage.SetRGBA(x, y, pixelRGBColor)
		}
	}
	fmt.Println(colorImage.Opaque()) //-> False
	fmt.Println(grayImage.Opaque()) //-> True

	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}}, grayImage, image.Point{0, 0}, draw.Src)
	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}}, colorImage, image.Point{0, 0}, draw.Over)

	return finalImage
}

func (img *Image) Save(saveFileName string, grayImage image.Image) {

	srcFilename := path.Base(img.srcPath)
	//fileExtension := path.Ext(img.srcPath)

	os.Mkdir("grayscaled", 0777)

	os.Mkdir("grayscaled/"+srcFilename, 0777)

	outfile, err := os.Create("./grayscaled/" + srcFilename + "/" + saveFileName + ".png")
	if err != nil {
		log.Fatalln(err)
	}
	defer outfile.Close()

	png.Encode(outfile, grayImage)
}

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

//Incomplete. Doesn't works Correctly.
//I want to create a color filter which allows only one color R,G or B
//And other colors are grey
//TODO: Find better way to get pixels where red color is significantly more visible
//Make other pixels transparent
func KeepRedOnly(r, g, b, a uint32) color.RGBA {
	if !(r > b ) || !(r > g ) {
		return color.RGBA{uint8(255), uint8(255), uint8(255), uint8(0)}
	}
	//avg := BasicImproved(r, g, b)
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
