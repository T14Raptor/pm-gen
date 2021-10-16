package rand

const (
	x = 123456789
	y = 362436069
	z = 521288629
	w = 88675123
)

type Random struct {
	x, y, z, w int64
	seed       int64
}

func NewRandom(seed int64) *Random {
	return &Random{
		seed: seed,
		x:    x ^ seed,
		y:    y ^ (seed << 17) | ((seed>>15)&0x7fffffff)&0xffffffff,
		z:    z ^ (seed << 31) | ((seed>>1)&0x7fffffff)&0xffffffff,
		w:    w ^ (seed << 18) | ((seed>>14)&0x7fffffff)&0xffffffff,
	}
}

func (r *Random) SetSeed(seed int64) {
	r.seed = seed
	r.x = x ^ seed
	r.y = y ^ (seed << 17) | ((seed>>15)&0x7fffffff)&0xffffffff
	r.z = z ^ (seed << 31) | ((seed>>1)&0x7fffffff)&0xffffffff
	r.w = w ^ (seed << 18) | ((seed>>14)&0x7fffffff)&0xffffffff
}

func (r *Random) Int31() int32 {
	return r.Int32() & 0x7fffffff
}

func (r *Random) Int32() int32 {
	t := (r.x ^ (r.x << 11)) & 0xffffffff

	r.x = r.y
	r.y = r.z
	r.z = r.w
	r.w = (r.w ^ ((r.w >> 19) & 0x7fffffff) ^ (t ^ ((t >> 8) & 0x7fffffff))) & 0xffffffff

	return int32(r.w)
}

func (r *Random) Int31n(n int32) int32 {
	return r.Int31() % n
}

func (r *Random) Range(start, end int32) int32 {
	return start + r.Int31n(end+1-start)
}

func (r *Random) Float64() float64 {
	return float64(r.Int31()) / 0x7fffffff
}
