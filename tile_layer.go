package tmx

import (
	"encoding/xml"
	"fmt"
)

type TileLayer struct {
	baseLayer
	TileData
}

func (layer *TileLayer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		if handled, err := layer.baseLayer.xmlAttr(attr); err != nil {
			return err
		} else if handled {
			continue
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
						if err := d.DecodeElement(&layer.TileData, &start); err != nil {
							return err
						}
					default:
						logElem(child.Name.Local, start.Name.Local)
					}
				}
			}
			token, err = d.Token()
		}
	}

	return layer.processTiles()
}

func (layer *TileLayer) processTiles() error {
	if len(layer.Chunks) > 0 {
		for i, chunk := range layer.Chunks {
			area := chunk.Width * chunk.Height
			if len(chunk.tileData) == 0 {
				if len(chunk.Tiles) != area {
					return fmt.Errorf("not enough tiles in chunk [%d, %d]", chunk.X, chunk.Y)
				}
				continue
			}

			layer.Chunks[i].Tiles = make([]TileID, area)
			if err := layer.decode(chunk.tileData, chunk.Tiles); err != nil {
				return err
			}
			layer.Chunks[i].tileData = nil
		}

		return nil
	}

	area := layer.Width * layer.Height
	if len(layer.Tiles) > 0 {
		if len(layer.Tiles) != area {
			return fmt.Errorf(`not enough tiles in tile layer "%s"`, layer.Name)
		} else {
			return nil
		}
	}

	layer.Tiles = make([]TileID, area)
	if err := layer.decode(layer.tileData, layer.Tiles); err != nil {
		return err
	} else {
		layer.tileData = nil
	}

	return nil
}

// vim: ts=4
