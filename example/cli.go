package main

import (
	"bufio"
	"fmt"
	"github.com/ishanjain28/goimagexp"
	"image/png"
	"log"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {

	imgPath := ""
	fmt.Printf("Enter Path to Image: ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()

		imgPath = s
		break

	}
	if err := scanner.Err(); err != nil {
		os.Exit(1)
	}
	test, err := imagexp.ColorTransform(imagexp.BlueFilter, imgPath)
	if err != nil {
		log.Fatal(err)
	}
	file, _ := os.Create("test.png")

	defer file.Close()

	png.Encode(file, test)

}
