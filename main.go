package main

import (
	"sync"
	"image"
	"fmt"
	"os"
	"log"
)

var wg sync.WaitGroup

type Image struct {
	srcPath  string
	srcImage image.Image
	width    int
	height   int
}

var Options [][]string

var GrayscaleOptions = []string{
	"Grayscale (Basic)",
	"Grayscale (Improved yet Basic)",
	"Grayscale (Desaturated)",
	"Grayscale (Decomposition (MAX))",
	"Grayscale (Decomposition (MIN))",
	"Grayscale (Single Channel (RED))",
	"Grayscale (Single Channel (GREEN))",
	"Grayscale (Single Channel (BLUE))",
}

var FilterOptions = []string{
	"RED Only Filter",
	"GREEN Only Filter",
	"BLUE Only Filter",
}

func init() {
	Options = append(Options, GrayscaleOptions)
	Options = append(Options, FilterOptions)
}

func main() {
	img := Image{}
	var choice int

	//Read Metafile
	//metafile, err := os.Open(".meta.txt")
	//if err != nil {
	//	if err !== os.ErrNotExist  {
	//		log.Fatalln(err)
	//	} else {
	//		img.srcPath = ""
	//	}
	//}
	//metafilereader := bufio.NewReader(metafile)
	//srcPathInMeta, _, err := metafilereader.ReadLine()

	//Ask to enter File path
	fmt.Printf("Enter File Path: ")
	fmt.Scanf("%s\n", &img.srcPath)

	//Open File
	ipImage, err := os.Open(img.srcPath)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer ipImage.Close()

	img.IdentifyDecoder(ipImage)
	//Set Width and height
	img.width, img.height = img.srcImage.Bounds().Max.X, img.srcImage.Bounds().Max.Y

	//Print Resolution
	fmt.Printf("Image Resolution is %dx%d\n", img.width, img.height)

	innerCount := 1
	for _, v := range Options {
		//fmt.Printf("\n", k)

		for _, v1 := range v {
			fmt.Printf("  %d) %s\n", innerCount, v1)
			innerCount++
		}
	}

	fmt.Printf("Choice: ")
	fmt.Scanf("%d\n", &choice)

	if choice > innerCount || choice == 0 {
		log.Fatalln("\nInvalid Choice")
	}
	img.ProcessImage(choice)

	//	Store this choice and reuse it in next turn
	//metaFile, err := os.Create(".meta.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer metaFile.Close()
	//metaFile.WriteString(img.srcPath)
}

func (img *Image) ProcessImage(choice int) {
	var grayImage *image.Gray16
	var colorImage *image.RGBA
	switch choice {
	case 1:
		//Basic
		grayImage = img.CreateGrayImage(0)
		img.Save(Options[0][0], grayImage)
	case 2:
		//Basic Improved
		grayImage = img.CreateGrayImage(1)
		img.Save(Options[0][1], grayImage)
	case 3:
		//Desaturation
		grayImage = img.CreateGrayImage(2)
		img.Save(Options[0][2], grayImage)
	case 4:
		//Decomposition Max
		grayImage = img.CreateGrayImage(3)
		img.Save(Options[0][3], grayImage)
	case 5:
		//Decomposition Min
		grayImage = img.CreateGrayImage(4)
		img.Save(Options[0][4], grayImage)
	case 6:
		//Single Channel Red
		grayImage = img.CreateGrayImage(5)
		img.Save(Options[0][5], grayImage)
	case 7:
		//Single Channel Green
		grayImage = img.CreateGrayImage(6)
		img.Save(Options[0][6], grayImage)
	case 8:
		//Single Channel Blue
		grayImage = img.CreateGrayImage(7)
		img.Save(Options[0][7], grayImage)
	case 9:
		filter := Filter{img, "red"}
		colorImage = filter.ApplyFilter()
		img.Save(Options[1][0], colorImage)
	case 10:
		filter := Filter{img, "green"}
		colorImage = filter.ApplyFilter()
		img.Save(Options[1][1], colorImage)
	case 11:
		filter := Filter{img, "blue"}
		colorImage = filter.ApplyFilter()
		img.Save(Options[1][2], colorImage)
	default:
		log.Fatalf("%d is an invalid choice", choice)
	}
}
