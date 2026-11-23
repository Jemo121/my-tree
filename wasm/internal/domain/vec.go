package domain

import "math"

// Vec2 is a 2D vector or point.
type Vec2 struct{ X, Y float64 }

// Add returns v + u.
func (v Vec2) Add(u Vec2) Vec2 { return Vec2{v.X + u.X, v.Y + u.Y} }

// Sub returns v - u.
func (v Vec2) Sub(u Vec2) Vec2 { return Vec2{v.X - u.X, v.Y - u.Y} }

// Scale returns v * s.
func (v Vec2) Scale(s float64) Vec2 { return Vec2{v.X * s, v.Y * s} }

// Len returns the magnitude of v.
func (v Vec2) Len() float64 { return math.Sqrt(v.X*v.X + v.Y*v.Y) }

// Norm returns the unit vector of v. Returns zero vector if magnitude is zero.
func (v Vec2) Norm() Vec2 {
	l := v.Len()
	if l == 0 {
		return Vec2{}
	}
	return Vec2{v.X / l, v.Y / l}
}

// DistSq returns the squared distance between v and u.
func (v Vec2) DistSq(u Vec2) float64 {
	dx, dy := v.X-u.X, v.Y-u.Y
	return dx*dx + dy*dy
}

// Rect is an axis-aligned bounding box.
type Rect struct{ Min, Max Vec2 }

// Center returns the center of r.
func (r Rect) Center() Vec2 {
	return Vec2{(r.Min.X + r.Max.X) / 2, (r.Min.Y + r.Max.Y) / 2}
}

// Contains reports whether p lies inside r.
func (r Rect) Contains(p Vec2) bool {
	return p.X >= r.Min.X && p.X <= r.Max.X &&
		p.Y >= r.Min.Y && p.Y <= r.Max.Y
}

// Intersects reports whether r and s overlap.
func (r Rect) Intersects(s Rect) bool {
	return r.Min.X <= s.Max.X && r.Max.X >= s.Min.X &&
		r.Min.Y <= s.Max.Y && r.Max.Y >= s.Min.Y
}

// RectFromCenter returns a Rect centered at c with given half-width and half-height.
func RectFromCenter(c Vec2, hw, hh float64) Rect {
	return Rect{Min: Vec2{c.X - hw, c.Y - hh}, Max: Vec2{c.X + hw, c.Y + hh}}
}
