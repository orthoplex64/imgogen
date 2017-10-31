package logo

import (
	"fmt"
	"math"
	"os"

	"github.com/orthoplex64/imgogen/util"
)

// DrawCmurrLogo genereates an SVG of a design I came up with several years ago (as of 2017)
func DrawCmurrLogo(outPath string, scale float64, useViewBox bool, ringColors [4]string) {
	sqrt3, sqrt7, sqrt8 := math.Sqrt(3), math.Sqrt(7), math.Sqrt(8)
	outfile, err := os.Create(outPath)
	util.Msv(err)
	defer outfile.Close()
	_, err = fmt.Fprintf(outfile, `<svg xmlns="http://www.w3.org/2000/svg" `)
	util.Msv(err)
	if useViewBox {
		_, err = fmt.Fprintf(outfile, `viewBox="0 0 %[1]v %[1]v">`, scale)
	} else {
		_, err = fmt.Fprintf(outfile, `width="%[1]v" height="%[1]v">`, scale)
	}
	util.Msv(err)
	// create ring color classes
	_, err = fmt.Fprintf(outfile, "<defs><style>")
	util.Msv(err)
	for i := 0; i < 4; i++ {
		_, err = fmt.Fprintf(outfile, ".ring-%v{fill:#%v;}", i, ringColors[i])
		util.Msv(err)
	}
	_, err = fmt.Fprintf(outfile, "</style></defs><title>CMurr Logo</title>")
	util.Msv(err)
	// rotates point clockwise numTimes times
	rotatePoint := func(point [2]float64, numTimes int) [2]float64 {
		for i := 0; i < numTimes; i++ {
			point[0], point[1] = -point[1], point[0]
		}
		return point
	}
	// converts point from ([-4,4],[-4,4]) to ([0,scale],[0,scale])
	toSvgCoords := func(point [2]float64) [2]float64 {
		point[0] = (point[0] + 4) * scale / 8
		point[1] = (point[1] + 4) * scale / 8
		return point
	}
	type pathPoint struct {
		point     [2]float64
		arcRadius float64
		isConvex  bool // refers to the arc drawn between this point and the previous
	}
	// draws the part specified by path. the first and last point coords must be the same.
	drawRingPart := func(path []pathPoint, shouldDrawMirroredAsWell bool) {
		var paths [][]pathPoint
		if shouldDrawMirroredAsWell {
			mirroredPath := make([]pathPoint, len(path))
			for i, pPoint := range path {
				mirroredPath[i].point[0] = pPoint.point[1]
				mirroredPath[i].point[1] = pPoint.point[0]
				mirroredPath[i].arcRadius = pPoint.arcRadius
				mirroredPath[i].isConvex = !pPoint.isConvex
			}
			paths = [][]pathPoint{path, mirroredPath}
		} else {
			paths = [][]pathPoint{path}
		}
		// iterate paths
		for _, iPath := range paths {
			// iterate rings
			for ring := 0; ring < 4; ring++ {
				_, err = fmt.Fprintf(outfile, `<path class="ring-%d" d="`, ring)
				util.Msv(err)
				for i, pPoint := range iPath {
					point := toSvgCoords(rotatePoint(pPoint.point, ring))
					if i == 0 {
						_, err = fmt.Fprintf(outfile, "M %v %v", point[0], point[1])
						util.Msv(err)
					} else {
						var sweepFlag int
						if pPoint.isConvex {
							sweepFlag = 1
						} else {
							sweepFlag = 0
						}
						_, err = fmt.Fprintf(outfile, " A %[1]v %[1]v 0 0 %v %v %v", pPoint.arcRadius*scale/8, sweepFlag, point[0], point[1])
						util.Msv(err)
					}
				}
				_, err = fmt.Fprint(outfile, `"/>`)
				util.Msv(err)
			}
		}
	}
	// draw the 4 large outer parts
	drawRingPart([]pathPoint{
		{point: [2]float64{-sqrt7 - 1, 0}},
		{point: [2]float64{0, -sqrt7 - 1}, arcRadius: sqrt8, isConvex: true},
		{point: [2]float64{-1, -3}, arcRadius: sqrt8, isConvex: false},
		{point: [2]float64{-3, -1}, arcRadius: 2, isConvex: false},
		{point: [2]float64{-sqrt7 - 1, 0}, arcRadius: sqrt8, isConvex: false},
	}, false)
	// draw the 8 next largest parts
	drawRingPart([]pathPoint{
		{point: [2]float64{1, -3}},
		{point: [2]float64{sqrt3, -sqrt3}, arcRadius: sqrt8, isConvex: true},
		{point: [2]float64{(sqrt7 - 1) / 2, (-sqrt7 - 1) / 2}, arcRadius: sqrt8, isConvex: false},
		{point: [2]float64{0, -sqrt3 - 1}, arcRadius: 2, isConvex: false},
		{point: [2]float64{1, -3}, arcRadius: 2, isConvex: true},
	}, true)
	// draw the 8 parts around the middle
	drawRingPart([]pathPoint{
		{point: [2]float64{(1 + sqrt7) / 2, (1 - sqrt7) / 2}},
		{point: [2]float64{sqrt7 - 1, 0}, arcRadius: sqrt8, isConvex: true},
		{point: [2]float64{1, -1}, arcRadius: sqrt8, isConvex: false},
		{point: [2]float64{(1 + sqrt7) / 2, (1 - sqrt7) / 2}, arcRadius: 2, isConvex: true},
	}, true)
	// draw the 4 center parts
	drawRingPart([]pathPoint{
		{point: [2]float64{1, 1}},
		{point: [2]float64{0, sqrt3 - 1}, arcRadius: 2, isConvex: true},
		{point: [2]float64{sqrt3 - 1, 0}, arcRadius: 2, isConvex: false},
		{point: [2]float64{1, 1}, arcRadius: 2, isConvex: true},
	}, false)
	_, err = fmt.Fprintln(outfile, "</svg>")
	util.Msv(err)
}
