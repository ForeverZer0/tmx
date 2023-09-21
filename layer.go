package tmx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
)

// Layer describes a map layer of any type, providing doubly linked-list functionality for
// iteration. The concrete type of the interface will be one of the following:
//
//   - *TileLayer
//   - *ObjectLayer
//   - *ImageLayer
//   - *GroupLayer
type Layer interface {
	// Map returns the top-level map that contains the layer.
	Map() *Map
	// Type returns a constant describing the type of layer this interface represents.
	Type() LayerType
	// Container returns the direct parent container of the layer. For layers within a group,
	// this will be the GroupLayer, otherwise it will be the parent Map.
	Container() Container
	// Next returns the next map layer, or nil when called by the tail layer.
	Next() Layer
	// Prev returns the previous map layer, or nil when called by the head layer.
	Prev() Layer

	setPrev(layer Layer)
	setNext(layer Layer)
	setParent(parent *Map)
	setContainer(container Container)
}

// jsonLayer is used internally to marshal JSON-formatted layers. The differences between the
// structure of the JSON/XML are significantly different, hence it is much easier to just
// use a dedicated struct for it and let Go handle the marshal logic.
type jsonLayer struct {
	Rect
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Class      string     `json:"class"`
	Type       LayerType  `json:"type"`
	Offset     Vec2       `json:"offsetx"`
	Parallax   Vec2       `json:"parallaxy"`
	Opacity    float32    `json:"opacity"`
	Start      Point      `json:"startx"`
	TintColor  Color      `json:"tintcolor"`
	Visible    bool       `json:"visible"`
	Properties Properties `json:"properties"`

	// Tile layer only
	Data TileData

	// Image layer only
	TransparentColor Color  `json:"transparentcolor"`
	Image            string `json:"image"`
	RepeatX          bool   `json:"repeatx"`
	RepeatY          bool   `json:"repeaty"`

	// Group only
	Layers []jsonLayer `json:"layers"`

	// ObjectGroup
	Objects   []Object  `json:"objects"`
	DrawOrder DrawOrder `json:"draworder"`

	cache *Cache
}

func (l *jsonLayer) UnmarshalJSON(data []byte) error {
	l.Opacity = 1.0
	l.Visible = true

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
		case "id":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.ID = int(value)
			}
		case "name":
			if l.Name, err = jsonProp[string](d); err != nil {
				return err
			}
		case "class":
			if l.Class, err = jsonProp[string](d); err != nil {
				return err
			}
		case "type":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if value, err := parseLayerType(str); err != nil {
				return err
			} else {
				l.Type = value
			}
		case "x":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.X = int(value)
			}
		case "y":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.Y = int(value)
			}
		case "width":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.Width = int(value)
			}
		case "height":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.Height = int(value)
			}
		case "offsetx":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.Offset.X = float32(value)
			}
		case "offsety":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.Offset.Y = float32(value)
			}
		case "parallaxx":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.Parallax.X = float32(value)
			}
		case "parallaxy":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.Parallax.Y = float32(value)
			}
		case "opacity":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.Opacity = float32(value)
			}
		case "startx":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.Start.X = int(value)
			}
		case "starty":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				l.Start.Y = int(value)
			}
		case "tintcolor":
			if value, err := jsonProp[string](d); err != nil {
				return err
			} else if color, err := ParseColor(value); err != nil {
				return err
			} else {
				l.TintColor = color
			}
		case "visible":
			if value, err := jsonProp[bool](d); err != nil {
				return err
			} else {
				l.Visible = value
			}
		case "properties":
			props := make(Properties)
			if err := d.Decode(&props); err != nil {
				return err
			}
			l.Properties = props
		case "compression":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if value, err := parseCompression(str); err != nil {
				return err
			} else {
				l.Data.Compression = value
			}
		case "encoding":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if value, err := parseEncoding(str); err != nil {
				return err
			} else {
				l.Data.Encoding = value
			}
		case "draworder":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if value, err := parseDrawOrder(str); err != nil {
				return err
			} else {
				l.DrawOrder = value
			}
		case "transparentcolor":
			if value, err := jsonProp[string](d); err != nil {
				return err
			} else if color, err := ParseColor(value); err != nil {
				return err
			} else {
				l.TransparentColor = color
			}
		case "image":
			if l.Image, err = jsonProp[string](d); err != nil {
				return err
			}
		case "repeatx":
			if value, err := jsonProp[bool](d); err != nil {
				return err
			} else {
				l.RepeatX = value
			}
		case "repeaty":
			if value, err := jsonProp[bool](d); err != nil {
				return err
			} else {
				l.RepeatY = value
			}
		case "objects":
			if token, err = d.Token(); err != nil {
				return err
			} else if token != json.Delim('[') {
				return ErrExpectedArray
			}
			for d.More() {
				var obj Object
				obj.cache = l.cache
				if err = d.Decode(&obj); err != nil {
					return err
				}
				l.Objects = append(l.Objects, obj)
			}
			// Position to next token ']'
			if token, err = d.Token(); err != nil {
				return err
			}
		case "layers":
			if token, err = d.Token(); err != nil {
				return err
			} else if token != json.Delim('[') {
				return ErrExpectedArray
			}
			for d.More() {
				var child jsonLayer
				child.cache = l.cache
				if err = d.Decode(&child); err != nil {
					return err
				}
				l.Layers = append(l.Layers, child)
			}
			// Position to next token ']'
			if token, err = d.Token(); err != nil {
				return err
			}
		case "chunks":
			if token, err = d.Token(); err != nil {
				return err
			} else if token != json.Delim('[') {
				return ErrExpectedArray
			}
			for d.More() {
				var chunk Chunk
				if err = d.Decode(&chunk); err != nil {
					return err
				}
				l.Data.Chunks = append(l.Data.Chunks, chunk)
			}
			// Position to next token ']'
			if token, err = d.Token(); err != nil {
				return err
			}
		case "data":
			if token, err = d.Token(); err != nil {
				return err
			}

			if token == json.Delim('[') {
				for d.More() {
					if value, err := jsonProp[float64](d); err != nil {
						return err
					} else {
						l.Data.Tiles = append(l.Data.Tiles, TileID(value))
					}
				}
				// Position to next token ']'
				if token, err = d.Token(); err != nil {
					return err
				}

			} else if str, ok := token.(string); ok {
				l.Data.tileData = []byte(str)
			}
		default:
			jsonSkip(d)
		}
	}

	if l.Type == LayerTile {
		return l.Data.postProcess(l.Width * l.Height)
	}

	return nil
}

