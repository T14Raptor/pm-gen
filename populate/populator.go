package populate

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/t14raptor/pm-gen/rand"
)

type Populator interface {
	Populate(w *world.World, pos world.ChunkPos, chunk *chunk.Chunk, r *rand.Random)
}
