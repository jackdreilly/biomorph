package biomorph

import (
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
}
