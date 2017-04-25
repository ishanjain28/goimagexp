package main

import (
	"fmt"
	"github.com/ishanjain28/goimagexp"
)

func main() {
	fmt.Println("hello")

	imagexp.TransformImage(imagexp.REDONLYFILTER, "image.png")


}