package tmx

import "fmt"

// Point descibes a location in 2D space.
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

// Size decsribes dimensions in 2D space.
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

// Rect describes a location and size in 2D space.
type Rect struct {
	// Point is the location of the rectangle.
	Point
	// Size is the dimensions of the rectangle.
	Size
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
