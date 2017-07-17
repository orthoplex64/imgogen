package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math"
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
	// for x := 0; x < width; x++ {
	// 	for y := 0; y < height; y++ {
	// 		outimg.Set(x, y, ColorToFColor(inimg.At(x, y)))
	// 	}
	// }
	colorFG, colorBG := ARGB32ToFColor(0xFFDC5037), ARGB32ToFColor(0xFFFFFFFF)
	distFGBG := colorFG.RGBDist(colorBG)
	RenderByFunc(inimg, outimg, 5, func(inimg image.Image, outimg draw.Image, ix int, iy int, x float64, y float64) *FColor {
		c := ColorToFColor(inimg.At(ix, iy))
		distFG := c.RGBDist(colorFG)
		if math.Mod(((math.Atan2(y, x)+math.Pi)/math.Pi/2+1.0/80.0+1.0/80.0*0.1)*4*10, 1) < 0.2 || math.Hypot(x, y) < 120 {
			c = ARGB32ToFColor(0xFFFFCC00)
			return FColorLerp(c, colorBG, distFG/distFGBG)
		}
		return c
	})
}

// RenderByFunc draws each pixel of outimg according to the colors returned by f.
// aa is the width of the square of per-pixel samples. f is called as
// f(inimg, outimg, ix, iy, x, y) where ix and iy are pixel coordinates.
func RenderByFunc(inimg image.Image, outimg draw.Image, aa int,
	f func(image.Image, draw.Image, int, int, float64, float64) *FColor) {
	bounds := outimg.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	for ix := 0; ix < width; ix++ {
		for iy := 0; iy < height; iy++ {
			c := FColor{}
			for dx := 0; dx < aa; dx++ {
				for dy := 0; dy < aa; dy++ {
					x, y := float64(ix)+float64(dx)/float64(aa), float64(iy)+float64(dy)/float64(aa)
					c.Add(f(inimg, outimg, ix, iy, x, y))
				}
			}
			c.Divide(float64(aa * aa))
			outimg.Set(ix, iy, &c)
		}
	}
}

// Msv - short for "must've"
func Msv(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
