package biome

import (
	"github.com/t14raptor/pm-gen/populate"
)

type BirchForest struct {
	grassy
}

func (f BirchForest) Populators() []populate.Populator {
	return []populate.Populator{populate.Tree{Type: populate.BirchTree{}, BaseAmount: 5}, populate.TallGrass{Amount: 3}}
}

func (f BirchForest) ID() uint8 {
	return IDBirchForest
}

func (f BirchForest) Elevation() (min, max int) {
	return 63, 81
}

func (f BirchForest) Temperature() float64 {
	return 0.6
}

func (f BirchForest) Rainfall() float64 {
	return 0.5
}
