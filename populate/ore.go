package populate

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/t14raptor/pm-gen/rand"
	"math"
)

type Ore struct {
	Types []OreType
}

func (o Ore) Populate(w *world.World, pos world.ChunkPos, chunk *chunk.Chunk, r *rand.Random) {
	for _, ore := range o.Types {
		for i := 0; i < ore.ClusterCount; i++ {
			pos := cube.Pos{
				int(r.Range(pos[0]<<4, (pos[0]<<4)+15)),
				int(r.Range(int32(ore.MinHeight), int32(ore.MaxHeight))),
				int(r.Range(pos[1]<<4, (pos[1]<<4)+15)),
			}
			if w.Block(pos) == ore.Replaces {
				ore.Place(w, pos, r)
			}
		}
	}
}

type OreType struct {
	Material, Replaces        world.Block
	ClusterCount, ClusterSize int
	MinHeight, MaxHeight      int
}

func (o OreType) Place(w *world.World, pos cube.Pos, r *rand.Random) {
	chunkminx := (pos.X() >> 4) << 4
	chunkminz := (pos.Z() >> 4) << 4
	chunkmaxx := chunkminx + 15
	chunkmaxz := chunkminz + 15

	clusterSize := float64(o.ClusterSize)
	angle := r.Float64() * math.Pi
	offset := mgl64.Vec2{math.Cos(angle), math.Sin(angle)}.Mul(clusterSize / 8)
	x1, x2 := float64(pos.X())+8+offset[0], float64(pos.X())+8-offset[0]
	z1, z2 := float64(pos.Z())+8+offset[1], float64(pos.Z())+8-offset[1]
	y1, y2 := float64(pos.Y())+float64(r.Int31n(3))+2, float64(pos.Y())+float64(r.Int31n(3))+2
	for i := float64(0); i <= clusterSize; i++ {
		seedX := x1 + (x2-x1)*i/clusterSize
		seedY := y1 + (y2-y1)*i/clusterSize
		seedZ := z1 + (z2-z1)*i/clusterSize
		size := ((math.Sin(i*(math.Pi/clusterSize))+1)*r.Float64()*clusterSize/16 + 1) / 2

		startX := float64(int(seedX - size))
		startY := float64(int(seedY - size))
		startZ := float64(int(seedZ - size))
		endX := float64(int(seedX + size))
		endY := float64(int(seedY + size))
		endZ := float64(int(seedZ + size))

		for xx := startX; xx <= endX; xx++ {
			if int(xx) < chunkminx {
				continue
			}
			if int(xx) > chunkmaxx {
				continue
			}
			sizeX := (xx + 0.5 - seedX) / size
			sizeX *= sizeX

			if sizeX < 1 {
				for yy := startY; yy <= endY; yy++ {
					sizeY := (yy + 0.5 - seedY) / size
					sizeY *= sizeY

					if yy > 0 && (sizeX+sizeY) < 1 {
						for zz := startZ; zz <= endZ; zz++ {
							if int(zz) < chunkminz {
								continue
							}
							if int(zz) > chunkmaxz {
								continue
							}

							sizeZ := (zz + 0.5 - seedZ) / size
							sizeZ *= sizeZ

							pos := cube.Pos{int(xx), int(yy), int(zz)}

							if (sizeX+sizeY+sizeZ) < 1 && w.Block(pos) == o.Replaces {
								w.SetBlock(pos, o.Material, &world.SetOpts{DisableBlockUpdates: true, DisableLiquidDisplacement: true})
							}
						}
					}
				}
			}
		}
	}
}
