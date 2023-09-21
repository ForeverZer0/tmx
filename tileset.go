package tmx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strconv"
)

// Tileset is a the core tileset implementation that contains no map-specific fields, such
// FirstGID, parent Map, etc.
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
				tile.cache = ts.cache
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
				// TODO: Parse
				log.Println("terraintypes are no longer supported, use wangsets instead")
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

	// TODO
	for i, tile := range ts.Tiles {

		if tile.Width == 0 {
			ts.Tiles[i].Width = ts.TileSize.Width
		}
		if tile.Height == 0 {
			ts.Tiles[i].Height = ts.TileSize.Height
		}

		x := int(tile.ID) % ts.Columns
		y := int(tile.ID) / ts.Columns
		ts.Tiles[i].Point = Point{X: x, Y: y}
	}

	return nil
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
		if impl, err := OpenTileset(source, ts.cache); err == nil {
			ts.Tileset = impl
		} else {
			return err
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
		if tileset, err := OpenTileset(temp.Source, ts.cache); err != nil {
			return err
		} else {
			ts.Tileset = tileset
		}
	}

	return nil
}

func (ts *Tileset) UnmarshalJSON(data []byte) error {
	// TODO
	type jsonTileset struct {
		FirstGID         TileID           `json:"firstgid"`
		Source           string           `json:"source"`
		BackgroundColor  Color            `json:"backgroundcolor"`
		Class            string           `json:"class"`
		Columns          int              `json:"columns"`
		FillMode         FillMode         `json:"fill_mode"`
		Grid             *Grid            `json:"grid"`
		Image            string           `json:"image"`
		ImageWidth       int              `json:"imagewidth"`
		ImageHeight      int              `json:"imageheight"`
		Margin           int              `json:"margin"`
		Name             string           `json:"name"`
		ObjectAlignment  Align            `json:"objectalignment"`
		Properties       Properties       `json:"properties"`
		Spacing          int              `json:"spacing"`
		Terrains         []int            `json:"terrains"`
		TileCount        int              `json:"tilecount"`
		TiledVersion     string           `json:"tiledversion"`
		TileHeight       int              `json:"tileheight"`
		TileWidth        int              `json:"tilewidth"`
		TileRenderSize   TileRender       `json:"tilerendersize"`
		Tiles            []Tile           `json:"tiles"`
		Wangsets         []WangSet        `json:"wangsets"`
		Version          string           `json:"version"`
		TransparentColor Color            `json:"transparentcolor"`
		Transformations  *Transformations `json:"transformations"`
		Type             string           `json:"type"`
		TileOffset       Point            `json:"tileoffset"`
	}

	var temp jsonTileset
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	ts.Name = temp.Name
	ts.Class = temp.Class
	ts.Version = temp.Version
	ts.TiledVersion = temp.TiledVersion
	ts.Columns = temp.Columns
	ts.FillMode = temp.FillMode
	ts.Grid = temp.Grid
	ts.Margin = temp.Margin
	ts.Spacing = temp.Spacing
	ts.Tiles = temp.Tiles
	ts.ObjectAlign = temp.ObjectAlignment
	ts.Count = temp.TileCount
	ts.Offset = temp.TileOffset
	ts.Properties = temp.Properties
	ts.RenderSize = temp.TileRenderSize
	ts.TileSize = Size{Width: temp.TileWidth, Height: temp.TileHeight}
	ts.WangSets = temp.Wangsets
	ts.Transforms = temp.Transformations
	ts.BackgroundColor = temp.BackgroundColor

	if temp.Source != "" {
		ts.Image = &Image{
			Source: temp.Source,
			Size: Size{
				Width:  temp.ImageWidth,
				Height: temp.ImageHeight,
			},
			Transparency: temp.TransparentColor,
		}
	}

	// Terrains         []int            `json:"terrains"`

	return nil
}

// OpenMap reads a tilemap from a file, automatically detecting it format.
//
// An optional cache can be supplied that maintains references to tilesets and
// templates to prevent frequent re-processing of them.
func OpenTileset(path string, cache *Cache) (*Tileset, error) {
	return OpenTilesetFormat(path, detectFileExt(path), cache)
}

// OpenMapFormat reads a tilemap from a file, using the specified format.
//
// An optional cache can be supplied that maintains references to tilesets and
// templates to prevent frequent re-processing of them.
func OpenTilesetFormat(path string, format Format, cache *Cache) (*Tileset, error) {
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
	// tileset.cache = cache

	if err := ReadTilesetFormat(reader, format, &tileset); err != nil {
		return nil, err
	}

	if cache != nil {
		cache.AddTileset(abs, &tileset)
	}
	return &tileset, nil
}

// ReadMap reads a tilemap from the current position in the reader.
func ReadTileset(r io.ReadSeeker, tileset *Tileset) error {
	return ReadTilesetFormat(r, detectReader(r), tileset)
}

// ReadMapFormat reads a tilemap from the current position in the reader using
// the specified format.
func ReadTilesetFormat(r io.Reader, format Format, tileset *Tileset) error {
	switch format {
	case FormatXML:
		d := xml.NewDecoder(r)
		if err := d.Decode(tileset); err != nil {
			return err
		}
	case FormatJSON:
		d := json.NewDecoder(r)
		if err := d.Decode(tileset); err != nil {
			return err
		}
	default:
		return errInvalidEnum("Format", fmt.Sprintf("Format(%d)", format))
	}

	return nil
}

// vim: ts=4
