package populate

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/t14raptor/pm-gen/rand"
)

type TallGrass struct {
	Amount int
}

var (
	air       = block.Air{}
	grass     = block.Grass{}
	tallGrass = block.TallGrass{}
)

func (t TallGrass) Populate(w *world.World, pos world.ChunkPos, chunk *chunk.Chunk, r *rand.Random) {
	amount := r.Int31n(2) + int32(t.Amount)
	for i := int32(0); i < amount; i++ {
		x, z := int(r.Range(pos[0]*16, pos[0]*16+15)), int(r.Range(pos[1]*16, pos[1]*16+15))
		if y, ok := t.highestWorkableBlock(w, x, z); ok {
			w.SetBlock(cube.Pos{x, y, z}, tallGrass, &world.SetOpts{DisableBlockUpdates: true, DisableLiquidDisplacement: true})
		}
	}
}

func (t TallGrass) highestWorkableBlock(w *world.World, x, z int) (int, bool) {
	var next world.Block
	for y := 127; y >= 0; y-- {
		b := next
		if b == nil {
			b = w.Block(cube.Pos{x, y, z})
		}
		next = w.Block(cube.Pos{x, y - 1, z})
		if b == air && next == grass {
			return y, true
		}
	}
	return 0, false
}
