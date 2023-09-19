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
					layer.linkLayer(&value)
					layer.TileLayers = append(layer.TileLayers, value)
				case "objectgroup":
					var value ObjectLayer
					if err := d.DecodeElement(&value, &child); err != nil {
						return err
					}
					layer.linkLayer(&value)
					layer.ObjectLayers = append(layer.ObjectLayers, value)
				case "imagelayer":
					var value ImageLayer
					if err := d.DecodeElement(&value, &child); err != nil {
						return err
					}
					layer.linkLayer(&value)
					layer.ImageLayers = append(layer.ImageLayers, value)
				case "group":
					var value GroupLayer
					if err := d.DecodeElement(&value, &child); err != nil {
						return err
					}
					layer.linkLayer(&value)
					layer.GroupLayers = append(layer.GroupLayers, value)
				default:
					logElem(child.Name.Local, start.Name.Local)
				}
			}
		}
		token, err = d.Token()
	}

	return nil
}

// linkLayer configures the Prev/Next values of new layer, as well as the Head/Tail of the map.
func (g *GroupLayer) linkLayer(layer Layer) {
	if g.head == nil {
		g.head = layer
	}

	if g.tail != nil {
		g.tail.setNext(layer)
		layer.setPrev(g.tail)
	}
	g.tail = layer
	g.head.setParent(g.parent)
	g.head.setContainer(g)
}

// vim: ts=4
