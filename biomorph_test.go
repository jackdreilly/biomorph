package biomorph

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBiomorph(t *testing.T) {
	creature := NewCreature(NewTreeSpecies())
	orig_value := creature.Values[0].Value
	gene := creature.CreatureSpecies.Genes[0]
	assert.Equal(t, 0.5*(gene.Range.Max+gene.Range.Min), orig_value)
	mutated_creature := MutateCreature(creature)
	assert.NotEqual(t, orig_value, mutated_creature.Values[0].Value)
	assert.Equal(t, orig_value, creature.Values[0].Value)
	for i := 0; i < 100; i++ {
		img := DrawTreeCreature(creature)
		f, err := os.Create(fmt.Sprintf("pngs/creature_%2.0d.png", i))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		png.Encode(f, img)
		creature = MutateCreature(creature)
	}
}
