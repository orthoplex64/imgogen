package logo

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/orthoplex64/imgogen/util"
)

// DrawCmurrLogoPNG genereates a PNG file of a design I came up with several years ago (as of 2017)
func DrawCmurrLogoPNG(outPath string, width, height, aa int, ringColors [4]*util.FColor, bgColor *util.FColor) {
	outfile, err := os.Create(outPath)
	util.Msv(err)
	defer outfile.Close()
	outimg := image.NewRGBA(image.Rect(0, 0, width, height))
	allRings := [][2]float64{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}}
	isInRing := func(x, y float64, ring [2]float64) bool {
		distSquared := (x-ring[0])*(x-ring[0]) + (y-ring[1])*(y-ring[1])
		return distSquared >= 2*2 && distSquared <= 8 // 8 is 2*sqrt2 * 2*sqrt2
	}
	minDim := float64(util.Min(width, height))
	util.RenderByFunc(nil, outimg, aa, func(inimg image.Image, outimg draw.Image, ix, iy int, x, y float64) *util.FColor {
		x, y = (x-float64(width)/2)*8/minDim, (y-float64(height)/2)*8/minDim
		inRing := -1
		for i, ring := range allRings {
			if isInRing(x, y, ring) {
				if inRing != -1 {
					return bgColor
				}
				inRing = i
			}
		}
		if inRing == -1 {
			return bgColor
		}
		return ringColors[inRing]
	})
	util.Msv(png.Encode(outfile, outimg))
}
