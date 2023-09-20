package tmx

import "encoding/xml"

// TileLayer describes a map layer that is composed of tile data from a Tileset.
type TileLayer struct {
	baseLayer
	TileData
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (layer *TileLayer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	layer.initDefaults(LayerTile)

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
				case "data":
					if err := d.DecodeElement(&layer.TileData, &child); err != nil {
						return err
					}
				default:
					logElem(child.Name.Local, start.Name.Local)
				}
			}
		}
		token, err = d.Token()
	}
	return layer.postProcess(layer.Width * layer.Height)
}

// vim: ts=4
