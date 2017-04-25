package imagexp

//func (img *ImageXP) CreateGrayImage(index int) *image.Gray16 {
//	grayImage := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{img.width, img.height}})
//	for x := 0; x <= img.width; x++ {
//		for y := 0; y <= img.height; y++ {
//			point := img.srcImage.At(x, y)
//			r, g, b, _ := point.RGBA()
//			var avg float64
//			 This switch uses different algorithms to create values that'll be used in image
//switch index {
//case 0:
//	avg = basic(r, g, b)
//case 1:
//	avg = BasicImproved(r, g, b)
//case 2:
//	avg = desaturation(r, g, b)
//case 3:
//	avg = float64(maxOfThree(r, g, b))
//case 4:
//	avg = float64(minOfThree(r, g, b))
//case 5:
//	avg = float64(r)
//case 6:
//	avg = float64(g)
//case 7:
//	avg = float64(b)
//}
//
//grayColor := color.Gray16{uint16(math.Ceil(avg))}
//Set the color of pixel
//grayImage.Set(x, y, grayColor)
//}
//}
//return grayImage
//
//}
