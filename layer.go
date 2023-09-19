package tmx

import (
	"encoding/xml"
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
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Class      string     `json:"class"`
	Type       LayerType  `json:"type"`
	X          int        `json:"x"`
	Y          int        `json:"y"`
	Width      int        `json:"width"`
	Height     int        `json:"height"`
	OffsetX    float32    `json:"offsetx"`
	OffsetY    float32    `json:"offsety"`
	ParallaxX  float32    `json:"parallaxx"`
	ParallaxY  float32    `json:"parallaxy"`
	Opacity    float32    `json:"opacity"`
	StartX     int        `json:"startx"`
	StartY     int        `json:"starty"`
	TintColor  Color      `json:"tintcolor"`
	Visible    bool       `json:"visible"`
	Properties Properties `json:"properties"`

	// Tile layer only
	Chunks      []Chunk     `json:"chunks"`
	Compression Compression `json:"compression"`
	Encoding    Encoding    `json:"encoding"`
	Data        []byte      `json:"data"`

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

// vim: ts=4
