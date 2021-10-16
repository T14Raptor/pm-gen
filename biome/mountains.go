package biome

import "github.com/t14raptor/pm-gen/populate"

type Mountains struct {
	grassy
}

func (m Mountains) Populators() []populate.Populator {
	return nil
}

func (m Mountains) ID() uint8 {
	return IDMountains
}

func (m Mountains) Elevation() (min, max int) {
	return 63, 127
}

func (m Mountains) Temperature() float64 {
	return 0.4
}

func (m Mountains) Rainfall() float64 {
	return 0.5
}
