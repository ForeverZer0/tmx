package tmx

import "encoding/xml"

// GroupLayer is a map layer that acts as a container for other map layers. Its offset,
// visibility, opacity, and tint recursively affect child layers.
type GroupLayer struct {
	baseLayer
	container
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (layer *GroupLayer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	layer.initDefaults(LayerGroup)

	for _, attr := range start.Attr {
		if handled, err := layer.baseLayer.xmlAttr(attr); err != nil {
			return err
		} else if handled {
			continue
		}
		logAttr(attr.Name.Local, start.Name.Local)
	}

	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		}

		if child, ok := token.(xml.StartElement); ok {
			if handled, err := layer.baseLayer.xmlProp(d, child); err != nil {
				return err
			} else if !handled {
				switch child.Name.Local {
				case "layer":
					var value TileLayer
					if err := d.DecodeElement(&value, &child); err != nil {
						return err
					}
					layer.AddLayer(&value)
				case "objectgroup":
					var value ObjectLayer
					if err := d.DecodeElement(&value, &child); err != nil {
						return err
					}
					layer.AddLayer(&value)
				case "imagelayer":
					var value ImageLayer
					if err := d.DecodeElement(&value, &child); err != nil {
						return err
					}
					layer.AddLayer(&value)
				case "group":
					var value GroupLayer
					if err := d.DecodeElement(&value, &child); err != nil {
						return err
					}
					layer.AddLayer(&value)
				default:
					logElem(child.Name.Local, start.Name.Local)
				}
			}
		}
		token, err = d.Token()
	}

	return nil
}

// AddLayer appends a new layer to the group.
func (g *GroupLayer) AddLayer(layer Layer) {
	g.container.AddLayer(layer)
	g.head.setParent(g.parent)
}

// vim: ts=4
