package biome

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/t14raptor/pm-gen/populate"
)

type River struct{}

func (r River) Populators() []populate.Populator {
	return []populate.Populator{populate.TallGrass{Amount: 5}}
}

func (r River) ID() uint8 {
	return IDRiver
}

func (r River) Elevation() (min, max int) {
	return 58, 62
}

func (r River) GroundCover() []world.Block {
	return []world.Block{
		block.Dirt{},
		block.Dirt{},
		block.Dirt{},
		block.Dirt{},
		block.Dirt{},
	}
}

func (r River) Temperature() float64 {
	return 0.5
}

func (r River) Rainfall() float64 {
	return 0.7
}
