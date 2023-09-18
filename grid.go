package tmx

import "fmt"

// Grid describes grid settings used for tiles in a tileset.
type Grid struct {
	// Size is the dimensions of a tile cell in the grid.
	Size
	// Orientation indicates the orientation of the grid.
	Orientation Orientation `json:"orientation" xml:"orientation,attr"`
}

// IsEmpty indicates if the grid has any defined values or is the default/empty "zero" value.
func (g Grid) IsEmpty() bool {
	return g.Width == 0 && g.Height == 0 && g.Orientation == 0
}

// String implements the Stringer interface.
func (g Grid) String() string {
	if g.IsEmpty() {
		return ""
	}
	return fmt.Sprintf("%dx%d (%s)", g.Width, g.Height, g.Orientation.String())
}

// vim: ts=4
