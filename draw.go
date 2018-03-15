package biomorph

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

func DrawTreeCreature(tree *Creature) image.Image {
	dc := gg.NewContext(ImageSize, ImageSize)
	dc.SetColor(color.White)
	dc.DrawRectangle(0, 0, ImageSize, ImageSize)
	dc.Fill()
	dc.SetColor(color.Black)
	dc.SetLineWidth(2.0)
	rand.Seed(0)
	var drawTreeGen func(tree *Creature, gen int, radians float64, p point, branch_size float64, branch_angle float64)
	drawTreeGen = func(tree *Creature, gen int, radians float64, p point, branch_size float64, branch_angle float64) {
		radians += (rand.ExpFloat64() * AngleNoise(tree)) * 0.0
		ba := branch_angle * (1 + LengthNoise(tree)*rand.ExpFloat64())
		bs := branch_size * (1 + LengthNoise(tree)*rand.ExpFloat64()*0.0)
		if gen == 0 {
			return
		}
		new_point := point{p.x - bs*math.Sin(radians), p.y - bs*math.Cos(radians)}
		dc.DrawLine(p.x, p.y, new_point.x, new_point.y)
		dc.Stroke()
		for i := 0; i < NumBranches(tree); i++ {
			drawTreeGen(tree, gen-1, radians-ba/2.0+ba*float64(i)/float64(NumBranches(tree)-1), new_point, branch_size*BranchIncrease(tree), branch_angle*AngleIncrease(tree))
		}
	}
	drawTreeGen(tree, NumGens(tree), 0, point{float64(dc.Width()) / 2, float64(dc.Height()) * 9 / 10}, BranchLength(tree), BranchAngle(tree))
	rand.Seed(time.Now().UnixNano())
	return dc.Image()
}