func (j *jsonLayer) toLayer() Layer {
	// TODO: StartX, StartY? The are documented, but no setting in Tiled uses them, nor are they
	// ever actually present(?)

	var base *baseLayer
	var layer Layer

	switch j.Type {
	case LayerTile:
		impl := TileLayer{TileData: j.Data}
		base = &impl.baseLayer
		layer = &impl
	case LayerImage:
		var impl ImageLayer
		layer = &impl
		base = &impl.baseLayer
		impl.Image = &Image{
			Source:       j.Image,
			Transparency: j.TransparentColor,
		}
		impl.RepeatX = j.RepeatX
		impl.RepeatY = j.RepeatY
	case LayerObject:
		var impl ObjectLayer
		layer = &impl
		base = &impl.baseLayer
		impl.Objects = j.Objects
		impl.DrawOrder = j.DrawOrder
	case LayerGroup:
		var impl GroupLayer
		layer = &impl
		base = &impl.baseLayer
		for i := range j.Layers {
			child := j.Layers[i].toLayer()
			impl.AddLayer(child)
		}
	}

	base.ID = j.ID
	base.Name = j.Name
	base.Class = j.Class
	base.layerType = j.Type
	base.Offset = j.Offset
	base.Parallax = j.Parallax
	base.Opacity = j.Opacity
	base.Visible = j.Visible
	base.TintColor = j.TintColor
	base.Properties = j.Properties
	base.Rect = j.Rect

	return layer
}

type baseLayer struct {
	// ID is a unique identifier for the layer. A value of 0 is invalid, and no two layers
	// will share the same ID.
	ID int
	// Name is the user-defined name of the layer.
	Name string
	// Class is the user-defined class of the layer.
	Class string
	// Offset is the rendering offset of the layer on each axis in pixel units.
	Offset Vec2
	// Parallax defines the parallax factor for the layer on each axis.
	Parallax Vec2
	// Opacity defines the transparency factor of the layer, where 0.0 is fully opaque and
	// and invisible, and 1.0 indicates no transparency.
	Opacity float32
	// Visible indicates whether the layer should be rendered or not.
	Visible bool
	// layerType describes the type of the outer layer.
	layerType LayerType
	// TintColor is a color that is multiplied by with any graphics used by this layer.
	TintColor Color
	// Rect is the offset and size of the layer in tile units.
	//
	// The location of the rectangle is always <0,0>, and cannot be edited in Tiled.
	Rect
	// Properties contain arbitrary key-value pairs of data to associate with the object.
	Properties
	// parent maintains a reference to the parent map.
	parent *Map
	// container maintains a reference to the container this layer is within.
	container Container
	// next maintains a reference to the next layer in the linked-list.
	next Layer
	// next maintains a reference to the previous layer in the linked-list.
	prev Layer
	// cache is a reference to the parent map's Cache.
	cache *Cache
}

