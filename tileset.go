package tmx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
)

// BottomLeftOrigin is a global configuration specifically for generating UV coordinates
// in Tile objects.
//
// By default, the origin (0,0) is the top-left corner, with y increasing as it moves down,
// which is the most common for 2D graphics. When set to true, UV coordinates will be
// calculated using the origin at the bottom-left, with y increasing as it moves up.
var BottomLeftOrigin bool

// Tileset is a the core tileset implementation, unbound to any specific map.
type Tileset struct {
	// Source is the location of an external file that contains tileset definition, or empty string
	// for the case of an embedded tileset.
	Source string
	// Version is the TMX format version.
	Version string
	// TiledVersion is the Tiled version used to save the file. May be a date (for snapshot
	// builds).
	TiledVersion string
	// Name is the user-defined name of this tileset.
	Name string
	// Class is the user-defined he class of this tileset.
	Class string
	// TileSize is the (maximum) dimensions of tiles in the tileset. Irrelevant for image
	// collection tilesets.
	TileSize Size
	// Spacing is the number of pixels between tiles in the tileset.
	Spacing int
	// Margin is the spacing around tiles in the tileset (applies to the tileset image).
	Margin int
	// Count is the number of tiles in this tileset. Note that there can be tiles with a higher
	// ID than the tile count, in case the tileset is an image collection from which tiles have
	// been removed.
	Count int
	// Columns is the number of tile columns in the tileset.
	//
	// For image collection tilesets it is editable and is used when displaying the tileset.
	Columns int
	// ObjectAlign controls the alignment for tile objects.
	//
	// For compatibility reasons, default value is unspecified. When unspecified, tile objects
	// use BottomLeft in orthogonal mode and Bottom in isometric mode.
	ObjectAlign Align
	// RenderSize is the size to use when rendering tiles from this tileset on a tile layer.
	RenderSize TileRender
	// FillMode is the fill mode to use when rendering tiles from this tileset. This field is
	// only relevant when the tiles are not rendered at their native size or in combination of
	// the RenderSize being set to RenderGrid.
	FillMode FillMode
	// Offset is a pixel offset for drawing tiles.
	Offset Point
	// Image is the source graphics of the tileset. When nil, and each tile defines its own
	// image.
	Image *Image
	// Tiles contains the definitions for all tiles that tileset is made from.
	Tiles []Tile
	// WangSets is a collection of WangSet objects.
	WangSets []WangSet
	// Properties contain arbitrary key-value pairs of data to associate with the object.
	Properties
	// Grid is used for isometric orientation, and determines how tile overlays for terrain and
	// collision information are rendered.
	Grid *Grid
	// Transforms is an optional field that describes which transformations can be applied to the
	// tiles (e.g. to extend a Wang set by transforming existing tiles).
	Transforms *Transformations
	// BackgroundColor is an optional color that is used when displaying the tileset in the
	// Tiled editor. Defaults to full transparency, and is typically of little relevance in
	// regards to tilemap rendering.
	BackgroundColor Color
	// cache is a resource cache that maintains references to shared objects.
	cache *Cache
}

// MapTileset describes the source tiles/graphics used by tilemaps.
type MapTileset struct {
	// FirstGID is the first global tile ID of this tileset (this global ID maps to the first
	// tile in this tileset).
	FirstGID TileID
	// Map is the parent tilemap this tileset is being used in.
	Map *Map
	// Tileset is the actual tileset implementation, and is unbound by the map-specific fields,
	// allowing it to be cached and reused with different maps.
	*Tileset
	// cache is a resource cache that maintains references to shared objects.
	cache *Cache
}

