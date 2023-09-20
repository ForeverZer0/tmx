package tmx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// WangSet defines a list of colors and any number of Wang tiles using these colors.
type WangSet struct {
	// Name is the user-defined name of the Wang set.
	Name string `xml:"name,attr"`
	// Class is the user-defined class of the Wang set.
	Class string `xml:"class,attr,omitempty"`
	// Tile is the tile ID of the tile representing the Wang set.
	Tile TileID `xml:"tile,attr"`
	// Type indicates the behavior of terrain generation.
	Type WangType `xml:"type,attr"`
	// Colors is a collection of up to 254 colors used by the Wang set.
	Colors []WangColor `xml:"wangcolor"`
	// Tiles is a collection of tiles used by the Wang set.
	Tiles []WangTile `xml:"wangtile"`
	// Properties contain arbitrary key-value pairs of data to associate with the object.
	Properties `xml:"properties"`
}

func (w *WangSet) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "name":
			w.Name = attr.Value
		case "class":
			w.Class = attr.Value
		case "tile":
			var id TileID
			if err := id.UnmarshalText([]byte(attr.Value)); err != nil {
				return err
			} else {
				w.Tile = id
			}		
		case "type":
			if value, err := parseWangType(attr.Value); err != nil {
				return err
			} else {
				w.Type = value
			}
		default:
			logAttr(attr.Name.Local, start.Name.Local)
		}
	}	

	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		}

		if child, ok := token.(xml.StartElement); ok {
			switch child.Name.Local {
			case "wangcolor":
				var color WangColor
				if err := d.DecodeElement(&color, &child); err != nil {
					return err
				} else {
					w.Colors = append(w.Colors, color)
				}
			case "wangtile":
				var tile WangTile
				if err := d.DecodeElement(&tile, &child); err != nil {
					return err
				} else {
					w.Tiles = append(w.Tiles, tile)
				}
			case "properties":
				props := make(Properties)
				if err := d.DecodeElement(&props, &child); err != nil {
					return err
				} else {
					w.Properties = props
				}
			default:
				logElem(child.Name.Local, start.Name.Local)
			}
		}

		token, err = d.Token()
	}

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (w *WangSet) UnmarshalJSON(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	token, err := d.Token()
	if err != nil {
		return err
	} else if token != json.Delim('{') {
		return errors.New("expected JSON object")
	}

	for {
		if token, err = d.Token(); err != nil {
			return err
		} else if token == json.Delim('}') {
			break
		}

		name := token.(string)
		switch name {
		case "colors", "wangtiles", "properties":
		default:
			if token, err = d.Token(); err != nil {
				return err
			}
		}

		switch name {
		case "name":
			w.Name = token.(string)
		case "class":
			w.Class = token.(string)
		case "tile":
			w.Tile = TileID(token.(float64))
		case "type":
			if value, err := parseWangType(token.(string)); err != nil {
				return err
			} else {
				w.Type = value
			}
		case "colors":
			var colors []WangColor
			if err = d.Decode(&colors); err != nil {
				return err
			}
			w.Colors = colors
		case "wangtiles":
			var tiles []WangTile
			if err = d.Decode(&tiles); err != nil {
				return err
			}
			w.Tiles = tiles
		case "properties":
			props := make(Properties)
			if err := d.Decode(&props); err != nil {
				return err
			}
			w.Properties = props
		default:
			logProp(name, "wangset")
		}
	}
	return nil
}

// WangColor is a color that can be used to define the corner and/or edge of a Wang tile.
type WangColor struct {
	// Name is the user-defined name of the Wang color.
	Name string
	// Class is the user-defined class of the Wang color.
	Class string
	// Color is the RGB color used to represent the Wang color.
	Color Color
	// Tile is the tile ID of the tile representing the Wang color.
	Tile TileID
	// Probability is the relative probability that this color is chosen over others
	// in case of multiple options (defaults to 0).
	Probability float64
	// Properties contain arbitrary key-value pairs of data to associate with the object.
	Properties
}

