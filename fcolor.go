package main

import (
	"bytes"
	"fmt"
	"image/color"
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
	res := FColor{}
	res.r = float64(ri) / 0xFFFF
	res.g = float64(gi) / 0xFFFF
	res.b = float64(bi) / 0xFFFF
	res.a = float64(ai) / 0xFFFF
	return &res
}

// ARGB32ToFColor constructs an FColor from a uint32 in format 0xAARRGGBB
func ARGB32ToFColor(argb uint32) *FColor {
	res := FColor{}
	res.a = float64(argb>>(8*3)&0xFF) / 0xFF
	res.r = float64(argb>>(8*2)&0xFF) / 0xFF
	res.g = float64(argb>>(8*1)&0xFF) / 0xFF
	res.b = float64(argb>>(8*0)&0xFF) / 0xFF
	return &res
}

// Cl creates a copy of the color
func (c *FColor) Cl() *FColor {
	return &FColor{c.r, c.g, c.b, c.a}
}

func trimFloat64(n *float64) {
	*n = math.Min(math.Max(*n, 0), 1)
}

// Trim caps each value to [0, 1]
func (c *FColor) Trim() *FColor {
	trimFloat64(&c.r)
	trimFloat64(&c.g)
	trimFloat64(&c.b)
	trimFloat64(&c.a)
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

// RGBDist returns the distance between 2 colors in RGB color space
func RGBDist(a, b *FColor) float64 {
	return math.Sqrt(math.Pow(a.r-b.r, 2) + math.Pow(a.g-b.g, 2) + math.Pow(a.b-b.b, 2))
}

// FColorLerp returns the linear interpolation between a and b.
// n closer to 0 means more of a; n closer to 1 means more of b.
func FColorLerp(a, b FColor, n float64) *FColor {
	trimFloat64(&n)
	return a.Cl().Multiply(1 - n).Add(b.Cl().Multiply(n))
}

// RGBA to implement color.Color
func (c *FColor) RGBA() (r, g, b, a uint32) {
	c = c.Cl().Trim()
	return uint32(c.r * 0xFFFF), uint32(c.g * 0xFFFF),
		uint32(c.b * 0xFFFF), uint32(c.a * 0xFFFF)
}

// String to implement fmt.Stringer
func (c *FColor) String() string {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "FColor{r=%.2f,g=%.2f,b=%.2f,a=%.2f}", c.r, c.g, c.b, c.a)
	return buf.String()
}
