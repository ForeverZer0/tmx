package tmx

import (
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

func jsonLayer(d *json.Decoder, cache *Cache) (Layer, error) {
	var base baseLayer
	base.cache = cache
	base.Opacity = 1.0
	base.Visible = true

	var objects []Object
	var tileData TileData
	var layers []Layer
	var image Image
	var repeatX, repeatY bool
	var order DrawOrder
	var start Point // TODO: Start is not actually used, but is documented in spec

	token, err := d.Token()
	if err != nil {
		return nil, err
	} else if token != json.Delim('{') {
		return nil, ErrExpectedObject
	}

	for {
		if token, err = d.Token(); err != nil {
			return nil, err
		} else if token == json.Delim('}') {
			break
		}

		name := token.(string)
		switch name {
		case "id":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				base.ID = int(value)
			}
		case "name":
			if base.Name, err = jsonProp[string](d); err != nil {
				return nil, err
			}
		case "class":
			if base.Class, err = jsonProp[string](d); err != nil {
				return nil, err
			}
		case "type":
			if str, err := jsonProp[string](d); err != nil {
				return nil, err
			} else if value, err := parseLayerType(str); err != nil {
				return nil, err
			} else {
				base.layerType = value
			}
		case "x":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				base.X = int(value)
			}
		case "y":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				base.Y = int(value)
			}
		case "width":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				base.Width = int(value)
			}
		case "height":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				base.Height = int(value)
			}
		case "offsetx":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				base.Offset.X = float32(value)
			}
		case "offsety":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				base.Offset.Y = float32(value)
			}
		case "parallaxx":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				base.Parallax.X = float32(value)
			}
		case "parallaxy":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				base.Parallax.Y = float32(value)
			}
		case "opacity":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				base.Opacity = float32(value)
			}
		case "startx":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				start.X = int(value)
			}
		case "starty":
			if value, err := jsonProp[float64](d); err != nil {
				return nil, err
			} else {
				start.Y = int(value)
			}
		case "tintcolor":
			if value, err := jsonProp[string](d); err != nil {
				return nil, err
			} else if color, err := ParseColor(value); err != nil {
				return nil, err
			} else {
				base.TintColor = color
			}
		case "visible":
			if value, err := jsonProp[bool](d); err != nil {
				return nil, err
			} else {
				base.Visible = value
			}
		case "properties":
			props := make(Properties)
			if err := d.Decode(&props); err != nil {
				return nil, err
			}
			base.Properties = props
		case "compression":
			if str, err := jsonProp[string](d); err != nil {
				return nil, err
			} else if value, err := parseCompression(str); err != nil {
				return nil, err
			} else {
				tileData.Compression = value
			}
		case "encoding":
			if str, err := jsonProp[string](d); err != nil {
				return nil, err
			} else if value, err := parseEncoding(str); err != nil {
				return nil, err
			} else {
				tileData.Encoding = value
			}
		case "draworder":
			if str, err := jsonProp[string](d); err != nil {
				return nil, err
			} else if value, err := parseDrawOrder(str); err != nil {
				return nil, err
			} else {
				order = value
			}
		case "transparentcolor":
			if value, err := jsonProp[string](d); err != nil {
				return nil, err
			} else if color, err := ParseColor(value); err != nil {
				return nil, err
			} else {
				image.Transparency = color
			}
		case "image":
			if image.Source, err = jsonProp[string](d); err != nil {
				return nil, err
			}
		case "repeatx":
			if value, err := jsonProp[bool](d); err != nil {
				return nil, err
			} else {
				repeatX = value
			}
		case "repeaty":
			if value, err := jsonProp[bool](d); err != nil {
				return nil, err
			} else {
				repeatY = value
			}
		case "objects":
			if token, err = d.Token(); err != nil {
				return nil, err
			} else if token != json.Delim('[') {
				return nil, ErrExpectedArray
			}
			for d.More() {
				var obj Object
				obj.cache = cache
				if err = d.Decode(&obj); err != nil {
					return nil, err
				}
				objects = append(objects, obj)
			}
			// Position to next token ']'
			if token, err = d.Token(); err != nil {
				return nil, err
			}
		case "layers":
			if token, err = d.Token(); err != nil {
				return nil, err
			} else if token != json.Delim('[') {
				return nil, ErrExpectedArray
			}
			for d.More() {
				if child, err := jsonLayer(d, cache); err != nil {
					return nil, err
				} else {
					layers = append(layers, child)
				}
			}
			// Position to next token ']'
			if token, err = d.Token(); err != nil {
				return nil, err
			}
		case "chunks":
			if token, err = d.Token(); err != nil {
				return nil, err
			} else if token != json.Delim('[') {
				return nil, ErrExpectedArray
			}
			for d.More() {
				var chunk Chunk
				if err = d.Decode(&chunk); err != nil {
					return nil, err
				}
				tileData.Chunks = append(tileData.Chunks, chunk)
			}
			// Position to next token ']'
			if token, err = d.Token(); err != nil {
				return nil, err
			}
		case "data":
			if token, err = d.Token(); err != nil {
				return nil, err
			}

			if token == json.Delim('[') {
				for d.More() {
					if value, err := jsonProp[float64](d); err != nil {
						return nil, err
					} else {
						tileData.Tiles = append(tileData.Tiles, TileID(value))
					}
				}
				// Position to next token ']'
				if token, err = d.Token(); err != nil {
					return nil, err
				}

			} else if str, ok := token.(string); ok {
				tileData.tileData = []byte(str)
			}
		default:
			jsonSkip(d)
		}
	}
	
	switch base.layerType {
	case LayerTile:
		impl := TileLayer{baseLayer: base, TileData: tileData}
		if err = impl.TileData.postProcess(base.Area()); err != nil {
			return nil, err
		}
		impl.calcChunks()
		return &impl, nil
	case LayerImage:
		impl := ImageLayer{baseLayer: base, Image: &image, RepeatX: repeatX, RepeatY: repeatY}
		return &impl, nil
	case LayerObject:
		impl := ObjectLayer{baseLayer: base, Objects: objects, DrawOrder: order}
		return &impl, nil
	case LayerGroup:
		impl := GroupLayer{baseLayer: base}
		for _, child := range layers {
			impl.AddLayer(child)
		}
		return &impl, nil
	default:
		return nil, errInvalidEnum("LayerType", base.layerType.String())
	}
}

// vim: ts=4
