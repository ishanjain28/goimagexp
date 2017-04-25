package imagexp

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

//func (img *ImageXP) Set(path string) {
//	img.path = path
//}
//
//func (img *ImageXP) Get() string {
//	return img.path
//}
