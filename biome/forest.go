package biome

import (
	"github.com/t14raptor/pm-gen/populate"
)

type Forest struct {
	grassy
}

func (f Forest) Populators() []populate.Populator {
	return []populate.Populator{populate.Tree{Type: populate.OakTree{}, BaseAmount: 5}, populate.TallGrass{Amount: 3}}
}

func (f Forest) ID() uint8 {
	return IDForest
}

func (f Forest) Elevation() (min, max int) {
	return 63, 81
}

func (f Forest) Temperature() float64 {
	return 0.7
}

func (f Forest) Rainfall() float64 {
	return 0.8
}
