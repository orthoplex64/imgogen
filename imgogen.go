package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"

	"github.com/orthoplex64/imgogen/logo"
	"github.com/orthoplex64/imgogen/util"
)

func main() {
	fmt.Println("Start main")
	fmt.Println("End main")
}

func drawOfficialCmurrLogoSVG(outPath string, scale float64, useViewBox bool) {
	logo.DrawCmurrLogoSVG(outPath, scale, useViewBox, [4]string{"A05", "FA0", "0A5", "05A"})
}

func drawOfficialCmurrLogoPNG(outPath string, width, height int) {
	logo.DrawCmurrLogoPNG(outPath, width, height, 16, [4]*util.FColor{
		util.ARGB32ToFColor(0xFFAA0055),
		util.ARGB32ToFColor(0xFFFFAA00),
		util.ARGB32ToFColor(0xFF00AA55),
		util.ARGB32ToFColor(0xFF0055AA),
	}, util.ARGB32ToFColor(0x00000000))
}

func drawDiscourseCard(outPath string) {
	outfile, err := os.Create(outPath)
	util.Msv(err)
	defer outfile.Close()
	outimg := image.NewRGBA(image.Rect(0, 0, 600, 600))
	bounds := outimg.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	util.RenderByFunc(nil, outimg, 1, func(inimg image.Image, outimg draw.Image, ix, iy int, x, y float64) *util.FColor {
		x, y = x-float64(width)/2, y-float64(height)*1.5
		return util.HSVToFColor((math.Atan2(y, x)+math.Pi)/math.Pi/2*8+1.0/2+math.Sin(math.Hypot(x, y)/50), 1, 1)
		//c := SinHueToFColor((math.Atan2(y, x)+math.Pi)/math.Pi/2*8 + 1.0/2 + math.Sin(math.Hypot(x, y)/50))
		//hue, sat, val := c.ToHSV()
		//return HSVToFColor(hue, sat*0+1, val*0.5+0.5)
		//return ColorToFColor(colorful.Hcl(((math.Atan2(y, x)+math.Pi)/math.Pi/2*8+1.0/2+math.Sin(math.Hypot(x, y)/50))*360, 0.5, 0.5))
	})
	util.Msv(png.Encode(outfile, outimg))
}

func drawSkcraftSummerIcon(inPath, outPath string) {
	infile, err := os.Open(inPath)
	util.Msv(err)
	defer infile.Close()
	inimg, err := png.Decode(infile)
	util.Msv(err)
	outfile, err := os.Create(outPath)
	util.Msv(err)
	defer outfile.Close()
	outimg := image.NewRGBA(inimg.Bounds())
	bounds := inimg.Bounds()
	fmt.Println("Bounds:", bounds)
	width, height := bounds.Dx(), bounds.Dy()
	fmt.Println("Dims:", width, height)
	colorFG, colorBG := util.ARGB32ToFColor(0xFFDC5037), util.ARGB32ToFColor(0xFFFFFFFF)
	distFGBG := colorFG.RGBDist(colorBG)
	util.RenderByFunc(inimg, outimg, 5, func(inimg image.Image, outimg draw.Image, ix, iy int, x, y float64) *util.FColor {
		c := util.ColorToFColor(inimg.At(ix, iy))
		distFG := c.RGBDist(colorFG)
		if math.Mod(((math.Atan2(y, x)+math.Pi)/math.Pi/2+1.0/80.0+1.0/80.0*0.1)*4*10, 1) < 0.2 || math.Hypot(x, y) < 120 {
			c = util.ARGB32ToFColor(0xFFFFCC00)
			return util.FColorLerp(c, colorBG, distFG/distFGBG)
		}
		return c
	})
	util.Msv(png.Encode(outfile, outimg))
}
