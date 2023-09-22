package tmx

import "encoding/xml"

// TileLayer describes a map layer that is composed of tile data from a Tileset.
type TileLayer struct {
	baseLayer
	// TileLayer contains the tile/chunk data for the layer.
	TileData
	// ChunkSize is the total size of the layer in tile units before it "wraps".
	//
	// Only valid for infinite maps.
	ChunkSize Size

	// chunkRows is the number of rows of chunks.
	chunkRows int
	// chunkCols is the number of columns of chunks.
	chunkCols int
	// chunkSz is the size of an individual chunk.
	chunkSz Size
}

// GetGID returns a the global tile ID for the specified map coordinates.
//
// For infinte maps, the given position is unrestricted and can include negative values,
// otherwise it must be within the bounds of the map. A zero value will be returned
// for positions outside the map bounds or when no tile is defined at the given position.
func (layer *TileLayer) GetGID(x, y int) TileID {
	if len(layer.Chunks) > 0 {
		chunk, x, y := layer.ChunkAt(x, y)
		return chunk.Tiles[x+(y*chunk.Width)]
	} else if x < 0 || x >= layer.Width || y < 0 || y >= layer.Height {
		return 0
	}

	return layer.Tiles[x+(y*layer.Width)]
}

// TileAt returns the tile and the GID (with flip/rotate bits still set) at the
// specified map coordinates.
//
// For infinte maps, the given position is unrestricted and can include negative values,
// otherwise it must be within the bounds of the map. A nil value will be returned
// for positions outside the map bounds or when no tile is defined at the given position.
func (layer *TileLayer) TileAt(x, y int) (*Tile, TileID) {
	if gid := layer.GetGID(x, y); gid != 0 {
		if ts, id := layer.parent.Tileset(gid); id > 0 {
			return &ts.Tiles[id], gid
		}
	}
	return nil, 0
}

// ChunkAt returns the chunk the Chunk and localized coordinates for the
// given position. The given values can be positive or negative.
//
// Only valid for infinte maps, otherwise returns nil.
func (layer *TileLayer) ChunkAt(x, y int) (*Chunk, int, int) {
	if len(layer.Chunks) == 0 {
		return nil, 0, 0
	}

	// Normalize coordinates to the chunk boundaries
	x %= layer.ChunkSize.Width
	y %= layer.ChunkSize.Height
	if x < 0 {
		x += layer.ChunkSize.Width
	}
	if y < 0 {
		y += layer.ChunkSize.Height
	}

	// Calculate chunk index
	i := (x / layer.ChunkSize.Width) + ((y * layer.ChunkSize.Height) * layer.chunkCols)
	return &layer.Chunks[i], x % layer.chunkSz.Width, y % layer.chunkSz.Height
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

	if err := layer.postProcess(layer.Area()); err != nil {
		return err
	}

	if len(layer.Chunks) > 0 {
		last := layer.Chunks[len(layer.Chunks)-1]
		layer.chunkSz = last.Size
		layer.ChunkSize = Size{Width: last.Right(), Height: last.Bottom()}

		layer.chunkCols = layer.ChunkSize.Width / layer.chunkSz.Width
		layer.chunkRows = layer.ChunkSize.Height / layer.chunkSz.Height
	}

	return nil
}

// vim: ts=4
