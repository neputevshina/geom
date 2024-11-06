package geom

import (
	"fmt"
	"image"
	"math"
)

type Point struct {
	X, Y float64
}

func Pt(x, y float64) Point {
	return Point{X: x, Y: y}
}

// PromotePt converts an integral image.Point to Point.
func PromotePt(p image.Point) Point {
	return Point{X: float64(p.X), Y: float64(p.Y)}
}

func (p Point) String() string {
	return fmt.Sprintf("%fy%f", p.X, p.Y)
}

// Add returns the vector p+q.
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
// Use 1/k for division.
func (p Point) Mul(k float64) Point {
	return Point{p.X * k, p.Y * k}
}

// Pmul returns the result of per-element multiplication of vectors p and q.
func (p Point) Pmul(q Point) Point {
	return Point{p.X * q.X, p.Y * q.Y}
}

func (p Point) Abs() Point {
	return Point{math.Abs(p.X), math.Abs(p.Y)}
}

func (p Point) Length() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// Dot returns the vector p•q.
func (p Point) Dot(q Point) float64 {
	return p.X*q.X + p.Y*q.Y
}

// Cross returns the vector p×q.
func (p Point) Cross(q Point) Point3d {
	return Pt3(p.X, p.Y, 1).Cross(Pt3(q.X, q.Y, 1))
}

// Triple returns the vector p•(q×d).
func (p Point) Triple(q, d Point) float64 {
	return p.To3().Dot(q.Cross(d))
}

// In reports whether p is in r.
func (p Point) In(r Rectangle) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

// Floor returns the greatest integer vector less than or equal to p.
func (p Point) Floor() Point {
	return Pt(math.Floor(p.X), math.Floor(p.Y))
}

// Ceil returns the least integer vector greater than or equal to x.
func (p Point) Ceil() Point {
	return Pt(math.Ceil(p.X), math.Ceil(p.Y))
}

// Mix linearly interpolates between points p and b by the parameter t.
func (p Point) Mix(b, t Point) Point {
	x := (1-t.X)*p.X + t.X*b.X
	y := (1-t.Y)*p.Y + t.Y*b.Y
	return Pt(x, y)
}

// Pmax returns the result of per-element maximum of vectors p and q.
func (p Point) Pmax(q Point) Point {
	return Point{max(p.X, q.X), max(p.Y, q.Y)}
}

// Pmin returns the result of per-element minimum of vectors p and q.
func (p Point) Pmin(q Point) Point {
	return Point{min(p.X, q.X), min(p.Y, q.Y)}
}

// Degrade returns the nearest integer vector as image.Point.
func (p Point) Degrade() image.Point {
	return image.Pt(int(math.Round(p.X)), int(math.Round(p.Y)))
}

// To3 promotes vector to three dimensions.
func (p Point) To3() Point3d { return Pt3(p.X, p.Y, 1) }

// Point3d is a X, Y and Z coordinate triple.
type Point3d struct {
	X, Y, Z float64
}

func Pt3(x, y, z float64) Point3d {
	return Point3d{X: x, Y: y, Z: z}
}

func (p Point3d) String() string {
	return fmt.Sprintf("%fy%fz%f", p.X, p.Y, p.Z)
}

// Add returns the vector p+q.
func (p Point3d) Add(q Point3d) Point3d {
	return Point3d{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

// Sub returns the vector p-q.
func (p Point3d) Sub(q Point3d) Point3d {
	return Point3d{p.X - q.X, p.Y - q.Y, p.Z - q.Z}
}

// Mul returns the vector p*k.
// Use 1/k for division.
func (p Point3d) Mul(k float64) Point3d {
	return Point3d{p.X * k, p.Y * k, p.Z * k}
}

// Dot returns the vector p⋅d.
func (p Point3d) Dot(q Point3d) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

// Cross returns the vector p×q.
func (p Point3d) Cross(q Point3d) Point3d {
	return Point3d{p.Y*q.Z - p.Z*q.Y, p.Z*q.X - p.X*q.Z, p.X*q.Y - p.Y*q.X}
}

// In reports whether p is in r.
func (p Point3d) In(r Rectangle) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

type Rectangle struct {
	Min, Max Point
}

func Rect(x0, y0, x1, y1 float64) Rectangle {
	return Rectangle{Min: Pt(x0, y0), Max: Pt(x1, y1)}
}

func PromoteRect(p image.Rectangle) Rectangle {
	return Rectangle{Min: PromotePt(p.Min), Max: PromotePt(p.Max)}
}

// Dx returns r's width.
func (r Rectangle) Dx() float64 {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectangle) Dy() float64 {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r Rectangle) Size() Point {
	return Point{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Add returns the rectangle r translated by p.
func (r Rectangle) Add(p Point) Rectangle {
	return Rectangle{
		Point{r.Min.X + p.X, r.Min.Y + p.Y},
		Point{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r Rectangle) Sub(p Point) Rectangle {
	return Rectangle{
		Point{r.Min.X - p.X, r.Min.Y - p.Y},
		Point{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r Rectangle) Inset(n float64) Rectangle {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Rectangle) Intersect(s Rectangle) Rectangle {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	// Letting r0 and s0 be the values of r and s at the time that the method
	// is called, this next line is equivalent to:
	//
	// if max(r0.Min.X, s0.Min.X) >= min(r0.Max.X, s0.Max.X) || likewiseForY { etc }
	if r.Empty() {
		return Rectangle{}
	}
	return r
}

// Union returns the smallest rectangle that contains both r and s.
func (r Rectangle) Union(s Rectangle) Rectangle {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

// Empty reports whether the rectangle contains no points.
func (r Rectangle) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r Rectangle) Eq(s Rectangle) bool {
	return r == s || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rectangle) Overlaps(s Rectangle) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r Rectangle) In(s Rectangle) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

// Canon returns the canonical version of r. The returned rectangle has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (r Rectangle) Canon() Rectangle {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// Center returns the geometric center of r.
func (r Rectangle) Center() Point {
	return Pt((r.Min.X+r.Max.X)/2, (r.Min.Y+r.Max.Y)/2)
}

// Distance evaluates signed distance function of the rectangle at the given point.
// If signed distance is negative, then the point is inside the rectangle.
// If it is zero, then it is on its boundary.
func (r Rectangle) Distance(at Point) float64 {
	// var x, y float64
	r = r.Sub(at)
	d := Pt(math.Abs(r.Min.X+r.Max.X)-r.Max.X+r.Min.X, math.Abs(r.Min.Y+r.Max.Y)-r.Max.Y+r.Min.Y)
	if d.X < 0 {
		d.X = 0
	}
	if d.Y < 0 {
		d.Y = 0
	}
	return d.Mul(0.5).Length()
}

// Degrade returns the nearest image.Rectangle.
func (r Rectangle) Degrade() image.Rectangle {
	return image.Rectangle{Min: r.Min.Degrade(), Max: r.Max.Degrade()}
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