// String implements the Stringer interface.
func (ts *Tileset) String() string {
	return fmt.Sprintf(`Tileset("%s")`, ts.Name)
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (ts *Tileset) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "version":
			ts.Version = attr.Value
		case "tiledversion":
			ts.TiledVersion = attr.Value
		case "name":
			ts.Name = attr.Value
		case "class":
			ts.Class = attr.Value
		case "tilewidth":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				ts.TileSize.Width = value
			} else {
				return err
			}
		case "tileheight":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				ts.TileSize.Height = value
			} else {
				return err
			}
		case "spacing":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				ts.Spacing = value
			} else {
				return err
			}
		case "margin":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				ts.Margin = value
			} else {
				return err
			}
		case "tilecount":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				ts.Count = value
			} else {
				return err
			}
		case "columns":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				ts.Columns = value
			} else {
				return err
			}
		case "objectalignment":
			if value, err := parseAlign(attr.Value); err != nil {
				return err
			} else {
				ts.ObjectAlign = value
			}
		case "tilerendersize":
			if value, err := parseTileRender(attr.Value); err != nil {
				return err
			} else {
				ts.RenderSize = value
			}
		case "fillmode":
			if value, err := parseFillMode(attr.Value); err != nil {
				return err
			} else {
				ts.FillMode = value
			}
		case "backgroundcolor":
			if value, err := ParseColor(attr.Value); err != nil {
				// Log this error, but do not terminate parsing over it.
				log.Println(err)
			} else {
				ts.BackgroundColor = value
			}
		case "source":
			if ts.Source == "" {
				ts.Source = attr.Value
			}
		case "firstgid":
			// Skip
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
			case "tile":
				var tile Tile
				tile.Tileset = ts
				if err := tile.UnmarshalXML(d, child); err != nil {
					return err
				}
				ts.Tiles = append(ts.Tiles, tile)
			case "properties":
				ts.Properties = make(Properties)
				if err := ts.Properties.UnmarshalXML(d, child); err != nil {
					return err
				}
			case "image":
				var image Image
				if err := image.UnmarshalXML(d, child); err != nil {
					return err
				}
				ts.Image = &image
			case "tileoffset":
				if err := d.DecodeElement(&ts.Offset, &child); err != nil {
					return err
				}
			case "grid":
				var grid Grid
				if err := d.DecodeElement(&grid, &child); err != nil {
					return err
				}
				ts.Grid = &grid
			case "terraintypes":
				logTerrain()
			case "wangsets":
				type wangsets struct {
					Values []WangSet `xml:"wangset"`
				}
				var wang wangsets
				if err := d.DecodeElement(&wang, &child); err != nil {
					return err
				}
				ts.WangSets = wang.Values
			case "transformations":
				var trans Transformations
				if d.DecodeElement(&trans, &child); err != nil {
					return err
				}
				ts.Transforms = &trans
			default:
				logElem(child.Name.Local, start.Name.Local)
			}
		}

		token, err = d.Token()
	}

	ts.postProcess()
	return nil
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (ts *MapTileset) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var source string
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "firstgid":
			if value, err := strconv.ParseUint(attr.Value, 10, 32); err == nil {
				ts.FirstGID = TileID(value)
			} else {
				return err
			}
		case "source":
			source = attr.Value
		}

		if ts.FirstGID > 0 && source != "" {
			break
		}
	}

	if source == "" {
		// Embedded tileset
		var impl Tileset
		if err := impl.UnmarshalXML(d, start); err != nil {
			return err
		}
		ts.Tileset = &impl
	} else {
		if impl, err := OpenTileset(source, DetectExt(source), ts.cache); err == nil {
			ts.Tileset = impl
		} else {
			return err
		}
	}

	// Ensure the element is fully consumed
	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		}
		if child, ok := token.(xml.StartElement); ok {
			logElem(child.Name.Local, start.Name.Local)
		}
	}

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ts *MapTileset) UnmarshalJSON(data []byte) error {
	type mapTileset struct {
		FirstGID TileID `json:"firstgid"`
		Source   string `json:"source"`
	}

	var temp mapTileset
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.Source == "" {
		var tileset Tileset
		if err := json.Unmarshal(data, &tileset); err != nil {
			return err
		}
		ts.Tileset = &tileset
	} else {
		if tileset, err := OpenTileset(temp.Source, DetectExt(temp.Source), ts.cache); err != nil {
			return err
		} else {
			ts.Tileset = tileset
		}
	}

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ts *Tileset) UnmarshalJSON(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	token, err := d.Token()
	if err != nil {
		return err
	} else if token != json.Delim('{') {
		return ErrExpectedObject
	}

	for {
		if token, err = d.Token(); err != nil {
			return err
		} else if token == json.Delim('}') {
			break
		}

		name := token.(string)
		switch name {
		case "backgroundcolor":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if ts.BackgroundColor, err = ParseColor(str); err != nil {
				return err
			}
		case "class":
			if ts.Class, err = jsonProp[string](d); err != nil {
				return err
			}
		case "columns":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				ts.Columns = int(value)
			}
		case "fill_mode":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if ts.FillMode, err = parseFillMode(str); err != nil {
				return err
			}
		case "grid":
			var grid Grid
			if err = d.Decode(&grid); err != nil {
				return err
			}
			ts.Grid = &grid
		case "image":
			if ts.Image == nil {
				ts.Image = &Image{}
			}
			if ts.Image.Source, err = jsonProp[string](d); err != nil {
				return err
			}
		case "imagewidth":
			if ts.Image == nil {
				ts.Image = &Image{}
			}
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				ts.Image.Size.Width = int(value)
			}
		case "imageheight":
			if ts.Image == nil {
				ts.Image = &Image{}
			}
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				ts.Image.Size.Height = int(value)
			}
		case "margin":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				ts.Margin = int(value)
			}
		case "name":
			if ts.Name, err = jsonProp[string](d); err != nil {
				return err
			}
		case "objectalignment":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if ts.ObjectAlign, err = parseAlign(str); err != nil {
				return err
			}
		case "properties":
			props := make(Properties)
			if err = d.Decode(&props); err != nil {
				return err
			}
			ts.Properties = props
		case "spacing":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				ts.Spacing = int(value)
			}
		case "tilecount":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				ts.Count = int(value)
			}
		case "tiledversion":
			if ts.TiledVersion, err = jsonProp[string](d); err != nil {
				return err
			}
		case "tileheight":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				ts.TileSize.Height = int(value)
			}
		case "tilewidth":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				ts.TileSize.Width = int(value)
			}
		case "tilerendersize":
			if value, err := jsonProp[string](d); err != nil {
				return err
			} else if ts.RenderSize, err = parseTileRender(value); err != nil {
				return err
			}
		case "tiles":
			if token, err = d.Token(); err != nil {
				return err
			} else if token != json.Delim('[') {
				return ErrExpectedArray
			}
			for d.More() {
				tile := Tile{Tileset: ts}
				if err = d.Decode(&tile); err != nil {
					return err
				}
				ts.Tiles = append(ts.Tiles, tile)
			}
			// Consume the closing ']'
			if token, err = d.Token(); err != nil {
				return err
			}
		case "wangsets":
			if err = d.Decode(&ts.WangSets); err != nil {
				return err
			}
		case "version":
			if ts.Version, err = jsonProp[string](d); err != nil {
				return err
			}
		case "transparentcolor":
			if ts.Image == nil {
				ts.Image = &Image{}
			}
			if value, err := jsonProp[string](d); err != nil {
				return err
			} else if ts.Image.Transparency, err = ParseColor(value); err != nil {
				return err
			}
		case "transformations":
			var transforms Transformations
			if err = d.Decode(&transforms); err != nil {
				return err
			}
			ts.Transforms = &transforms
		case "tileoffset":
			if err = d.Decode(&ts.Offset); err != nil {
				return err
			}
		case "terrains":
			logTerrain()
			jsonSkip(d)
		case "firstgid", "source", "type":
			jsonSkip(d)
		default:
			logProp(name, "tileset")
			jsonSkip(d)
		}
	}

	ts.postProcess()
	return nil
}

