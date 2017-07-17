package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func main() {
	fmt.Println("Start main")
	infile, err := os.Open("C:/Users/Christopher/Pictures/skcraftIcon.png")
	Msv(err)
	defer infile.Close()
	inimg, err := png.Decode(infile)
	Msv(err)
	outfile, err := os.Create("C:/Users/Christopher/Pictures/skcraftIconSummer.png")
	Msv(err)
	defer outfile.Close()
	outimg := image.NewRGBA(inimg.Bounds())
	drawSkcraftSummerIcon(inimg, outimg)
	Msv(png.Encode(outfile, outimg))
	fmt.Println("End main")
}

func drawSkcraftSummerIcon(inimg image.Image, outimg draw.Image) {
	bounds := inimg.Bounds()
	fmt.Println("Bounds:", bounds)
	width, height := bounds.Dx(), bounds.Dy()
	fmt.Println("Dims:", width, height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			outimg.Set(x, y, ColorToFColor(inimg.At(x, y)))
		}
	}
}

// Msv - short for "must've"
func Msv(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
