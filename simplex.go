package pmgen

import (
	"github.com/t14raptor/pm-gen/rand"
)

type Grad3 struct {
	x, y, z float64
}

var grad3Table = [...]Grad3{
	{1, 1, 0},
	{-1, 1, 0},
	{1, -1, 0},
	{-1, -1, 0},
	{1, 0, 1},
	{-1, 0, 1},
	{1, 0, -1},
	{-1, 0, -1},
	{0, 1, 1},
	{0, -1, 1},
	{0, 1, -1},
	{0, -1, -1},
}

func Dot3(g Grad3, x, y, z float64) float64 {
	return g.x*x + g.y*y + g.z*z
}

func Dot2(g Grad3, x, y float64) float64 {
	return g.x*x + g.y*y
}

const (
	S3 = 1.7320508075688772
	F2 = 0.5 * (S3 - 1)
	G2 = (3 - S3) / 6
	F3 = 1.0 / 3.0
	G3 = 1.0 / 6.0
)

type simplex struct {
	perm                      []int
	permMod12                 []int
	offsetX, offsetY, offsetZ float64
	octaves                   int
	persistence               float64
	expansion                 float64
}

func newSimplex(r *rand.Random, octaves int, persistence, expansion float64) *simplex {
	s := &simplex{
		octaves:     octaves,
		persistence: persistence,
		expansion:   expansion,
		offsetX:     r.Float64() * 256,
		offsetY:     r.Float64() * 256,
		offsetZ:     r.Float64() * 256,
	}

	s.perm = make([]int, 512)
	for i := 0; i < 256; i++ {
		s.perm[i] = int(r.Int31()) % 256
	}

	for i := int32(0); i < 256; i++ {
		pos := (r.Int31() % (256 - i)) + i
		old := s.perm[i]

		s.perm[i] = s.perm[pos]
		s.perm[pos] = old
		s.perm[i+256] = s.perm[i]
	}

	s.permMod12 = make([]int, 512)
	for i, p := range s.perm {
		s.permMod12[i] = p % 12
	}

	r.Float64()

	return s
}

func (s *simplex) getFastNoise3D(xSize, ySize, zSize int64, xSamplingRate, ySamplingRate, zSamplingRate int64, x, y, z int64) [][][]float64 {
	noiseArray := make([][][]float64, xSize+1)
	for i := int64(0); i <= (xSize); i++ {
		noiseArray[i] = make([][]float64, zSize+1)
		for j := int64(0); j <= (zSize); j++ {
			noiseArray[i][j] = make([]float64, ySize+1)
		}
	}
	for xx := int64(0); xx <= xSize; xx += xSamplingRate {
		for zz := int64(0); zz <= zSize; zz += zSamplingRate {
			for yy := int64(0); yy <= ySize; yy += ySamplingRate {
				noiseArray[xx][zz][yy] = s.OctaveNoise3D(float64(xx+x), float64(y+yy), float64(z+zz), true)
			}
		}
	}

	for xx := int64(0); xx < (xSize); xx++ {
		for zz := int64(0); zz < (zSize); zz++ {
			for yy := int64(0); yy < (ySize); yy++ {
				if xx%xSamplingRate != 0 || yy%ySamplingRate != 0 || zz%zSamplingRate != 0 {
					nx := (xx / xSamplingRate) * xSamplingRate
					ny := (yy / ySamplingRate) * ySamplingRate
					nz := (zz / zSamplingRate) * zSamplingRate

					noiseArray[xx][zz][yy] = triLerp(float64(xx), float64(yy), float64(zz),
						noiseArray[nx][nz][ny], noiseArray[nx][nz][ny+ySamplingRate],
						noiseArray[nx][nz+zSamplingRate][ny], noiseArray[nx][nz+zSamplingRate][ny+ySamplingRate],
						noiseArray[nx+xSamplingRate][nz][ny], noiseArray[nx+xSamplingRate][nz][ny+ySamplingRate],
						noiseArray[nx+xSamplingRate][nz+zSamplingRate][ny], noiseArray[nx+xSamplingRate][nz+zSamplingRate][ny+ySamplingRate],
						float64(nx), float64(nx+xSamplingRate), float64(ny), float64(ny+ySamplingRate), float64(nz), float64(nz+zSamplingRate))
				}
			}
		}
	}

	return noiseArray
}

func triLerp(x, y, z, q000, q001, q010, q011, q100, q101, q110, q111, x1, x2, y1, y2, z1, z2 float64) float64 {
	q00 := lerp(x, x1, x2, q000, q100)
	q01 := lerp(x, x1, x2, q010, q110)
	q10 := lerp(x, x1, x2, q001, q101)
	q11 := lerp(x, x1, x2, q011, q111)
	q0 := lerp(y, y1, y2, q00, q10)
	q1 := lerp(y, y1, y2, q01, q11)
	return lerp(z, z1, z2, q0, q1)
}

func lerp(x, x1, x2, q0, q1 float64) float64 {
	return ((x2-x)/(x2-x1))*q0 + ((x-x1)/(x2-x1))*q1
}

func (s *simplex) OctaveNoise3D(xin, yin, zin float64, normalized bool) (result float64) {
	freq := 1.0
	amp := 1.0
	max := 0.0

	xin, yin, zin = xin*s.expansion, yin*s.expansion, zin*s.expansion

	for i := 0; i < s.octaves; i++ {
		result += s.Noise3D(xin*freq, yin*freq, zin*freq) * amp
		max += amp
		freq *= 2.0
		amp *= s.persistence
	}

	if normalized {
		result /= max
	}
	return
}

