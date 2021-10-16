package biome

import (
	"github.com/t14raptor/pm-gen/populate"
)

type IcePlains struct {
	snowy
}

func (i IcePlains) Populators() []populate.Populator {
	return []populate.Populator{populate.TallGrass{Amount: 5}}
}

func (i IcePlains) ID() uint8 {
	return IDIcePlains
}

func (i IcePlains) Elevation() (min, max int) {
	return 63, 74
}

func (i IcePlains) Temperature() float64 {
	return 0.05
}

func (i IcePlains) Rainfall() float64 {
	return 0.8
}
