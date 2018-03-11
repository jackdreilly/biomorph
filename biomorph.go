package biomorph

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

const (
	image_size = 150
)

type GeneRange struct {
	Min float64
	Max float64
}

type Gene struct {
	Range GeneRange
	Name  string
}

type GeneValue struct {
	Gene  *Gene
	Value float64
}

func (g *GeneValue) Mutate() *GeneValue {
	new_value := g.Value + (0.5-rand.Float64())*0.2*(g.Gene.Range.Max-g.Gene.Range.Min)
	if new_value > g.Gene.Range.Max {
		new_value = g.Gene.Range.Max
	}
	if new_value < g.Gene.Range.Min {
		new_value = g.Gene.Range.Min
	}
	return &GeneValue{g.Gene, new_value}
}

type Species struct {
	Genes []*Gene
}

func (c *Creature) SetValuesFromMap(m map[string]float64) {
	for k, v := range m {
		c.GetGeneValue(k).Value = v
	}
}

type Creature struct {
	CreatureSpecies *Species
	Values          []*GeneValue
}

func (c *Creature) ValuesMap() map[string]float64 {
	a := map[string]float64{}
	for _, v := range c.Values {
		a[v.Gene.Name] = v.Value
	}
	return a
}

func TreeGenes() (genes []*Gene) {
	genes = append(genes, &Gene{GeneRange{float64(image_size) * 0.1, float64(image_size) * 0.1}, "branch_length"})
	genes = append(genes, &Gene{GeneRange{2, 7}, "num_gens"})
	genes = append(genes, &Gene{GeneRange{0.1, 1.5}, "branch_angle"})
	genes = append(genes, &Gene{GeneRange{0.6, 1.5}, "branch_increase"})
	genes = append(genes, &Gene{GeneRange{0.6, 1.5}, "angle_increase"})
	genes = append(genes, &Gene{GeneRange{1, 4}, "num_branches"})
	return
}

func NewSpecies(genes []*Gene) *Species {
	return &Species{genes}
}

func NewTreeSpecies() *Species {
	return NewSpecies(TreeGenes())
}

func NewCreature(species *Species) *Creature {
	c := &Creature{species, make([]*GeneValue, len(species.Genes))}
	for i, gene := range species.Genes {
		c.Values[i] = &GeneValue{gene, 0.5 * (gene.Range.Max + gene.Range.Min)}
	}
	return c
}

func MutateCreature(creature *Creature) *Creature {
	new_creature := NewCreature(creature.CreatureSpecies)
	for i, v := range creature.Values {
		new_creature.Values[i] = v.Mutate()
	}
	return new_creature
}

type point struct {
	x float64
	y float64
}

type line struct {
	a point
	b point
}

func DrawLine(img *image.Gray16, l *line) {
	x := float64(l.a.x)
	y := float64(l.a.y)
	distance := math.Sqrt(math.Pow(float64(l.a.x-l.b.x), 2) + math.Pow(float64(l.a.y-l.b.y), 2))
	slope := float64((l.b.y - l.a.y)) / float64((l.b.x - l.a.x))
	dir_y := 1.0
	if l.b.y < l.a.y {
		dir_y = -1.0
	}
	dir_x := dir_y
	if slope < 0 {
		dir_x = -dir_y
	}
	dx := dir_x
	dy := dir_y
	if slope < 0 {
		slope *= -1
	}
	if slope > 0 {
		dx /= slope
	} else {
		dy *= slope
	}
	for i := 0; i < int(distance); i++ {
		img.Set(int(x), int(y), color.Black)
		x += dx
		y += dy
	}
}

func (c *Creature) GetGeneValue(name string) *GeneValue {
	for _, gene := range c.Values {
		if gene.Gene.Name == name {
			return gene
		}
	}
	return nil
}

func (c *Creature) GetValue(name string) float64 {
	return c.GetGeneValue(name).Value
}

func NumGens(tree *Creature) int {
	return int(tree.GetValue("num_gens"))
}

func NumBranches(tree *Creature) int {
	return int(tree.GetValue("num_branches"))
}

func BranchAngle(tree *Creature) float64 {
	return tree.GetValue("branch_angle")
}

func BranchLength(tree *Creature) float64 {
	return tree.GetValue("branch_length")
}

func BranchIncrease(tree *Creature) float64 {
	return tree.GetValue("branch_increase")
}

func AngleIncrease(tree *Creature) float64 {
	return tree.GetValue("angle_increase")
}

func drawTreeGen(img *image.Gray16, tree *Creature, gen int, radians float64, p point, branch_size float64, branch_angle float64) {
	if gen == 0 {
		return
	}
	new_point := point{p.x - branch_size*math.Sin(radians), p.y - branch_size*math.Cos(radians)}
	DrawLine(img, &line{p, new_point})
	for i := 0; i < NumBranches(tree)+1; i++ {
		drawTreeGen(img, tree, gen-1, radians-branch_angle/2.0+branch_angle*float64(i)/float64(NumBranches(tree)), new_point, branch_size*BranchIncrease(tree), branch_angle*AngleIncrease(tree))
	}
}

func DrawTreeCreature(tree *Creature) *image.Gray16 {
	img := image.NewGray16(image.Rect(0, 0, image_size, image_size))
	for i := 0; i < image_size; i++ {
		for j := 0; j < image_size; j++ {
			img.Set(i, j, color.White)
		}
	}
	drawTreeGen(img, tree, NumGens(tree), 0, point{float64(img.Rect.Size().X) / 2, float64(img.Rect.Size().Y) * 9 / 10}, BranchLength(tree), BranchAngle(tree))
	return img
}
