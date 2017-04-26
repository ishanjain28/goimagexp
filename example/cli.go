package main

import (
	"bufio"
	"fmt"
	"github.com/ishanjain28/goimagexp"
	"image/png"
	"log"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	str := []string{
		imagexp.DESATURATION,
		imagexp.BASIC,
		imagexp.BASICIMPROVED,
		imagexp.SINGLEBLUE,
		imagexp.SINGLEGREEN,
		imagexp.SINGLERED,
		imagexp.REDONLYFILTER,
		imagexp.GREENONLYFILTER,
		imagexp.BLUEONLYFILTER,
		imagexp.DECOMPOSITIONMAX,
		imagexp.DECOMPOSITIONMIN,
	}

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

	for _, v := range str {
		wg.Add(1)
		create(v, imgPath)
	}

	wg.Wait()
}

func create(v, imgPath string) {
	image, err := imagexp.TransformImage(v, imgPath)
	if err != nil {
		log.Fatalf("Error Occurred: %s\n", err)
		//Keep the command prompt open for 5 more seconds
		time.Sleep(5 * time.Second)
	}

	file, _ := os.Create(v + ".png")

	defer file.Close()

	png.Encode(file, image)
	wg.Done()
}