// String implements the Stringer interface.
func (l *baseLayer) String() string {
	return fmt.Sprintf(`Tileset("%s")`, l.Name)
}

// xmlAttr attempts to process the given attribute into the base layer type, returning whether
// it was handled or not and if an error occurred.
func (layer *baseLayer) xmlAttr(attr xml.Attr) (bool, error) {
	switch attr.Name.Local {
	case "id":
		if value, err := strconv.Atoi(attr.Value); err != nil {
			return false, err
		} else {
			layer.ID = value
		}
	case "name":
		layer.Class = attr.Value
	case "class":
		layer.Class = attr.Value
	case "tintcolor":
		if color, err := ParseColor(attr.Value); err != nil {
			return false, err
		} else {
			layer.TintColor = color
		}
	case "offsetx":
		if value, err := strconv.ParseFloat(attr.Value, 32); err != nil {
			return false, err
		} else {
			layer.Offset.X = float32(value)
		}
	case "offsety":
		if value, err := strconv.ParseFloat(attr.Value, 32); err != nil {
			return false, err
		} else {
			layer.Offset.Y = float32(value)
		}
	case "parallaxx":
		if value, err := strconv.ParseFloat(attr.Value, 32); err != nil {
			return false, err
		} else {
			layer.Parallax.X = float32(value)
		}
	case "parallaxy":
		if value, err := strconv.ParseFloat(attr.Value, 32); err != nil {
			return false, err
		} else {
			layer.Parallax.Y = float32(value)
		}
	case "opacity":
		if value, err := strconv.ParseFloat(attr.Value, 32); err != nil {
			return false, err
		} else {
			layer.Opacity = float32(value)
		}
	case "visible":
		if value, err := strconv.ParseBool(attr.Value); err != nil {
			return false, err
		} else {
			layer.Visible = value
		}
	case "x":
		if value, err := strconv.Atoi(attr.Value); err != nil {
			return false, err
		} else {
			layer.X = value
		}
	case "y":
		if value, err := strconv.Atoi(attr.Value); err != nil {
			return false, err
		} else {
			layer.Y = value
		}
	case "width":
		if value, err := strconv.Atoi(attr.Value); err != nil {
			return false, err
		} else {
			layer.Width = value
		}
	case "height":
		if value, err := strconv.Atoi(attr.Value); err != nil {
			return false, err
		} else {
			layer.Height = value
		}
	default:
		return false, nil
	}
	return true, nil
}

// xmlProp attempts to process the given element into the base layer type, returning whether
// it was handled or not and if an error occurred.
func (layer *baseLayer) xmlProp(d *xml.Decoder, start xml.StartElement) (bool, error) {
	switch start.Name.Local {
	case "properties":
		layer.Properties = make(Properties)
		if err := d.DecodeElement(&layer.Properties, &start); err != nil {
			return false, err
		}
	default:
		return false, nil
	}
	return true, nil
}

// Type returns a constant describing the type of layer.
func (layer *baseLayer) Type() LayerType {
	return layer.layerType
}

// Next returns the next map layer, or nil when called by the tail layer.
func (layer *baseLayer) Next() Layer {
	return layer.next
}

// Prev returns the previous map layer, or nil when called by the head layer.
func (layer *baseLayer) Prev() Layer {
	return layer.prev
}

// Map returns the top-level map that contains the layer.
func (layer *baseLayer) Map() *Map {
	return layer.parent
}

// Container returns the direct parent container of the layer. For layers within a group,
// this will be the GroupLayer, otherwise it will be the parent Map.
func (layer *baseLayer) Container() Container {
	if layer.container != nil {
		return layer.container
	}
	return layer.parent
}

// setPrev implements the Layer interface.
func (layer *baseLayer) setPrev(prev Layer) {
	layer.prev = prev
}

// setNext implements the Layer interface.
func (layer *baseLayer) setNext(next Layer) {
	layer.next = next
}

// setParent implements the Layer interface.
func (layer *baseLayer) setParent(parent *Map) {
	layer.parent = parent
}

// setContainer implements the Layer interface.
func (layer *baseLayer) setContainer(container Container) {
	layer.container = container
}

// initDefaults initializes default values of a layer.
func (layer *baseLayer) initDefaults(lt LayerType) {
	layer.layerType = lt
	layer.Opacity = 1.0
	layer.Visible = true
}

// vim: ts=4
