package tmx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strconv"
)

// Tile defines a single tile in a Tileset.
type Tile struct {
	// Properties contain arbitrary key-value pairs of data to associate with the object.
	Properties
	// Rect describes the sub-rectangle representing this tile.
	Rect
	// ID is the local tile ID within its tileset.
	ID TileID
	// Class is the user-defined class of the tile. Is inherited by tile objects.
	Class string
	// Probability is the percentage indicating the probability that this tile is chosen when
	// it competes with others while editing with the terrain tool. (defaults to 0.0)
	Probability float64
	// Image is the optional image associated with this tile for image-based tilesets. For
	// tile-based tilesets, the source image is defined in the parent Tileset.
	Image *Image
	// Animation contains frames defining timings and tile IDs to produce an animation.
	Animation []Frame
	// Collision contains the map objects that define collision information for the tile, or nil
	// when none is defined.
	Collision *Collision
	// UV0 is the first texture coordinate for the tile.
	//
	// Initially calculated based on the image (or parent Tileset image) size.
	UV0 Vec2
	// UV1 is the second texture coordinate for the tile.
	//
	// Initially calculated based on the image (or parent Tileset image) size.
	UV1 Vec2
	// Tileset is a reference to the parent tilset.
	Tileset *Tileset
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (t *Tile) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			if value, err := strconv.ParseUint(attr.Value, 10, 32); err == nil {
				t.ID = TileID(value)
			} else {
				return err
			}
		case "type", "class":
			t.Class = attr.Value
		case "probability":
			if value, err := strconv.ParseFloat(attr.Value, 64); err == nil {
				t.Probability = value
			} else {
				return err
			}
		case "x":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				t.Rect.X = value
			} else {
				return err
			}
		case "y":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				t.Rect.Y = value
			} else {
				return err
			}
		case "width":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				t.Rect.Width = value
			} else {
				return err
			}
		case "height":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				t.Rect.Height = value
			} else {
				return err
			}
		case "terrain":
			logTerrain()
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
			switch child.Name.Local {
			case "properties":
				t.Properties = make(Properties)
				if err := t.Properties.UnmarshalXML(d, child); err != nil {
					return err
				}
			case "image":
				var image Image
				if err := image.UnmarshalXML(d, child); err != nil {
					return err
				}
				t.Image = &image
			case "objectgroup":
				var collision Collision
				collision.cache = t.Tileset.cache
				if err := collision.UnmarshalXML(d, child); err != nil {
					return err
				}
				t.Collision = &collision
			case "animation":
				if err := t.readFramesXML(d, child); err != nil {
					return err
				}
			default:
				logElem(child.Name.Local, start.Name.Local)
			}
		}

		token, err = d.Token()
	}

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Tile) UnmarshalJSON(data []byte) error {
	d := json.NewDecoder(bytes.NewBuffer(data))
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
				t.ID = TileID(value)
			}
		case "type":
			if t.Class, err = jsonProp[string](d); err != nil {
				return err
			}
		case "x":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				t.X = int(value)
			}
		case "y":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				t.Y = int(value)
			}
		case "width":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				t.Width = int(value)
			}
		case "height":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				t.Height = int(value)
			}
		case "animation":
			var frames []Frame
			if err := d.Decode(&frames); err != nil {
				return err
			} else {
				t.Animation = frames
			}
		case "image":
			if value, err := jsonProp[string](d); err != nil {
				return err
			} else {
				if t.Image == nil {
					t.Image = &Image{}
				}
				t.Image.Source = value
			}
		case "imagewidth":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				if t.Image == nil {
					t.Image = &Image{}
				}
				t.Image.Width = int(value)
			}
		case "imageheight":
			if value, err := jsonProp[float64](d); err != nil {
				return err
			} else {
				if t.Image == nil {
					t.Image = &Image{}
				}
				t.Image.Height = int(value)
			}
		case "objectgroup":
			var collision Collision
			if err = d.Decode(&collision); err != nil {
				return err
			}
			t.Collision = &collision
		case "probability":
			if t.Probability, err = jsonProp[float64](d); err != nil {
				return err
			}
		case "properties":
			props := make(Properties)
			if err = d.Decode(&props); err != nil {
				return err
			}
			t.Properties = props
		case "terrain":
			logTerrain()
			jsonSkip(d)
		default:
			logProp(name, "tile")
			jsonSkip(d)
		}
	}

	return nil
}

func (t *Tile) readFramesXML(d *xml.Decoder, start xml.StartElement) error {
	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		}

		if child, ok := token.(xml.StartElement); ok {
			if child.Name.Local != "frame" {
				logElem(child.Name.Local, start.Name.Local)
			} else {
				var frame Frame
				if err := frame.UnmarshalXML(d, child); err != nil {
					return err
				}
				t.Animation = append(t.Animation, frame)
			}
		}

		token, err = d.Token()
	}
	return nil
}

// vim: ts=4
