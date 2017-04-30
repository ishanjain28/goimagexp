package main

import (
	"github.com/ishanjain28/goimagexp"
	"image/png"
	"log"
	"os"
)

func main() {

	//imgPath := ""
	//fmt.Printf("Enter Path to Image: ")

	//scanner := bufio.NewScanner(os.Stdin)
	//for scanner.Scan() {
	//	s := scanner.Text()
	//
	//	imgPath = s
	//	break
	//
	//}
	//if err := scanner.Err(); err != nil {
	//	os.Exit(1)
	//}
	test, err := imagexp.ColorTransform(imagexp.RedFilter, "image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	file, _ := os.Create("test.png")

	defer file.Close()

	png.Encode(file, test)
}
