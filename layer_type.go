package tmx

import "fmt"

// ENUM(none, tilelayer, objectgroup, imagelayer, group, all)
type LayerType byte

const (
	// LayerNone is a LayerType of type None.
	LayerNone LayerType = 0x00
	// LayerTile is a LayerType of type Tile.
	LayerTile LayerType = 0x01
	// LayerImage is a LayerType of type Image.
	LayerImage LayerType = 0x02
	// LayerObject is a LayerType of type Object.
	LayerObject LayerType = 0x04
	// LayerGroup is a LayerType of type Group.
	LayerGroup LayerType = 0x08
	// LayerAll is a LayerType of type All.
	LayerAll LayerType = 0xFF
)

const _LayerTypeName = "nonetilelayerobjectgroupimagelayergroupall"

var _LayerTypeMap = map[LayerType]string{
	LayerNone:   _LayerTypeName[0:4],
	LayerTile:   _LayerTypeName[4:13],
	LayerObject: _LayerTypeName[13:24],
	LayerImage:  _LayerTypeName[24:34],
	LayerGroup:  _LayerTypeName[34:39],
	LayerAll:    _LayerTypeName[39:42],
}

// String implements the Stringer interface.
func (x LayerType) String() string {
	if str, ok := _LayerTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("LayerType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x LayerType) IsValid() bool {
	_, ok := _LayerTypeMap[x]
	return ok
}

var _LayerTypeValue = map[string]LayerType{
	_LayerTypeName[0:4]:   LayerNone,
	_LayerTypeName[4:13]:  LayerTile,
	_LayerTypeName[13:24]: LayerObject,
	_LayerTypeName[24:34]: LayerImage,
	_LayerTypeName[34:39]: LayerGroup,
	_LayerTypeName[39:42]: LayerAll,
}

// parseLayerType attempts to convert a string to a LayerType.
func parseLayerType(name string) (LayerType, error) {
	if x, ok := _LayerTypeValue[name]; ok {
		return x, nil
	}
	return LayerType(0), errInvalidEnum("LayerType", name)
}

// MarshalText implements the text marshaller method.
func (x LayerType) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *LayerType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseLayerType(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
