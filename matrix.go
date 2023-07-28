package geom

import (
	"fmt"
	"math"
)

type Geom [3][2]float64

func seven(f float64) string { return fmt.Sprintf("%-#.10g", f) }

func (g *Geom) String() string {
	return fmt.Sprint("⎧ ", g[0][0], " ", g[0][1], " 0 ⎫\n") +
		fmt.Sprint("⎨ ", g[1][0], " ", g[1][1], " 0 ⎬\n") +
		fmt.Sprint("⎩ ", g[2][0], " ", g[2][1], " 1 ⎭")
}

func (g Geom) Mul(b Geom) (c Geom) {
	for i := range g {
		for j := range b[0] {
			for k := range g[0] {
				c[i][j] += g[i][k] * b[k][j]
			}
		}
	}
	return c
}

func (g Geom) Transpose() (b Geom) {
	for i := range g {
		for j := range g[0] {
			b[j][i] = g[i][j]
		}
	}
	return
}

func (g Geom) Inverse() (b Geom) {
	a, bb, c, d, e, f, gg, h, i :=
		g[0][0], g[0][1], 0.0,
		g[1][0], g[1][1], 0.0,
		g[2][0], g[2][1], 1.0
	b = Geom{
		{e*i - f*h, -(bb*i - c*h)},
		{-(d*i - f*gg), a*i - c*gg},
		{d*h - e*gg, -(a*h - bb*gg)},
	}
	det := g[0][0]*b[0][0] + g[0][1]*b[1][0] + 0.0*b[2][0]
	for i := range b {
		for j := range b[0] {
			b[i][j] /= det
		}
	}
	return
}

func (g Geom) Apply(vec [3]float64) (o [3]float64) {
	for i := range g {
		for j := range vec {
			o[i] += g[j][i] * vec[j]
		}
	}
	return o
}

func (g Geom) ApplyPt(p Point) Point {
	v := g.Apply([3]float64{p.X, p.Y, 1})
	return Pt(v[0], v[1])
}

func (g Geom) ApplyRect(r Rectangle) Rectangle {
	return Rectangle{Min: g.ApplyPt(r.Min), Max: g.ApplyPt(r.Max)}
}

func (g Geom) Scale(x, y float64) Geom {
	return g.Mul(Scale2d(x, y))
}

func (g Geom) Translate(x, y float64) Geom {
	return g.Mul(Translate2d(x, y))
}

func (g Geom) Rotate(theta float64) Geom {
	return g.Mul(Rotate2d(theta))
}

func (g Geom) Shear(x, y float64) Geom {
	return g.Mul(Shear2d(x, y))
}

func (g Geom) To3d() Geom3d {
	return Geom3d{
		{g[0][0], g[0][1], 0, 0},
		{g[1][0], g[1][1], 0, 0},
		{g[2][0], g[2][1], 1, 0},
		{0, 0, 0, 1},
	}
}

func Identity2d() (g Geom) {
	for i := range g {
		g[i][i] = 1
	}
	return
}

func Scale2d(x, y float64) (p Geom) {
	p[0][0] = x
	p[1][1] = y
	// p[2][2] = 1
	return p
}

func Translate2d(x, y float64) Geom {
	p := Identity2d()
	p[2][0] = x
	p[2][1] = y
	return p
}

func Rotate2d(theta float64) (p Geom) {
	// in 2d we have only roll
	// p[2][2] = 1

	p[0][0] = math.Cos(theta)
	p[0][1] = math.Sin(theta)

	p[1][0] = -math.Sin(theta)
	p[1][1] = +math.Cos(theta)

	return
}

func Shear2d(x, y float64) (p Geom) {
	p = Identity2d()
	p[0][1] = y
	p[1][0] = x
	return
}

func Window2d(w, h int) Geom {
	return Scale2d(2/float64(w), -2/float64(h)).Translate(-1, 1)
}

type Geom3d [4][4]float64

func (g Geom3d) Mul(b Geom3d) (c Geom3d) {
	for i := range g {
		for j := range b[0] {
			for k := range g[0] {
				c[i][j] += g[i][k] * b[k][j]
			}
		}
	}
	return c
}

