package tmx

// Transformations describe which transformations can be applied to the tiles in
// tileset (e.g. to extend a Wang set by transforming existing tiles).
type Transformations struct {
	// HFlip indicates whether the tiles in this set can be flipped horizontally.
	HFlip bool `json:"hflip" xml:"hflip,attr"`
	// VFlip indicates whether the tiles in this set can be flipped vertically.
	VFlip bool `json:"vflip" xml:"vflip,attr"`
	// Rotate indicates whether the tiles in this set can be rotated in 90 degree increments.
	Rotate bool `json:"rotate" xml:"rotate,attr"`
	// PreferUntransformed indicates whether untransformed tiles remain preferred, otherwise
	// transformed tiles are used to produce more variations
	PreferUntransformed bool `json:"preferuntransformed" xml:"preferuntransformed,attr"`
}

// vim: ts=4
