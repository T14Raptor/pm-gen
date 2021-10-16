package biome

import "github.com/t14raptor/pm-gen/populate"

type SmallMountains struct {
	grassy
}

func (m SmallMountains) Populators() []populate.Populator {
	return nil
}

func (m SmallMountains) ID() uint8 {
	return IDMountains
}

func (m SmallMountains) Elevation() (min, max int) {
	return 63, 97
}

func (m SmallMountains) Temperature() float64 {
	return 0.4
}

func (m SmallMountains) Rainfall() float64 {
	return 0.5
}
