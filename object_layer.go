package tmx

import (
	"encoding/xml"
	"strconv"
)

// ObjectLayer is a map layer that contains map objects.
type ObjectLayer struct {
	baseLayer
	// Color is the color used to display the objects in the group (optional).
	Color Color
	// DrawOrder determines the order in which objects in the layer are rendered.
	DrawOrder DrawOrder
	// Objects is the collection of objects to be rendered in this layer.
	Objects []Object
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (layer *ObjectLayer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	layer.initDefaults(LayerImage)

	for _, attr := range start.Attr {
		if handled, err := layer.xmlAttr(attr); err != nil {
			return err
		} else if handled {
			continue
		}

		switch attr.Name.Local {
		case "color":
			if value, err := ParseColor(attr.Value); err == nil {
				layer.Color = value
			} else {
				return err
			}
		case "draworder":
			if value, err := parseDrawOrder(attr.Value); err != nil {
				return err
			} else {
				layer.DrawOrder = value
			}
		case "x":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				layer.X = value
			} else {
				return err
			}
		case "y":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				layer.Y = value
			} else {
				return err
			}
		case "width":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				layer.Width = value
			} else {
				return err
			}
		case "height":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				layer.Height = value
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
		}

		next, ok := token.(xml.StartElement)
		if ok {
			switch next.Name.Local {
			case "properties":
				layer.Properties = make(Properties)
				if err := layer.Properties.UnmarshalXML(d, next); err != nil {
					return err
				}
			case "object":
				var obj Object
				obj.cache = layer.cache
				if err := obj.UnmarshalXML(d, next); err != nil {
					return err
				}
				layer.Objects = append(layer.Objects, obj)
			default:
				logElem(next.Name.Local, start.Name.Local)
			}
		}

		token, err = d.Token()
	}

	return nil
}

// vim: ts=4
