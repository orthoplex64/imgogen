package util

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
)

// FColor stores red, green, blue, and alpha values as float64's ranging from
// 0 to 1 inclusive. Mutable.
type FColor struct {
	r, g, b, a float64
}

// ColorToFColor constructs an FColor from an existing color.Color
func ColorToFColor(c color.Color) *FColor {
	ri, gi, bi, ai := c.RGBA()
	fc := FColor{}
	fc.r = float64(ri) / 0xFFFF
	fc.g = float64(gi) / 0xFFFF
	fc.b = float64(bi) / 0xFFFF
	fc.a = float64(ai) / 0xFFFF
	return &fc
}

// ARGB32ToFColor constructs an FColor from a uint32 in format 0xAARRGGBB
func ARGB32ToFColor(argb uint32) *FColor {
	c := FColor{}
	c.a = float64(argb>>(8*3)&0xFF) / 0xFF
	c.r = float64(argb>>(8*2)&0xFF) / 0xFF
	c.g = float64(argb>>(8*1)&0xFF) / 0xFF
	c.b = float64(argb>>(8*0)&0xFF) / 0xFF
	return &c
}

// HSVToFColor constructs an FColor from hue, saturation, and value
func HSVToFColor(hue float64, sat float64, val float64) *FColor {
	c := FColor{a: 1}
	if sat == 0 {
		c.r, c.g, c.b = val, val, val
		return &c
	}
	hue6 := (hue - math.Floor(hue)) * 6
	hue6r := hue6 - math.Floor(hue6)
	f1 := val * (1 - sat)
	f2 := val * (1 - sat*hue6r)
	f3 := val * (1 - sat*(1-hue6r))
	switch int(hue6) {
	case 0:
		c.r, c.g, c.b = val, f3, f1
	case 1:
		c.r, c.g, c.b = f2, val, f1
	case 2:
		c.r, c.g, c.b = f1, val, f3
	case 3:
		c.r, c.g, c.b = f1, f2, val
	case 4:
		c.r, c.g, c.b = f3, f1, val
	case 5:
		c.r, c.g, c.b = val, f1, f2
	default:
		log.Fatal("HSVToFColor: int(hue6) is", int(hue6))
	}
	return &c
}

// SinHueToFColor returns a color with sine wave-based hue.
func SinHueToFColor(h float64) *FColor {
	h *= math.Pi * 2
	c := FColor{a: 1}
	c.r = (math.Cos(h-math.Pi*2*0/3) + 1) / 2
	c.g = (math.Cos(h-math.Pi*2*1/3) + 1) / 2
	c.b = (math.Cos(h-math.Pi*2*2/3) + 1) / 2
	c.Clamp()
	return &c
}

// Cl creates a copy of the color
func (c *FColor) Cl() *FColor {
	return &FColor{c.r, c.g, c.b, c.a}
}

func clampFloat64(n *float64) {
	*n = math.Min(math.Max(*n, 0), 1)
}

// Clamp clamps each value to [0, 1]
func (c *FColor) Clamp() *FColor {
	clampFloat64(&c.r)
	clampFloat64(&c.g)
	clampFloat64(&c.b)
	clampFloat64(&c.a)
	return c
}

// Add adds the components of addend to the respective components of c
func (c *FColor) Add(addend *FColor) *FColor {
	c.r += addend.r
	c.g += addend.g
	c.b += addend.b
	c.a += addend.a
	return c
}

// Multiply multiplies the components of c by factor
func (c *FColor) Multiply(factor float64) *FColor {
	c.r *= factor
	c.g *= factor
	c.b *= factor
	c.a *= factor
	return c
}

// Divide divides the components of c by divisor
func (c *FColor) Divide(divisor float64) *FColor {
	c.r /= divisor
	c.g /= divisor
	c.b /= divisor
	c.a /= divisor
	return c
}

// ToHSV returns the hue, saturation, and value of c in HSV color space
func (c *FColor) ToHSV() (hue, sat, val float64) {
	i := math.Max(math.Max(c.r, c.g), c.b)
	j := math.Min(math.Min(c.r, c.g), c.b)
	val = i
	if i > 0 {
		sat = (i - j) / i
	} else {
		sat = 0
	}
	if sat == 0 {
		hue = 0
		return
	}
	fr, fg, fb := (i-c.r)/(i-j), (i-c.g)/(i-j), (i-c.b)/(i-j)
	switch i {
	case c.r:
		hue = fb - fg
	case c.g:
		hue = 2 + fr - fb
	default:
		hue = 4 + fg - fr
	}
	hue /= 6
	if hue < 0 {
		hue++
	}
	return
}

// RGBDist returns the distance between c and other in RGB color space
func (c *FColor) RGBDist(other *FColor) float64 {
	return math.Sqrt(math.Pow(c.r-other.r, 2) + math.Pow(c.g-other.g, 2) + math.Pow(c.b-other.b, 2))
}

// FColorLerp returns the linear interpolation between a and b.
// n closer to 0 means more of a; n closer to 1 means more of b.
func FColorLerp(a, b *FColor, n float64) *FColor {
	clampFloat64(&n)
	return a.Cl().Multiply(1 - n).Add(b.Cl().Multiply(n))
}

// ARGB32 returns c in a uint32 0xAARRGGBB
func (c *FColor) ARGB32() uint32 {
	c = c.Cl().Clamp()
	return uint32(c.a*0xFF)<<(8*3) | uint32(c.r*0xFF)<<(8*2) |
		uint32(c.g*0xFF)<<(8*1) | uint32(c.b*0xFF)<<(8*0)
}

// RGBA to implement color.Color
func (c *FColor) RGBA() (r, g, b, a uint32) {
	c = c.Cl().Clamp()
	return uint32(c.r * 0xFFFF), uint32(c.g * 0xFFFF),
		uint32(c.b * 0xFFFF), uint32(c.a * 0xFFFF)
}

// String to implement fmt.Stringer
func (c *FColor) String() string {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "FColor{r=%.2f,g=%.2f,b=%.2f,a=%.2f}", c.r, c.g, c.b, c.a)
	return buf.String()
}
