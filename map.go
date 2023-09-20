package tmx

import (
	"encoding/json"
	"encoding/xml"
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
	// nextLayerId stores the next available ID for new layers. This number is stored to prevent reuse of the same ID after layers have been removed. (since 1.2) (defaults to the highest layer id in the file + 1)
	nextLayerId int
	// nextobjectid stores the next available ID for new objects. This number is stored to prevent reuse of the same ID after objects have been removed. (since 0.11) (defaults to the highest object id in the file + 1)
	nextObjectId int
	// Infinite indicates whether this map is infinite. An infinite map has no fixed size and can grow in all directions. Its layer data is stored in chunks. (0 for false, 1 for true, defaults to 0)
	Infinite bool
	// Tilesets contains a collection of MapTileset objects used by the Map.
	Tilesets []*MapTileset
	// Properties contain arbitrary key-value pairs of data to associate with the object.
	Properties
	// container is the base container implementation for types that hold a collection of layers.
	container
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (m *Map) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	m.compressionlevel = -1

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
				m.nextLayerId = value
			} else {
				return err
			}
		case "nextobjectid":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				m.nextObjectId = value
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

func (m *Map) UnmarshalJSON(data []byte) error {
	type jsonMap struct {
		BackgroundColor  Color        `json:"backgroundcolor"`
		Class            string       `json:"class"`
		Compressionlevel int          `json:"compressionlevel"`
		Width            int          `json:"width"`
		Height           int          `json:"height"`
		HexSideLength    int          `json:"hexsidelength"`
		Infinite         bool         `json:"infinite"`
		Layers           []jsonLayer  `json:"layers"`
		NextLayerID      int          `json:"nextlayerid"`
		NextObjectID     int          `json:"nextobjectid"`
		Orientation      Orientation  `json:"orientation"`
		ParallaxOriginX  float64      `json:"parallaxoriginx"`
		ParallaxOriginY  float64      `json:"parallaxoriginy"`
		Properties       Properties   `json:"properties"`
		RenderOrder      RenderOrder  `json:"renderorder"`
		StaggerAxis      StaggerAxis  `json:"staggeraxis"`
		StaggerIndex     StaggerIndex `json:"staggerindex"`
		TiledVersion     string       `json:"tiledversion"`
		TileHeight       int          `json:"tileheight"`
		TileWidth        int          `json:"tilewidth"`
		Tilesets         []Tileset    `json:"tilesets"`
		Type             string       `json:"type"`
		Version          string       `json:"version"`
	}
	
	var temp jsonMap
	temp.Properties = make(Properties)

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	
	m.Version = temp.Version
	m.TiledVersion = temp.TiledVersion
	m.Class = temp.Class
	m.Orientation = temp.Orientation
	m.RenderOrder = temp.RenderOrder
	m.BackgroundColor = temp.BackgroundColor
	m.Size = Size{Width: temp.Width, Height: temp.Height}
	m.TileSize = Size{Width: temp.TileWidth, Height: temp.TileHeight}
	m.compressionlevel = temp.Compressionlevel
	m.HexSideLength = temp.HexSideLength
	m.StaggerAxis = temp.StaggerAxis
	m.StaggerIndex = temp.StaggerIndex
	m.ParallaxOrigin = Vec2{X: float32(temp.ParallaxOriginX), Y: float32(temp.ParallaxOriginY)}
	m.nextLayerId = temp.NextLayerID
	m.nextObjectId = temp.NextObjectID
	m.Infinite = temp.Infinite
	m.Properties = temp.Properties

	for _, layer := range temp.Layers {
		m.AddLayer(layer.toLayer())
	}

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
	m.head.setParent(nil)
	m.head.setContainer(m)
}

// OpenMap reads a tilemap from a file, automatically detecting it format.
func OpenMap(path string) (*Map, error) {
	reader, abs, ft, err := getStream(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	IncludePaths = append(IncludePaths, filepath.Dir(abs))
	defer func() { IncludePaths = IncludePaths[:len(IncludePaths)-1] }()

	return ReadMapFormat(reader, ft)
}

// OpenMapFormat reads a tilemap from a file, using the specified format.
func OpenMapFormat(path string, format Format) (tilemap *Map, err error) {
	reader, abs, _, err := getStream(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	IncludePaths = append(IncludePaths, filepath.Dir(abs))
	defer func() { IncludePaths = IncludePaths[:len(IncludePaths)-1] }()

	return ReadMapFormat(reader, format)
}

// ReadMap reads a tilemap from the current position in the reader.
func ReadMap(r io.ReadSeeker) (*Map, error) {
	return ReadMapFormat(r, detectReader(r))
}

// ReadMapFormat reads a tilemap from the current position in the reader using
// the specified format.
func ReadMapFormat(r io.Reader, format Format) (*Map, error) {
	var tilemap Map
	switch format {
	case FormatXML:
		d := xml.NewDecoder(r)
		if err := d.Decode(&tilemap); err != nil {
			return nil, err
		}
	case FormatJSON:
		d := json.NewDecoder(r)
		if err := d.Decode(&tilemap); err != nil {
			return nil, err
		}
	default:
		return nil, errInvalidEnum("Format", fmt.Sprintf("Format(%d)", format))
	}

	return &tilemap, nil
}

// vim: ts=4
