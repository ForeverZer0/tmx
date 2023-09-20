package tmx

import "fmt"

// Point describes a location in 2D space.
type Point struct {
	// X is the location on the horizontal x-axis.
	X int `xml:"x,attr" json:"x"`
	// Y is the location on the vertical y-axis.
	Y int `xml:"y,attr" json:"y"`
}

// String implements the Stringer interface.
func (p Point) String() string {
	return fmt.Sprintf("<%d, %d>", p.X, p.Y)
}

// Size descibes dimensions in 2D space.
type Size struct {
	// Width is the dimension on the horizontal x-axis.
	Width int `xml:"width,attr" json:"width"`
	// Height is the dimension on the vetical y-axis.
	Height int `xml:"height,attr" json:"height"`
}

// String implements the Stringer interface.
func (s Size) String() string {
	return fmt.Sprintf("<%d, %d>", s.Width, s.Height)
}

// Area returns the total spacial area.
func (s Size) Area() int {
	return s.Width * s.Height
}

// Rect describes a location and size in 2D space.
type Rect struct {
	// Point is the location of the rectangle.
	Point
	// Size is the dimensions of the rectangle.
	Size
}

// Left returns the left edge of the rectangle.
func (r Rect) Left() int {
	return r.X
}

// Left returns the right edge of the rectangle.
func (r Rect) Right() int {
	return r.X + r.Width
}

// Left returns the top edge of the rectangle.
func (r Rect) Top() int {
	return r.Y
}

// Left returns the bottom edge of the rectangle.
func (r Rect) Bottom() int {
	return r.Y + r.Height
}

// TopLeft returns the point at the top-left corner of the rectangle.
func (r Rect) TopLeft() Point {
	return Point{X: r.X, Y: r.Y}
}

// TopRight returns the point at the top-right corner of the rectangle.
func (r Rect) TopRight() Point {
	return Point{X: r.X + r.Width, Y: r.Y}
}

// BottomLeft returns the point at the bottom-left corner of the rectangle.
func (r Rect) BottomLeft() Point {
	return Point{X: r.X, Y: r.Y + r.Height}
}

// BottomRight returns the point at the bottom-right corner of the rectangle.
func (r Rect) BottomRight() Point {
	return Point{X: r.X + r.Width, Y: r.Y + r.Height}
}

// String implements the Stringer interface.
func (r Rect) String() string {
	return fmt.Sprintf("<%d, %d, %d, %d>", r.X, r.Y, r.Width, r.Height)
}

// Vec2 describes a vector with two 32-bit float components.
type Vec2 struct {
	// X is the x-component of the vector.
	X float32 `xml:"x,attr" json:"x"`
	// Y is the y-component of the vector.
	Y float32 `xml:"y,attr" json:"y"`
}

// String implements the Stringer interface.
func (v Vec2) String() string {
	return fmt.Sprintf("<%f, %f>", v.X, v.Y)
}

// vim: ts=4