func (g Geom3d) Transpose() (b Geom3d) {
	for i := range g {
		for j := range g[0] {
			b[j][i] = g[i][j]
		}
	}
	return
}

func (g Geom3d) Apply(vec [4]float64) (o [4]float64) {
	for i := range g {
		for j := range vec {
			o[i] += g[i][j] * vec[j]
		}
	}
	return o
}

func (g Geom3d) Scale(x, y, z float64) Geom3d {
	return g.Mul(Scale3d(x, y, z))
}

func (g Geom3d) Translate(x, y, z float64) Geom3d {
	return g.Mul(Translate3d(x, y, z))
}

func (g Geom3d) Rotate(pitch, yaw, roll float64) Geom3d {
	return g.Mul(Rotate3d(pitch, yaw, roll))
}

func Identity3d() (g Geom3d) {
	for i := range g {
		g[i][i] = 1
	}
	return
}

func Scale3d(x, y, z float64) (p Geom3d) {
	p[0][0] = x
	p[1][1] = y
	p[2][2] = z
	p[3][3] = 1
	return p
}

func Translate3d(x, y, z float64) Geom3d {
	p := Identity3d()
	p[3][0] = x
	p[3][1] = y
	p[3][2] = z
	return p
}

// Rotate3d generates a 3d rotation matrix based on Euler angles.
func Rotate3d(yaw, pitch, roll float64) (p Geom3d) {
	// https://math.stackexchange.com/questions/2796055/3d-coordinate-rotation-using-roll-pitch-yaw
	ϕ := yaw
	θ := pitch
	ψ := roll
	sin := math.Sin
	cos := math.Cos
	p[3][3] = 1

	p[0][0] = cos(θ) * cos(ψ)
	p[0][1] = cos(θ) * sin(ψ)
	p[0][2] = -sin(θ)

	p[1][0] = sin(ϕ)*sin(θ)*cos(ψ) - cos(ϕ)*sin(ψ)
	p[1][1] = sin(ϕ)*sin(θ)*sin(ψ) + cos(ϕ)*cos(ψ)
	p[1][2] = sin(ϕ) * cos(θ)

	p[2][0] = cos(ϕ)*sin(θ)*cos(ψ) + sin(ϕ)*sin(ψ)
	p[2][1] = cos(ϕ)*sin(θ)*sin(ψ) - sin(ϕ)*cos(ψ)
	p[2][2] = cos(ϕ) * cos(θ)

	return
}

func Frustum3d(left, right, bottom, top, near, far float64) Geom3d {
	return Geom3d{
		{2 * near / (right - left), 0, (right + left) / (right - left), 0},
		{0, 2 * near / (top - bottom), (top + bottom) / (top - bottom), 0},
		{0, 0, -(far + near) / (far - near), -2 * far * near / (far - near)},
		{0, 0, -1, 0},
	}.Transpose()
}

// Perspective3d generates a perspective projection matrix.
// Fov is the angle between top and bottom planes of the frustum.
// Aspect is its aspect ratio.
func Perspective3d(fov, aspect, near, far float64) Geom3d {
	// https://github.com/g-truc/glm/blob/master/glm/ext/matrix_clip_space.inl#L233
	halftan := math.Tan(fov / 2)
	return Geom3d{
		{1 / aspect / halftan, 0, 0, 0},
		{0, 1 / halftan, 0, 0},
		{0, 0, far / (near - far), -1},
		{0, 0, -far * near / (far - near), 0},
	}
}

// Ortho3d generates an orthographic projection matrix.
// Zoom is the relative size of one unit in world coordinates.
// If width and height are equal to window width and height respectively,
// zoom is the size of one uint in pixels.
func Ortho3d(zoom, width, height, near, far float64) Geom3d {
	zoom *= 2

	// To get this matrix, take a regular orthographic projection matrix and use following conversions:
	// left = -width/zoom
	// right = width/zoom
	// top = -height/zoom
	// bottom = height/zoom
	return Geom3d{
		{1 / (width / zoom), 0, 0, 0},
		{0, 1 / (height / zoom), 0, 0},
		{0, 0, -2 / (far - near), 0},
		{0, 0, -(far + near) / (far - near), 1},
	}
}
