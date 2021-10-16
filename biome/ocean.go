package biome

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/t14raptor/pm-gen/populate"
)

type Ocean struct{}

func (o Ocean) Populators() []populate.Populator {
	return []populate.Populator{populate.TallGrass{Amount: 5}}
}

func (o Ocean) ID() uint8 {
	return IDOcean
}

func (o Ocean) Elevation() (min, max int) {
	return 46, 58
}

func (o Ocean) GroundCover() []world.Block {
	return []world.Block{
		block.Gravel{},
		block.Gravel{},
		block.Gravel{},
		block.Gravel{},
		block.Gravel{},
	}
}

func (o Ocean) Temperature() float64 {
	panic("implement me")
}

func (o Ocean) Rainfall() float64 {
	panic("implement me")
}
