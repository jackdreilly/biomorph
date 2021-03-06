package biomorph

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

const (
	ImageSize = 150
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
	if rand.Float64() > 0.3 {
		return &GeneValue{g.Gene, g.Value}
	}
	new_value := g.Value + (0.5-rand.Float64())*0.3*(g.Gene.Range.Max-g.Gene.Range.Min)
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
	genes = append(genes, &Gene{GeneRange{float64(ImageSize) * 0.1, float64(ImageSize) * 0.4}, "branch_length"})
	genes = append(genes, &Gene{GeneRange{2, 5}, "num_gens"})
	genes = append(genes, &Gene{GeneRange{0.1, 5}, "branch_angle"})
	genes = append(genes, &Gene{GeneRange{0.1, 2}, "branch_increase"})
	genes = append(genes, &Gene{GeneRange{0.1, 2}, "angle_increase"})
	genes = append(genes, &Gene{GeneRange{2, 9}, "num_branches"})
	genes = append(genes, &Gene{GeneRange{-0.1, 0.1}, "angle_noise"})
	genes = append(genes, &Gene{GeneRange{-0.1, 0.1}, "length_noise"})
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

func AngleNoise(tree *Creature) float64 {
	return tree.GetValue("angle_noise")
}

func LengthNoise(tree *Creature) float64 {
	return tree.GetValue("length_noise")
}
