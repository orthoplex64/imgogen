package util

import (
	"image"
	"image/draw"
	"log"
)

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

// Min returns the minimum of a and b
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Msv - short for "must've"
func Msv(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
