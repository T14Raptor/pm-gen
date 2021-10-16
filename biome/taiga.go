package biome

import (
	"github.com/t14raptor/pm-gen/populate"
)

type Taiga struct {
	snowy
}

func (t Taiga) Populators() []populate.Populator {
	return []populate.Populator{populate.Tree{Type: populate.SpruceTree{}, BaseAmount: 10}, populate.TallGrass{Amount: 1}}
}

func (t Taiga) ID() uint8 {
	return IDTaiga
}

func (t Taiga) Elevation() (min, max int) {
	return 63, 81
}

func (t Taiga) Temperature() float64 {
	return 0.05
}

func (t Taiga) Rainfall() float64 {
	return 0.8
}
