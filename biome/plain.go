package biome

import (
	"github.com/t14raptor/pm-gen/populate"
)

type Plains struct {
	grassy
}

func (p Plains) Populators() []populate.Populator {
	return []populate.Populator{populate.TallGrass{Amount: 12}}
}

func (p Plains) ID() uint8 {
	return IDPlains
}

func (p Plains) Elevation() (min, max int) {
	return 63, 68
}

func (p Plains) Temperature() float64 {
	return 0.8
}

func (p Plains) Rainfall() float64 {
	return 0.4
}
