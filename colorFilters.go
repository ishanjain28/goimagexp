package imagexp

//	img := filter.img
//	colorImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})
//	grayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})
//
//	finalImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})
//
//	for x := 0; x <= img.width; x++ {
//		for y := 0; y <= img.height; y++ {
//			point := img.srcImage.At(x, y)
//			r, g, b, a := point.RGBA()
//
//			avg := BasicImproved(r, g, b)
//
//			grayColor := color.Gray16{uint16(math.Ceil(avg))}
//			var pixelRGBColor color.RGBA64
//			switch filter.filterType {
//			case "red":
//				pixelRGBColor = redFilter(r, g, b, a)
//			case "green":
//				pixelRGBColor = greenFilter(r, g, b, a)
//			case "blue":
//				pixelRGBColor = blueFilter(r, g, b, a)
//			}
//
//			grayImage.SetGray16(x, y, grayColor)
//			colorImage.Set(x, y, pixelRGBColor)
//		}
//	}
//
//	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}}, grayImage, image.Point{0, 0}, draw.Src)
//	draw.Draw(finalImage, image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}}, colorImage, image.Point{0, 0}, draw.Over)
//
//	return finalImage
//}