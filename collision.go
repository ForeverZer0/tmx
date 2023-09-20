package tmx

import (
	"encoding/xml"
	"strconv"
)

// Collision describes a collection of basic Objects that are used to defined shapes for
// tile collision.
type Collision struct {
	// ID is a unique identifier for the Collision instance. No two collisions within the same
	// tileset will have the same ID.
	ID int
	// DrawOrder specifies the order in which the shapes whould be "drawn".
	DrawOrder DrawOrder
	// Objects is a collection of shaped objects defining the collision.
	Objects []Object
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (c *Collision) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				c.ID = value
			} else {
				return err
			}
		case "draworder":
			switch attr.Value {
			case "index":
				c.DrawOrder = DrawIndex
			default:
				c.DrawOrder = DrawTopDown
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
			if child.Name.Local != "object" {
				logElem(child.Name.Local, start.Name.Local)
			} else {
				var object Object
				if err := object.UnmarshalXML(d, child); err != nil {
					return err
				}
				c.Objects = append(c.Objects, object)
			}
		}

		token, err = d.Token()
	}

	return nil
}

// vim: ts=4