func (s *simplex) Noise3D(xin, yin, zin float64) float64 {
	var (
		x, y, z, t [4]float64
		i, j, k    [3]int
	)

	xin += s.offsetX
	yin += s.offsetY
	zin += s.offsetZ

	skew := (xin + yin + zin) * F3
	i[0] = int(xin + skew)
	j[0] = int(yin + skew)
	k[0] = int(zin + skew)
	tt := float64(i[0]+j[0]+k[0]) * G3

	x[0] = xin - (float64(i[0]) - tt)
	y[0] = yin - (float64(j[0]) - tt)
	z[0] = zin - (float64(k[0]) - tt)

	if x[0] >= y[0] {
		if y[0] >= z[0] {
			i[1], j[1], k[1], i[2], j[2], k[2] = 1, 0, 0, 1, 1, 0
		} else if x[0] >= z[0] {
			i[1], j[1], k[1], i[2], j[2], k[2] = 1, 0, 0, 1, 0, 1
		} else {
			i[1], j[1], k[1], i[2], j[2], k[2] = 0, 0, 1, 1, 0, 1
		}
	} else {
		if y[0] < z[0] {
			i[1], j[1], k[1], i[2], j[2], k[2] = 0, 0, 1, 0, 1, 1
		} else if x[0] < z[0] {
			i[1], j[1], k[1], i[2], j[2], k[2] = 0, 1, 0, 0, 1, 1
		} else {
			i[1], j[1], k[1], i[2], j[2], k[2] = 0, 1, 0, 1, 1, 0
		}
	}

	x[1] = x[0] - float64(i[1]) + G3
	y[1] = y[0] - float64(j[1]) + G3
	z[1] = z[0] - float64(k[1]) + G3
	x[2] = x[0] - float64(i[2]) + 2.0*G3
	y[2] = y[0] - float64(j[2]) + 2.0*G3
	z[2] = z[0] - float64(k[2]) + 2.0*G3
	x[3] = x[0] - 1.0 + 3.0*G3
	y[3] = y[0] - 1.0 + 3.0*G3
	z[3] = z[0] - 1.0 + 3.0*G3

	ii := i[0] & 255
	jj := j[0] & 255
	kk := k[0] & 255

	n := 0.0

	t[0] = 0.6 - x[0]*x[0] - y[0]*y[0] - z[0]*z[0]
	if t[0] > 0 {
		n += t[0] * t[0] * t[0] * t[0] * Dot3(grad3Table[s.permMod12[ii+s.perm[jj+s.perm[kk]]]], x[0], y[0], z[0])
	}

	t[1] = 0.6 - x[1]*x[1] - y[1]*y[1] - z[1]*z[1]
	if t[1] > 0 {
		n += t[1] * t[1] * t[1] * t[1] * Dot3(grad3Table[s.permMod12[ii+i[1]+s.perm[jj+j[1]+s.perm[kk+k[1]]]]], x[1], y[1], z[1])
	}

	t[2] = 0.6 - x[2]*x[2] - y[2]*y[2] - z[2]*z[2]
	if t[2] > 0 {
		n += t[2] * t[2] * t[2] * t[2] * Dot3(grad3Table[s.permMod12[ii+i[2]+s.perm[jj+j[2]+s.perm[kk+k[2]]]]], x[2], y[2], z[2])
	}

	t[3] = 0.6 - x[3]*x[3] - y[3]*y[3] - z[3]*z[3]
	if t[3] > 0 {
		n += t[3] * t[3] * t[3] * t[3] * Dot3(grad3Table[s.permMod12[ii+1+s.perm[jj+1+s.perm[kk+1]]]], x[3], y[3], z[3])
	}

	return 32.0 * n
}

func (s *simplex) OctaveNoise2D(xin, zin float64, normalized bool) (result float64) {
	freq := 1.0
	amp := 1.0
	maxAmp := 0.0

	xin, zin = xin*s.expansion, zin*s.expansion

	for i := 0; i < s.octaves; i++ {
		result += s.Noise2D(xin*freq, zin*freq) * amp
		freq *= 2.0
		maxAmp += amp
		amp *= s.persistence
	}

	if normalized {
		result /= maxAmp
	}

	return
}

func (s *simplex) Noise2D(xin, yin float64) float64 {
	var (
		n, x, y, t [3]float64
		gi         [3]int
		i, j       [2]int
	)

	xin += s.offsetX
	yin += s.offsetY

	skew := (xin + yin) * F2
	i[0] = int(xin + skew)
	j[0] = int(yin + skew)
	tt := float64(i[0]+j[0]) * G2
	x[0] = xin - (float64(i[0]) - tt)
	y[0] = yin - (float64(j[0]) - tt)

	if x[0] > y[0] {
		i[1], j[1] = 1, 0
	} else {
		i[1], j[1] = 0, 1
	}

	x[1] = x[0] - float64(i[1]) + G2
	y[1] = y[0] - float64(j[1]) + G2
	x[2] = x[0] + (G2*2 - 1.0)
	y[2] = y[0] + (G2*2 - 1.0)

	ii := i[0] & 255
	jj := j[0] & 255
	gi[0] = s.permMod12[ii+s.perm[jj]]
	gi[1] = s.permMod12[ii+i[1]+s.perm[jj+j[1]]]
	gi[2] = s.permMod12[ii+1+s.perm[jj+1]]

	for i := 0; i < len(n); i++ {
		t[i] = 0.5 - x[i]*x[i] - y[i]*y[i]
		if t[i] > 0 {
			n[i] = t[i] * t[i] * t[i] * t[i] * Dot2(grad3Table[gi[i]], x[i], y[i])
		}
	}

	return 70.0 * (n[0] + n[1] + n[2])
}
