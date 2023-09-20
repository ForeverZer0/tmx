package tmx

import "fmt"

type TileRender int

const (
	// RenderTile is a TileRender of type Tile.
	RenderTile TileRender = iota
	// RenderGrid is a TileRender of type Grid.
	RenderGrid
)

const _TileRenderName = "tilegrid"

var _TileRenderMap = map[TileRender]string{
	RenderTile: _TileRenderName[0:4],
	RenderGrid: _TileRenderName[4:8],
}

// String implements the Stringer interface.
func (x TileRender) String() string {
	if str, ok := _TileRenderMap[x]; ok {
		return str
	}
	return fmt.Sprintf("TileRender(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x TileRender) IsValid() bool {
	_, ok := _TileRenderMap[x]
	return ok
}

var _TileRenderValue = map[string]TileRender{
	_TileRenderName[0:4]: RenderTile,
	_TileRenderName[4:8]: RenderGrid,
}

// parseTileRender attempts to convert a string to a TileRender.
func parseTileRender(name string) (TileRender, error) {
	if x, ok := _TileRenderValue[name]; ok {
		return x, nil
	}
	return TileRender(0), errInvalidEnum("TileRender", name)
}

// MarshalText implements the text marshaller method.
func (x TileRender) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *TileRender) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseTileRender(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// vim: ts=4
