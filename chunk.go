package tmx

import (
	"bytes"
	"encoding/json"
	"errors"
)

type Chunk struct {
	Rect
	Tiles    []TileID
	tileData []byte
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *Chunk) UnmarshalJSON(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	token, err := d.Token()
	if err != nil {
		return err
	} else if token != json.Delim('{') {
		return errors.New("expected JSON object")
	}

	for {
		if token, err = d.Token(); err != nil {
			return err
		} else if token == json.Delim('}') {
			break
		}

		name := token.(string)
		if token, err = d.Token(); err != nil {
			return err
		}

		switch name {
		case "x":
			c.X = int(token.(float64))
		case "y":
			c.Y = int(token.(float64))
		case "width":
			c.Width = int(token.(float64))
		case "height":
			c.Height = int(token.(float64))
		case "data":
			if token == json.Delim('[') {
				// An array of tile IDs
				c.Tiles = make([]TileID, 0, c.Width*c.Height)
				for {
					if token, err = d.Token(); err != nil {
						return err
					} else if token == json.Delim(']') {
						break
					}
					id := TileID(token.(float64))
					c.Tiles = append(c.Tiles, id)
				}
			} else {
				// Text data. Store for now, process later
				c.tileData = trimPayload([]byte(token.(string)))
			}
		default:
			logProp(name, "chunk")
			jsonSkip(d)
		}
	}

	return nil
}

// vim: ts=4
