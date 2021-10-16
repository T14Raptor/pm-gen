package pmgen

import (
	"github.com/t14raptor/pm-gen/biome"
	"github.com/t14raptor/pm-gen/rand"
)

type biomeSelector struct {
	temp, rainfall *simplex

	m [64 * 64]biome.Biome
}

func newBiomeSelector(r *rand.Random) *biomeSelector {
	return &biomeSelector{
		temp:     newSimplex(r, 2, 1.0/16, 1.0/512),
		rainfall: newSimplex(r, 2, 1.0/16, 1.0/512),
	}
}

func (s *biomeSelector) recalculate() {
	for i := 0; i < 64; i++ {
		for j := 0; j < 64; j++ {
			b := biome.BiomeByID(s.lookup(float64(i)/63, float64(j)/63))
			s.m[i+(j<<6)] = b
		}
	}
}

func (*biomeSelector) lookup(temp, rain float64) uint8 {
	if rain < 0.25 {
		if temp < 0.7 {
			return biome.IDOcean
		} else if temp < 0.85 {
			return biome.IDRiver
		} else {
			return biome.IDSwamp
		}
	} else if rain < 0.6 {
		if temp < 0.25 {
			return biome.IDIcePlains
		} else if temp < 0.75 {
			return biome.IDPlains
		} else {
			return biome.IDDesert
		}
	} else if rain < 0.8 {
		if temp < 0.25 {
			return biome.IDTaiga
		} else if temp < 0.75 {
			return biome.IDForest
		} else {
			return biome.IDBirchForest
		}
	} else {
		return biome.IDRiver
	}
}

func (b *biomeSelector) Temperature(x, z float64) float64 {
	return (b.temp.OctaveNoise2D(x, z, true) + 1) / 2
}

func (b *biomeSelector) Rainfall(x, z float64) float64 {
	return (b.rainfall.OctaveNoise2D(x, z, true) + 1) / 2
}

func (b *biomeSelector) pickBiome(x, z int64) biome.Biome {
	temp := int(b.Temperature(float64(x), float64(z)) * 63)
	rain := int(b.Rainfall(float64(x), float64(z)) * 63)
	return b.m[temp+(rain<<6)]
}