func (w *WangColor) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "name":
			w.Name = attr.Value
		case "class":
			w.Class = attr.Value
		case "color":
			if color, err := ParseColor(attr.Value); err != nil {
				return err
			} else {
				w.Color = color
			}
		case "tile":
			var id TileID
			if err := id.UnmarshalText([]byte(attr.Value)); err != nil {
				return err
			} else {
				w.Tile = id
			}
		case "probability":
			if value, err := strconv.ParseFloat(attr.Value, 64); err != nil {
				return err
			} else {
				w.Probability = value
			}
		default:
			logAttr(attr.Name.Local, start.Name.Local)
		}
	}

	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		}

		if child, ok := token.(xml.StartElement); ok {
			switch child.Name.Local {
			case "properties":
				props := make(Properties)
				if err := d.DecodeElement(&props, &child); err != nil {
					return err
				} else {
					w.Properties = props
				}
			default:
				logElem(child.Name.Local, start.Name.Local)
			}
		}
		token, err = d.Token()
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (w *WangColor) UnmarshalJSON(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	token, err := d.Token()
	if err != nil {
		return err
	} else if token != json.Delim('{') {
		return errors.New("expected JSON object")
	}

	for {
		if token, err = d.Token(); err != nil {
			return err
		} else if token == json.Delim('}') {
			break
		}

		name := token.(string)
		if name != "properties" {
			if token, err = d.Token(); err != nil {
				return err
			}
		}

		switch name {
		case "name":
			w.Name = token.(string)
		case "class":
			w.Name = token.(string)
		case "color":
			if color, err := ParseColor(token.(string)); err != nil {
				return err
			} else {
				w.Color = color
			}
		case "tile":
			// TODO: Check if valid before cast
			w.Tile = TileID(token.(float64))
		case "probability":
			w.Probability = token.(float64)
		case "properties":
			props := make(Properties)
			if err := d.Decode(&props); err != nil {
				return err
			}
			w.Properties = props
		default:
			logProp(name, "wangcolor")
		}
	}
	
	return nil
}

type WangTile struct {
	// Tile is the local tile ID used by the Wang tile.
	Tile TileID `json:"tileid"`
	// WangID is a list of indices (0-254) referring to the Wang colors in the Wang set in
	// the order: top, top-right, right, bottom-right, bottom, bottom-left, left, top-left.
	//
	// Index 0 means unset and index 1 refers to the first Wang color.
	WangID [8]uint8 `json:"wangid"`
	// Deprecated: Defaults to false and is now defined in Transformations.
	HFlip bool `json:"hflip"`
	// Deprecated: Defaults to false and is now defined in Transformations.
	VFlip bool `json:"vflip"`
	// Deprecated: Defaults to false and is now defined in Transformations.
	DFlip bool `json:"dflip"`
}

func (w *WangTile) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "tileid":
			var id TileID
			if err := id.UnmarshalText([]byte(attr.Value)); err != nil {
				return err
			} else {
				w.Tile = id
			}
		case "wangid":
			fields := strings.Split(attr.Value, ",")
			if len(fields) > len(w.WangID) {
				return errors.New("expected array of 8 elements or less in WangTile")
			}
			for i := 0; i < len(fields); i++ {
				if value, err := strconv.ParseUint(fields[i], 10, 8); err != nil {
					return err
				} else {
					w.WangID[i] = uint8(value)
				}
			}
		case "hflip":
			log.Println("WangTile: hflip is deprecated, use Transformations")
			if value, err := strconv.ParseBool(attr.Value); err != nil {
				return err
			} else {
				w.HFlip = value
			}
		case "vflip":
			log.Println("WangTile: vflip is deprecated, use Transformations")
			if value, err := strconv.ParseBool(attr.Value); err != nil {
				return err
			} else {
				w.VFlip = value
			}
		case "dlip":
			log.Println("WangTile: dflip is deprecated, use Transformations")
			if value, err := strconv.ParseBool(attr.Value); err != nil {
				return err
			} else {
				w.DFlip = value
			}
		default:
			logAttr(attr.Name.Local, start.Name.Local)
		}
	}

	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		}
		if child, ok := token.(xml.StartElement); ok {
			logElem(child.Name.Local, start.Name.Local)
		}
		token, err = d.Token()
	}

	return nil
}


// WangType describes the behavior of terrain generation.
type WangType int

const (
	// WangTypeCorner is a WangType of type Corner.
	WangTypeCorner WangType = iota
	// WangTypeEdge is a WangType of type Edge.
	WangTypeEdge
	// WangTypeMixed is a WangType of type Mixed.
	WangTypeMixed
)

const _WangTypeName = "corneredgemixed"

var _WangTypeMap = map[WangType]string{
	WangTypeCorner: _WangTypeName[0:6],
	WangTypeEdge:   _WangTypeName[6:10],
	WangTypeMixed:  _WangTypeName[10:15],
}

// String implements the Stringer interface.
func (e WangType) String() string {
	if str, ok := _WangTypeMap[e]; ok {
		return str
	}
	return fmt.Sprintf("WangType(%d)", e)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (e WangType) IsValid() bool {
	_, ok := _WangTypeMap[e]
	return ok
}

var _WangTypeValue = map[string]WangType{
	_WangTypeName[0:6]:   WangTypeCorner,
	_WangTypeName[6:10]:  WangTypeEdge,
	_WangTypeName[10:15]: WangTypeMixed,
}

// parseWangType attempts to convert a string to a WangType.
func parseWangType(name string) (WangType, error) {
	if x, ok := _WangTypeValue[name]; ok {
		return x, nil
	}
	return WangType(0), errInvalidEnum("WangType", name)
}

// MarshalText implements the text marshaller method.
func (e WangType) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (e *WangType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseWangType(name)
	if err != nil {
		return err
	}
	*e = tmp
	return nil
}

// vim: ts=4
