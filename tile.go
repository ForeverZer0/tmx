package tmx

import (
	"encoding/json"
	"encoding/xml"
	"log"
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
	// Image is the optional image associated with this tile.
	Image *Image
	// Animation contains frames defining timings and tile IDs to produce an animation.
	Animation []Frame
	// Collision contains the map objects that define collision information for the tile, or nil
	// when none is defined.
	Collision *Collision
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
			log.Println("terrains are no longer supported, and are replaced by wangsets")
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

	type jsonTile struct {
		ID          TileID     `json:"id"`
		Class       string     `json:"type"`
		X           int        `json:"x"`
		Y           int        `json:"y"`
		Width       int        `json:"width"`
		Height      int        `json:"height"`
		Animation   []Frame    `json:"animation"`
		Image       string     `json:"image"`
		ImageHeight int        `json:"imageheight"`
		ImageWidth  int        `json:"imagewidth"`
		Collision   *Collision `json:"objectgroup"`
		Probability float64    `json:"probability"`
		Properties  Properties `json:"properties"`
		// Terrain     []int      `json:"terrain"`
	}
	
	var temp jsonTile
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	t.ID = temp.ID
	t.Class = temp.Class
	t.Rect = Rect{Point{X: temp.X, Y: temp.Y}, Size{Width: temp.Width, Height: temp.Height}}
	t.Animation = temp.Animation
	t.Image = &Image{Source: temp.Image, Size: Size{Width: temp.ImageWidth, Height: temp.ImageHeight}}
	t.Probability = temp.Probability
	t.Collision = temp.Collision
	t.Properties = temp.Properties
	
	// TODO: Terrain

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
