package tmx

import (
	"bytes"
	"encoding/json"
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
	// cache is a resource cache that maintains references to shared objects.
	cache *Cache
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

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *Collision) UnmarshalJSON(data []byte) error {
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
				c.ID = int(value)
			}
		case "draworder":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if value, err := parseDrawOrder(str); err != nil {
				return err
			} else {
				c.DrawOrder = value
			}
		case "objects":
			if token, err = d.Token(); err != nil {
				return err
			} else if token != json.Delim('[') {
				return ErrExpectedArray
			}
			for d.More() {
				var obj Object
				obj.cache = c.cache
				if err = d.Decode(&obj); err != nil {
					return err
				}
				c.Objects = append(c.Objects, obj)
			}
			// Position to next token ']'
			if token, err = d.Token(); err != nil {
				return err
			}
		default:
			logProp(name, "objectgroup")
			jsonSkip(d)
		}
	}

	return nil
}

// vim: ts=4
