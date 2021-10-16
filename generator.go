package pmgen

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/t14raptor/pm-gen/biome"
	"github.com/t14raptor/pm-gen/populate"
	"github.com/t14raptor/pm-gen/rand"
)

const SmoothSize = 2

var gaussianKernel = [5][5]float64{
	{
		1.4715177646858,
		2.141045714076,
		2.4261226388505,
		2.141045714076,
		1.4715177646858,
	},
	{
		2.141045714076,
		3.1152031322856,
		3.5299876103384,
		3.1152031322856,
		2.141045714076,
	},
	{
		2.4261226388505,
		3.5299876103384,
		4,
		3.5299876103384,
		2.4261226388505,
	},
	{
		2.141045714076,
		3.1152031322856,
		3.5299876103384,
		3.1152031322856,
		2.141045714076,
	},
	{
		1.4715177646858,
		2.141045714076,
		2.4261226388505,
		2.141045714076,
		1.4715177646858,
	},
}

type Generator struct {
	seed        int64
	waterHeight int
	noise       *simplex
	selector    *biomeSelector

	populationQueue chan PopulationEntry
}

var (
	bedrock, _ = world.BlockRuntimeID(block.Bedrock{})
	stone, _   = world.BlockRuntimeID(block.Stone{})
	air, _     = world.BlockRuntimeID(block.Air{})
	water, _   = world.BlockRuntimeID(block.Water{Depth: 8, Still: true})
)

type PopulationEntry struct {
	ChunkPos  world.ChunkPos
	Chunk     *chunk.Chunk
	Populator populate.Populator
	Random    *rand.Random
}

func New(w *world.World, seed int64) *Generator {
	r := rand.NewRandom(seed)
	noise := newSimplex(r, 4, 1.0/4, 1.0/32)
	r.SetSeed(seed)
	selector := newBiomeSelector(r)
	selector.recalculate()

	g := &Generator{
		seed:            seed,
		noise:           noise,
		selector:        selector,
		populationQueue: make(chan PopulationEntry, 99999),
	}

	go g.populate(w)

	return g
}

func (g *Generator) populate(w *world.World) {
	for populator := range g.populationQueue {
		populator.Populator.Populate(w, populator.ChunkPos, populator.Chunk, populator.Random)
	}
}

func (g *Generator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
	r := rand.NewRandom(0xdeadbeef ^ (int64(pos[0]) << 8) ^ int64(pos[1]) ^ g.seed)

	noise := g.noise.getFastNoise3D(16, 128, 16, 4, 8, 4, int64(pos[0])*16, 0, int64(pos[1])*16)

	var biomeCache = make(map[[2]int64]biome.Biome)
	for x := int64(0); x < 16; x++ {
		for z := int64(0); z < 16; z++ {
			var minSum, maxSum, weightSum float64

			b := g.pickBiome(int64(pos[0])*16+x, int64(pos[1])*16+z)
			chunk.SetBiomeID(uint8(x), uint8(z), b.ID())

			for sx := int64(-SmoothSize); sx <= SmoothSize; sx++ {
				for sz := int64(-SmoothSize); sz <= SmoothSize; sz++ {
					weight := gaussianKernel[sx+SmoothSize][sz+SmoothSize]

					var adjacent biome.Biome
					if sx == 0 && sz == 0 {
						adjacent = b
					} else {
						i := [2]int64{int64(pos[0])*16 + x + sx, int64(pos[1])*16 + z + sz}
						if bc, ok := biomeCache[i]; ok {
							adjacent = bc
						} else {
							adjacent = g.pickBiome(i[0], i[1])
							biomeCache[i] = adjacent
						}
					}

					min, max := adjacent.Elevation()
					minSum += float64(min-1) * weight
					maxSum += float64(max) * weight

					weightSum += weight
				}
			}

			minSum /= weightSum
			maxSum /= weightSum

			smoothHeight := (maxSum - minSum) / 2

			for y := 0; y < 128; y++ {
				if y == 0 {
					chunk.SetRuntimeID(uint8(x), int16(y), uint8(z), 0, bedrock)
					continue
				}
				const waterHeight = 62

				noiseValue := noise[x][z][y] - 1.0/smoothHeight*(float64(y)-smoothHeight-minSum)
				if noiseValue > 0 {
					chunk.SetRuntimeID(uint8(x), int16(y), uint8(z), 0, stone)
				} else if y <= waterHeight {
					chunk.SetRuntimeID(uint8(x), int16(y), uint8(z), 0, water)
				}
			}
		}
	}

	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			b := biome.BiomeByID(chunk.BiomeID(x, z))
			c := b.GroundCover()
			if len(c) > 0 {
				var diffY int16
				if (c[0].Model() != model.Solid{}) {
					diffY = 1
				}

				start := min(127, chunk.HighestLightBlocker(x, z)+diffY)
				end := start - int16(len(c))
				for y := start; y > end && y >= 0; y-- {
					b := c[start-y]
					r := chunk.RuntimeID(x, y, z, 0)
					if r == air && (b.Model() == model.Solid{}) {
						break
					}
					if _, ok := b.(block.LiquidRemovable); ok {
						bl, _ := world.BlockByRuntimeID(r)
						if _, ok = bl.(world.Liquid); ok {
							continue
						}
					}

					rid, _ := world.BlockRuntimeID(b)
					chunk.SetRuntimeID(x, y, z, 0, rid)
				}
			}
		}
	}

	bi := biome.BiomeByID(chunk.BiomeID(7, 7))

	for _, populator := range append([]populate.Populator{populate.Ore{Types: []populate.OreType{
		{block.CoalOre{}, block.Stone{}, 20, 16, 0, 128},
		{block.IronOre{}, block.Stone{}, 20, 8, 0, 64},
		//{ block.RedstoneOre{}, block.Stone{}, 8, 7, 0, 16 }, // TODO
		{block.LapisOre{}, block.Stone{}, 1, 6, 0, 32},
		{block.GoldOre{}, block.Stone{}, 2, 8, 0, 32},
		{block.DiamondOre{}, block.Stone{}, 1, 7, 0, 16},
		{block.Dirt{}, block.Stone{}, 20, 32, 0, 128},
		{block.Gravel{}, block.Stone{}, 10, 16, 0, 128},
	}}}, bi.Populators()...) {
		g.populationQueue <- PopulationEntry{ChunkPos: pos, Chunk: chunk, Populator: populator, Random: r}
	}
}

func (g *Generator) pickBiome(x, z int64) biome.Biome {
	hash := x*2345803 ^ z*9236449 ^ g.seed
	hash *= hash + 223
	xNoise := hash >> 20 & 3
	zNoise := hash >> 22 & 3
	if xNoise == 3 {
		xNoise = 1
	}
	if zNoise == 3 {
		zNoise = 1
	}

	return g.selector.pickBiome(x+xNoise-1, z+zNoise-1)
}

func min(a, b int16) int16 {
	if a < b {
		return a
	}
	return b
}