func (ts *Tileset) postProcess() {
	var cx, cy float32
	if ts.Image != nil && ts.Image.Width > 0 && ts.Image.Height > 0 {
		cx = float32(ts.TileSize.Width) / float32(ts.Image.Width)
		cy = float32(ts.TileSize.Height) / float32(ts.Image.Height)
	}

	for i := range ts.Tiles {
		tile := &ts.Tiles[i]

		if tile.Width == 0 {
			tile.Width = ts.TileSize.Width
		}
		if tile.Height == 0 {
			tile.Height = ts.TileSize.Height
		}

		if tile.Image != nil {
			// TODO
			tile.UV0 = Vec2{0.0, 0.0}
			tile.UV1 = Vec2{1.0, 1.0}
		} else {
			tile.Point = Point{
				X: int(tile.ID) % ts.Columns,
				Y: int(tile.ID) / ts.Columns,
			}

			tile.UV0.X = float32(tile.X) * cx
			tile.UV1.X = min(float32(tile.X + 1) * cx, 1.0)

			if BottomLeftOrigin {
				tile.UV0.Y = 1.0 - max(float32(tile.Bottom()) * cy, 0.0)
				tile.UV1.Y = 1.0 - min(float32(tile.Bottom() - 1) * cy, 1.0)
			} else {
				tile.UV0.Y = float32(tile.Y) * cy
				tile.UV1.Y = min(float32(tile.Y + 1) * cy, 1.0)
			}
		}
	}
}



// OpenTileset reads a tileset from a file, using the specified format.
//
// An optional cache can be supplied that maintains references to tilesets and
// templates to prevent frequent re-processing of them.
func OpenTileset(path string, format Format, cache *Cache) (*Tileset, error) {
	var abs string
	var err error
	if abs, err = FindPath(path); err != nil {
		return nil, err
	}

	// Check cache
	if cache != nil {
		if tileset, ok := cache.Tileset(abs); ok {
			return tileset, nil
		}
	}

	reader, _, err := getStream(abs)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	IncludePaths = append(IncludePaths, filepath.Dir(abs))
	defer func() { IncludePaths = IncludePaths[:len(IncludePaths)-1] }()

	var tileset Tileset
	tileset.Source = abs
	tileset.cache = cache

	if err := Decode(reader, format, &tileset); err != nil {
		return nil, err
	}

	if cache != nil {
		cache.AddTileset(abs, &tileset)
	}
	return &tileset, nil
}

// vim: ts=4
