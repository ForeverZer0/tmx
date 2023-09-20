package tmx

import "strconv"

// TileID describes the the ID of a single tile on a map.
//
// The values is a bitfield with flags denoting orientation, rotation, etc, and as such these
// flags need masked before they can be used for indexing. See ClearMask for details.
type TileID uint32

const (
	// FlipH is a bitflag that indicates the tile is flipped on the horizotal axis.
	//
	//		var flipped bool = gid & FlipH != 0
	FlipH TileID = 0x80000000 // 32
	// FlipV is a bitflag that indicates the tile is flipped on the vertical axis.
	//
	//		var flipped bool = gid & FlipV != 0
	FlipV TileID = 0x40000000 // 31
	// FlipD is a bitflag that indicates the tile is diagonally.
	//
	//		var flipped bool = gid & FlipD != 0
	FlipD TileID = 0x20000000 // 30
	// RotateCCW is a bitflag that indicates the tile is rotate 60 degress counter-clockwise.
	// Only valid for hexagonal maps.
	//
	//		var rotatedCCW = gid & RotateCCW!= 0
	RotateCCW TileID = 0x10000000 // 29
	// RotateCW is a bitflag that indicates the tile is rotate 60 degress clockwise.
	// Only valid for hexagonal maps, otherwise it shares the same bit as FlipD.
	//
	//		var rotatedCW bool = gid & RotateCW!= 0
	RotateCW TileID = FlipD // 29
	// ClearMask is a bitflag that can be AND together with a TileID to remove all
	// flip/rotate flags and isolate the actual tile ID.
	//
	//		var clean TileID = gid & ClearMask
	ClearMask TileID = ^(FlipH | FlipV | FlipD | RotateCCW)
	// InvalidID is a strongly-typed value indicating the tile is invalid,
	// should be ignored, or a context-dependent "error" value.
	InvalidID TileID = 0xFFFFFFFF
)

// UnmarshalJSON implements the json.Unmarshaler interface.
func (id *TileID) UnmarshalJSON(data []byte) error {
	text := string(data)
	if value, err := strconv.ParseUint(text, 10, 32); err != nil {
		if text == "-1" {
			*id = InvalidID
			return nil
		}
		return err
	} else {
		*id = TileID(value)
	}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (id *TileID) UnmarshalText(text []byte) error {
	str := string(text)
	if value, err := strconv.ParseUint(str, 10, 32); err != nil {
		if str == "-1" {
			*id = InvalidID
			return nil
		}
		return err
	} else {
		*id = TileID(value)
	}
	return nil
}

// vim: ts=4
