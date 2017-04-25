package main

import (
	"image"
	"image/color"
	"math"
	"image/draw"
)

type Filter struct {
	img        *Image
	filterType string
}

//Applies Red Color Only Filter to an Image
func (filter *Filter) ApplyFilter() *image.RGBA {

	img := filter.img
	colorImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})
	grayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})

	finalImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})

	for x := 0; x <= img.width; x++ {
		for y := 0; y <= img.height; y++ {
			point := img.srcImage.At(x, y)
			r, g, b, a := point.RGBA()

			avg := BasicImproved(r, g, b)

			grayColor := color.Gray16{uint16(math.Ceil(avg))}
			var pixelRGBColor color.RGBA64
			switch filter.filterType {
			case "red":
				pixelRGBColor = redFilter(r, g, b, a)
			case "green":
				pixelRGBColor = greenFilter(r, g, b, a)
			case "blue":
				pixelRGBColor = blueFilter(r, g, b, a)
			}

			grayImage.SetGray16(x, y, grayColor)
			colorImage.Set(x, y, pixelRGBColor)
		}
	}

	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}}, grayImage, image.Point{0, 0}, draw.Src)
	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}}, colorImage, image.Point{0, 0}, draw.Over)

	return finalImage
}

func redFilter(r, g, b, a uint32) color.RGBA64 {
	if !(r > b ) || !(r > g ) {
		return color.RGBA64{uint16(255), uint16(255), uint16(255), uint16(0)}
	}
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func greenFilter(r, g, b, a uint32) color.RGBA64 {
	if !(g > b ) || !(g > r ) {
		return color.RGBA64{uint16(255), uint16(255), uint16(255), uint16(0)}
	}
	//avg := BasicImproved(r, g, b)
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

func blueFilter(r, g, b, a uint32) color.RGBA64 {
	if !(b > g ) || !(b > r ) {
		return color.RGBA64{uint16(255), uint16(255), uint16(255), uint16(0)}
	}
	//avg := BasicImproved(r, g, b)
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}
