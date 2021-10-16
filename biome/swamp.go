package biome

import "github.com/t14raptor/pm-gen/populate"

type Swamp struct {
	grassy
}

func (s Swamp) Populators() []populate.Populator {
	return nil
}

func (s Swamp) ID() uint8 {
	return IDSwamp
}

func (s Swamp) Elevation() (min, max int) {
	return 62, 63
}

func (s Swamp) Temperature() float64 {
	return 0.8
}

func (s Swamp) Rainfall() float64 {
	return 0.9
}
