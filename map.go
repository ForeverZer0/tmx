package tmx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
)

// Map is the top-level object defining a complete tilemap, composed of tilesets, layers, and
// objects.
type Map struct {
	// Source is the path of to the map definition if loaded from a file.
	Source string
	// Version is the TMX format version.
	Version string
	// TiledVersion is the Tiled version used to save the file. May be a date (for snapshot
	// builds).
	TiledVersion string
	// Class is the class of this map (since 1.9, defaults to “”).
	Class string
	// Orientation describes how tiles are oriented when drawing.
	Orientation Orientation
	// RenderOrder is the order in which tiles on tile layers are rendered. In all cases, the map is drawn row-by-row.
	// Currently is only supported for orthogonal maps.
	RenderOrder RenderOrder
	// Size is the size of the map in tile units.
	Size Size
	// TileSize is the dimensions of tiles on the map in pixel units.
	TileSize Size
	// compressionlevel is the compression level to use for tile layer data (defaults to -1, which means to use the algorithm default).
	compressionlevel int
	// HexSideLength determines the width or height (depending on the staggered axis) of the tile’s edge, in pixels.
	// Only for hexagonal maps.
	HexSideLength int
	// StaggerAxis determines which axis is staggered.
	// Onlye for or staggered and hexagonal maps.
	StaggerAxis StaggerAxis
	// StaggerIndex determines whether the even or odd indexes along the staggered axis are shifted.
	// Only for staggered and hexagonal maps.
	StaggerIndex StaggerIndex
	// ParallaxOrigin is the coordinate of the parallax origin in pixels.
	ParallaxOrigin Vec2
	// BackgroundColor is the background color for the map.
	BackgroundColor Color
	// NextLayerId stores the next available ID for new layers. This number is stored to prevent reuse of the same ID after layers have been removed. (since 1.2) (defaults to the highest layer id in the file + 1)
	NextLayerId int
	// nextobjectid stores the next available ID for new objects. This number is stored to prevent reuse of the same ID after objects have been removed. (since 0.11) (defaults to the highest object id in the file + 1)
	NextObjectId int
	// Infinite indicates whether this map is infinite. An infinite map has no fixed size and can grow in all directions. Its layer data is stored in chunks. (0 for false, 1 for true, defaults to 0)
	Infinite bool
	// Tilesets contains a collection of MapTileset objects used by the Map.
	Tilesets []*MapTileset
	// Properties contain arbitrary key-value pairs of data to associate with the object.
	Properties
	// container is the base container implementation for types that hold a collection of layers.
	container

	cache *Cache
}

func (m *Map) initDefault() {
	if m.cache == nil {
		m.cache = NewCache()
	}
	m.compressionlevel = -1
	m.TileSize = Size{Width: 16, Height: 16}
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (m *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	m.initDefault()
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "version":
			m.Version = attr.Value
		case "tiledversion":
			m.TiledVersion = attr.Value
		case "class":
			m.Class = attr.Value
		case "orientation":
			if value, err := parseOrientation(attr.Value); err != nil {
				return err
			} else {
				m.Orientation = value
			}
		case "renderorder":
			if value, err := parseRenderOrder(attr.Value); err != nil {
				return err
			} else {
				m.RenderOrder = value
			}
		case "compressionlevel":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				m.compressionlevel = value
			} else {
				return err
			}
		case "width":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				m.Size.Width = value
			} else {
				return err
			}
		case "height":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				m.Size.Height = value
			} else {
				return err
			}
		case "tilewidth":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				m.TileSize.Width = value
			} else {
				return err
			}
		case "tileheight":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				m.TileSize.Height = value
			} else {
				return err
			}
		case "hexsidelength":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				m.HexSideLength = value
			} else {
				return err
			}
		case "staggeraxis":
			if attr.Value == "y" {
				m.StaggerAxis = StaggerY
			} else {
				m.StaggerAxis = StaggerX
			}
		case "staggerindex":
			if attr.Value == "odd" {
				m.StaggerIndex = StaggerOdd
			} else {
				m.StaggerIndex = StaggerEven
			}
		case "parallaxoriginx":
			if value, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				m.ParallaxOrigin.X = float32(value)
			} else {
				return err
			}
		case "parallaxoriginy":
			if value, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				m.ParallaxOrigin.Y = float32(value)
			} else {
				return err
			}
		case "backgroundcolor":
			if value, err := ParseColor(attr.Value); err == nil {
				m.BackgroundColor = value
			} else {
				return err
			}
		case "nextlayerid":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				m.NextLayerId = value
			} else {
				return err
			}
		case "nextobjectid":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				m.NextObjectId = value
			} else {
				return err
			}
		case "infinite":
			if value, err := strconv.ParseBool(attr.Value); err == nil {
				m.Infinite = value
			} else {
				return err
			}
		default:
			logAttr(attr.Value, start.Name.Local)
		}
	}

	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		} else if err == io.EOF {
			break
		}

		child, ok := token.(xml.StartElement)
		if ok {
			switch child.Name.Local {
			case "properties":
				m.Properties = make(Properties)
				if err := m.Properties.UnmarshalXML(d, child); err != nil {
					return err
				}
			case "editorsettings":
				// Skip
			case "tileset":
				var tileset MapTileset
				tileset.Map = m
				tileset.cache = m.cache
				if err := tileset.UnmarshalXML(d, child); err != nil {
					return err
				}
				m.Tilesets = append(m.Tilesets, &tileset)
			case "layer":
				var layer TileLayer
				if err := layer.UnmarshalXML(d, child); err != nil {
					return err
				}
				m.AddLayer(&layer)
			case "objectgroup":
				var layer ObjectLayer
				if err := layer.UnmarshalXML(d, child); err != nil {
					return err
				}
				m.AddLayer(&layer)
			case "imagelayer":
				var layer ImageLayer
				if err := layer.UnmarshalXML(d, child); err != nil {
					return err
				}
				m.AddLayer(&layer)
			case "group":
				var layer GroupLayer
				if err := layer.UnmarshalXML(d, child); err != nil {
					return err
				}
				m.AddLayer(&layer)
			default:
				logElem(child.Name.Local, start.Name.Local)
			}
		}

		token, err = d.Token()
	}

	return nil
}


