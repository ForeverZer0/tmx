package tmx

import (
	"encoding/xml"
	"strconv"
)

// ImageLayer describes a map layer that contains a single image.
type ImageLayer struct {
	baseLayer
	// RepeatX determines whether the image drawn by this layer is repeated along the x-axis.
	RepeatX bool
	// RepeatY determines whether the image drawn by this layer is repeated along the y-axis.
	RepeatY bool
	// Image is the image used by this layer.
	Image *Image
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (layer *ImageLayer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if handled, err := layer.baseLayer.xmlAttr(attr); err != nil {
			return err
		} else if handled {
			continue
		}

		switch attr.Name.Local {
		case "repeatx":
			if value, err := strconv.ParseBool(attr.Value); err != nil {
				return err
			} else {
				layer.RepeatX = value
			}
		case "repeaty":
			if value, err := strconv.ParseBool(attr.Value); err != nil {
				return err
			} else {
				layer.RepeatY = value
			}
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
			if handled, err := layer.baseLayer.xmlProp(d, child); err != nil {
				return err
			} else if !handled {
				switch child.Name.Local {
				case "image":
					var img Image
					if err := img.UnmarshalXML(d, child); err != nil {
						return err
					}
					layer.Image = &img
				default:
					logElem(child.Name.Local, start.Name.Local)
				}
			}
		}
		token, err = d.Token()
	}
	return nil
}

// vim: ts=4
