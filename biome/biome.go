package biome

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/t14raptor/pm-gen/populate"
)

var biomes = make(map[uint8]Biome)

func Register(id uint8, b Biome) {
	biomes[id] = b
}

func BiomeByID(id uint8) Biome {
	return biomes[id]
}

func init() {
	Register(IDOcean, Ocean{})
	Register(IDPlains, Plains{})
	Register(IDDesert, Desert{})
	Register(IDMountains, Mountains{})
	Register(IDForest, Forest{})
	Register(IDTaiga, Taiga{})
	Register(IDSwamp, Swamp{})
	Register(IDRiver, River{})
	Register(IDIcePlains, IcePlains{})
	Register(IDSmallMountains, SmallMountains{})
	Register(IDBirchForest, BirchForest{})
}

const (
	IDOcean byte = iota
	IDPlains
	IDDesert
	IDMountains
	IDForest
	IDTaiga
	IDSwamp
	IDRiver
	IDIcePlains      = 12
	IDSmallMountains = 20
	IDBirchForest    = 27
)

type Biome interface {
	Populators() []populate.Populator
	ID() uint8
	Elevation() (min, max int)
	GroundCover() []world.Block
	Temperature() float64
	Rainfall() float64
}

type grassy struct{}

func (grassy) GroundCover() []world.Block {
	return []world.Block{
		block.Grass{},
		block.Dirt{},
		block.Dirt{},
		block.Dirt{},
		block.Dirt{},
	}
}

type sandy struct{}

func (sandy) GroundCover() []world.Block {
	return []world.Block{
		block.Sand{},
		block.Sand{},
		block.Sandstone{},
		block.Sandstone{},
		block.Sandstone{},
	}
}

type snowy struct{}

func (snowy) GroundCover() []world.Block {
	return []world.Block{
		block.Grass{},
		block.Dirt{},
		block.Dirt{},
		block.Dirt{},
	}
}