func jsonProp[T any](d *json.Decoder) (value T, err error) {
	var token json.Token
	token, err = d.Token()
	if err != nil {
		return 
	}

	var ok bool
	if value, ok = token.(T); !ok {
		err = errors.New("TODO: Fill this out")
	} 
	return
}

func jsonSkip(d *json.Decoder) error {

	var d1, d2 int

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch token {
		case json.Delim('}'):
			d1--
		case json.Delim(']'):
			d2--
		case json.Delim('{'):
			d1++	
		case json.Delim('['):
			d2++
		}
		if d1 == 0 && d2 == 0 {
			break 
		}
	}

	return nil
}

func (m *Map) UnmarshalJSON(data []byte) error {

	m.initDefault()
	d := json.NewDecoder(bytes.NewReader(data))
	token, err := d.Token()
	if token != json.Delim('{') {
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
		case "backgroundcolor":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if m.BackgroundColor, err = ParseColor(str); err != nil {
				return err
			}
		case "class":
			if m.Class, err = jsonProp[string](d); err != nil {
				return err
			}
		case "compressionlevel":	
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				m.compressionlevel = int(value)
			}
		case "width":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				m.Size.Width = int(value)
			}
		case "height":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				m.Size.Height = int(value)
			}
		case "hexsidelength":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				m.HexSideLength = int(value)
			}
		case "infinite":
			if m.Infinite, err = jsonProp[bool](d); err != nil {
				return err
			}
		case "nextlayerid":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				m.NextLayerId = int(value)
			}
		case "nextobjectid":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				m.NextObjectId = int(value)
			}
		case "tileheight":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				m.TileSize.Height = int(value)
			}
		case "tilewidth":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				m.TileSize.Width = int(value)
			}
		case "version":
			if m.Version, err = jsonProp[string](d); err != nil {
				return err
			}
		case "tiledversion":
			if m.TiledVersion, err = jsonProp[string](d); err != nil {
				return err
			}		
		case "parallaxoriginx":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				m.ParallaxOrigin.X = float32(value)
			}
		case "parallaxoriginy":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				m.ParallaxOrigin.Y = float32(value)
			}
		case "type":
			if _, err = d.Token(); err != nil {
				return err
			}
		case "orientation":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if m.Orientation, err = parseOrientation(str); err != nil {
				return err
			}
		case "renderorder":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if m.RenderOrder, err = parseRenderOrder(str); err != nil {
				return err
			}
		case "staggeraxis":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if m.StaggerAxis, err = parseStaggerAxis(str); err != nil {
				return err
			}
		case "staggerindex":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if m.StaggerIndex, err = parseStaggerIndex(str); err != nil {
				return err
			}
		case "properties":
			m.Properties = make(Properties)
			if err = d.Decode(&m.Properties); err != nil {
				return err
			}
		case "layers":
			if token, err = d.Token(); err != nil {
				return err
			} else if token != json.Delim('[') {
				return errors.New("expected JSON array")
			}
			for d.More() {
				var layer jsonLayer
				if err = d.Decode(&layer); err != nil {
					return err
				}
				m.AddLayer(layer.toLayer())
			}
			// Position to next token ']'
			if token, err = d.Token(); err != nil {
				return err
			}
		case "tilesets":
			if token, err = d.Token(); err != nil {
				return err
			} else if token != json.Delim('[') {
				return errors.New("expected JSON array")
			}
			for d.More() {
				var tileset MapTileset
				tileset.Map = m
				tileset.cache = m.cache

				if err = d.Decode(&tileset); err != nil {
					return err
				}
				m.Tilesets = append(m.Tilesets, &tileset)
			}
			// Position to next token ']'
			if token, err = d.Token(); err != nil {
				return err
			}
		default: 
			jsonSkip(d)
		}
	}


	//
	// for _, tileset := range temp.Tilesets {
	// 	// Skip embedded tilsets
	// 	if tileset.Source == "" {
	// 		continue
	// 	}
	// 	if ts, err := OpenTileset(tileset.Source, m.cache); err != nil {
	// 		return err
	// 	} else {
	// 			}
	// 	// TODO
	// }
	
	// for _, layer := range temp.Layers {
	// 	m.AddLayer(layer.toLayer())
	// }

	return nil
}

