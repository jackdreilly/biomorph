package biomorph

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

func DrawTreeCreature(tree *Creature) image.Image {
	dc := gg.NewContext(ImageSize, ImageSize)
	dc.SetColor(color.White)
	dc.DrawRectangle(0, 0, ImageSize, ImageSize)
	dc.Fill()
	dc.SetColor(color.Black)
	dc.SetLineWidth(3.0)
	var drawTreeGen func(tree *Creature, gen int, radians float64, p point, branch_size float64, branch_angle float64)
	drawTreeGen = func(tree *Creature, gen int, radians float64, p point, branch_size float64, branch_angle float64) {
		if gen == 0 {
			return
		}
		new_point := point{p.x - branch_size*math.Sin(radians), p.y - branch_size*math.Cos(radians)}
		dc.DrawLine(p.x, p.y, new_point.x, new_point.y)
		dc.Stroke()
		for i := 0; i < NumBranches(tree); i++ {
			drawTreeGen(tree, gen-1, radians-branch_angle/2.0+branch_angle*float64(i)/float64(NumBranches(tree)-1), new_point, branch_size*BranchIncrease(tree), branch_angle*AngleIncrease(tree))
		}
	}
	drawTreeGen(tree, NumGens(tree), 0, point{float64(dc.Width()) / 2, float64(dc.Height()) * 9 / 10}, BranchLength(tree), BranchAngle(tree))
	return dc.Image()
}
