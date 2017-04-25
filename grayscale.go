package main

import (
	"image"
	"image/color"
	"math"
)

//func (img *Image) CreateImages(index int, name string) {
//	var grayImage *image.Gray16
//	var RGBAImage *image.RGBA
//	switch index {
//	case 0:
//		grayImage = img.CreateGrayImage(0)
//		img.Save(name, grayImage)
//	case 1:
//		grayImage = img.CreateGrayImage(1)
//		img.Save(name, grayImage)
//	case 2:
//		grayImage = img.CreateGrayImage(2)
//		img.Save(name, grayImage)
//	case 3:
//		grayImage = img.CreateGrayImage(3)
//		img.Save(name, grayImage)
//	case 4:
//		grayImage = img.CreateGrayImage(4)
//		img.Save(name, grayImage)
//	case 5:
//		grayImage = img.CreateGrayImage(5)
//		img.Save(name, grayImage)
//	case 6:
//		grayImage = img.CreateGrayImage(6)
//		img.Save(name, grayImage)
//	case 7:
//		grayImage = img.CreateGrayImage(7)
//		img.Save(name, grayImage)
//	case 8:
//		RGBAImage = img.CreateRGBAImage()
//
//		redOnlyFile, _ := os.Create("test.png")
//		defer redOnlyFile.Close()
//		png.Encode(redOnlyFile, RGBAImage)
//
//	}
//
//	wg.Done()
//}
//
////Creates a GrayScaled Image
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
				avg = basic(r, g, b)
			case 1:
				avg = basicImproved(r, g, b)
			case 2:
				avg = desaturation(r, g, b)
			case 3:
				avg = float64(maxOfThree(r, g, b))
			case 4:
				avg = float64(minOfThree(r, g, b))
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