func (m *Map) AddTileset(ts *Tileset, first TileID) {
	value := MapTileset{Tileset: ts, FirstGID: first}
	m.Tilesets = append(m.Tilesets, &value)
}

// linkLayer configures the Prev/Next values of new layer, as well as the Head/Tail of the map.
func (m *Map) AddLayer(layer Layer) {
	switch v := layer.(type) {
	case *TileLayer:
		m.TileLayers = append(m.TileLayers, v)
	case *ImageLayer:
		m.ImageLayers = append(m.ImageLayers, v)
	case *ObjectLayer:
		m.ObjectLayers = append(m.ObjectLayers, v)
	case *GroupLayer:
		m.GroupLayers = append(m.GroupLayers, v)
	}

	if m.head == nil {
		m.head = layer
	}

	if m.tail != nil {
		m.tail.setNext(layer)
		layer.setPrev(m.tail)
	}
	m.tail = layer
	m.head.setParent(m)
	m.head.setContainer(m)
}

// OpenMap reads a tilemap from a file, automatically detecting it format.
//
// An optional cache can be supplied that maintains references to tilesets and
// templates to prevent frequent re-processing of them. When nil, an internal
// cache will be used that only exists for the lifetime of the map.
func OpenMap(path string, cache *Cache) (*Map, error) {
	return OpenMapFormat(path, detectFileExt(path), cache)
}

// OpenMapFormat reads a tilemap from a file, using the specified format.
//
// An optional cache can be supplied that maintains references to tilesets and
// templates to prevent frequent re-processing of them. When nil, an internal
// cache will be used that only exists for the lifetime of the map.
func OpenMapFormat(path string, format Format, cache *Cache) (*Map, error) {
	var abs string
	var err error
	if abs, err = FindPath(path); err != nil {
		return nil, err
	} 

	reader, _, err := getStream(abs)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	IncludePaths = append(IncludePaths, filepath.Dir(abs))
	defer func() { IncludePaths = IncludePaths[:len(IncludePaths)-1] }()
	
	var tilemap Map
	tilemap.Source = abs
	tilemap.cache = cache

	if err = ReadMapFormat(reader, format, &tilemap); err != nil {
		return nil, err
	}
	return &tilemap, nil
}

// ReadMap reads a tilemap from the current position in the reader.
func ReadMap(r io.ReadSeeker, tilemap *Map) (error) {
	return ReadMapFormat(r, detectReader(r), tilemap)
}

// ReadMapFormat reads a tilemap from the current position in the reader using
// the specified format.
func ReadMapFormat(r io.Reader, format Format, tilemap *Map) error {
	switch format {
	case FormatXML:
		d := xml.NewDecoder(r)
		if err := d.Decode(tilemap); err != nil {
			return err
		}
	case FormatJSON:
		d := json.NewDecoder(r)
		if err := d.Decode(tilemap); err != nil {
			return err
		}
	default:
		return errInvalidEnum("Format", fmt.Sprintf("Format(%d)", format))
	}

	return nil
}

// vim: ts=4
