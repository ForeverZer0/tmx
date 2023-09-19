package tmx

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/DataDog/zstd"
)

// Data is a container for arbitrary data that can be stored in a TMX document.
type Data struct {
	// Compression is the compression algorithm used to deflate the payload during serialization.
	Compression Compression `xml:"compression,attr"`
	// Encoding is the encoding used to encode the payload during serialization.
	Encoding Encoding `xml:"encoding,attr"`
	// Payload is a buffer containing the data. It will be stripped of leading/trailing
	// whitespace when present, decoded from Base64, but will not be decompressed.
	Payload []byte `xml:",chardata"`
}

// TileData contains tile data defining tile layers.
type TileData struct {
	// Compression is the compression algorithm used to deflate the payload during serialization.
	Compression Compression
	// Encoding is the encoding used to encode the payload during serialization.
	Encoding Encoding
	// Chunks contains the the chunk data for infinite maps, otherwise empty.
	Chunks []Chunk
	// Tiles contains the tile definitions, or empty for infinite maps.
	Tiles []TileID
	// tileData contains the raw data from the XML/JSON. After the document is read without
	// error, it is processed and then discarded.
	tileData []byte
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (data *Data) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Must use an alias type to avoid infinite recursion.
	type dataAlias Data
	var temp dataAlias
	if err := d.DecodeElement(&temp, &start); err != nil {
		return err
	}

	data.Encoding = temp.Encoding
	data.Compression = temp.Compression
	data.Payload = trimPayload(temp.Payload)

	if data.Encoding == EncodingBase64 {
		if decoded, err := decodeBase64(data.Payload); err != nil {
			return err
		} else {
			data.Payload = decoded
		}
	}
	return nil
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (data *TileData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type xmlTileData struct {
		Compression Compression `xml:"compression,attr"`
		Encoding    Encoding    `xml:"encoding,attr"`
		Payload     []byte      `xml:",chardata"`
		Tiles       []struct {
			Value TileID `xml:"gid,attr"`
		} `xml:"tile"`
		Chunks []struct {
			Rect
			Tiles []struct {
				Value TileID `xml:"gid,attr"`
			} `xml:"tile"`
			Payload []byte `xml:",chardata"`
		} `xml:"chunk"`
	}

	var temp xmlTileData
	if err := d.DecodeElement(&temp, &start); err != nil {
		return err
	}

	data.Encoding = temp.Encoding
	data.Compression = temp.Compression

	if len(temp.Chunks) > 0 {
		data.Chunks = make([]Chunk, len(temp.Chunks))
		for i, chunk := range temp.Chunks {
			result := Chunk{Rect: chunk.Rect}
			result.Tiles = make([]TileID, chunk.Width*chunk.Height)
			if len(chunk.Tiles) > 0 {
				for i, id := range chunk.Tiles {
					result.Tiles[i] = id.Value
				}
			} else {
				// Store for now, process later
				result.tileData = trimPayload(chunk.Payload)
				// trimmed := trimPayload(chunk.Payload)
				// if err := data.decode(trimmed, result.Tiles); err != nil {
				// 	return err
				// }
			}
			data.Chunks[i] = result
		}
	} else if len(temp.Tiles) > 0 {
		data.Tiles = make([]TileID, len(temp.Tiles))
		for i, id := range temp.Tiles {
			data.Tiles[i] = id.Value
		}
	} else {
		data.tileData = trimPayload(temp.Payload)
	}

	return nil
}

func (data *TileData) decode(raw []byte, gids []TileID) error {
	// Encoding: CSV
	if data.Encoding == EncodingCSV {
		if err := decodeCSV(raw, gids); err != nil {
			return err
		}
		return nil
	}

	if data.Encoding != EncodingBase64 {
		return errInvalidEnum("Encoding", data.Encoding.String())
	}

	decoded, err := decodeBase64(raw)
	if err != nil {
		return err
	}

	var buffer []byte
	if data.Compression == CompressionNone {
		buffer = decoded
	} else {
		buffer = make([]byte, len(gids)*4)
		if err := inflate(decoded, buffer, data.Compression); err != nil {
			return err
		}
	}

	reader := bytes.NewReader(buffer)
	if err := binary.Read(reader, binary.LittleEndian, gids); err != nil && err != io.EOF {
		return err
	}
	return nil
}

// decodeCSV decodes a buffer containing comma-separated tile IDs. The input tile buffer must be
// exactly the length of the values.
func decodeCSV(data []byte, gids []TileID) error {
	ids := strings.Split(string(data), ",")
	if len(ids) != len(gids) {
		return fmt.Errorf("expected %v tile values, actual: %v", len(gids), len(ids))
	}

	for i, id := range ids {
		if result, err := strconv.ParseUint(id, 10, 32); err != nil {
			gids[i] = TileID(result)
		} else {
			return err
		}
	}
	return nil
}

// isWhitespace tests whether the given buffer contains only whitespace characters.
func (data *Data) isWhitespace(payload []byte) bool {
	for i := 0; i < len(payload); {
		if r, n := utf8.DecodeRune(payload[i:]); !unicode.IsSpace(r) {
			return false
		} else {
			i += n
		}
	}
	return true
}

// trimPayload returns the sub-slice of the paylaod with leading/trailing whitespace stripped.
func trimPayload(payload []byte) []byte {
	var start, end int
	for start = 0; start < len(payload); {
		if r, n := utf8.DecodeRune(payload[start:]); !unicode.IsSpace(r) {
			break
		} else {
			start += n
		}
	}

	for end = len(payload) - 1; end > start; {
		if r, n := utf8.DecodeLastRune(payload[:end]); !unicode.IsSpace(r) {
			break
		} else {
			end -= n
		}
	}
	return payload[start : end+1]
}

// decodeBase64 decodes the given slice of base64-encoded data into a newly allocated slice.
func decodeBase64(data []byte) ([]byte, error) {
	byteLength := base64.StdEncoding.DecodedLen(len(data))
	decoded := make([]byte, byteLength)
	decodedLen, err := base64.StdEncoding.Decode(decoded, data)
	if err != nil {
		return nil, err
	}

	return decoded[:decodedLen], nil
}

// inflate decompresses a slice of bytes into the given buffer. The destination buffer must be
// allocated with enough size to accomodate the decompressed data else it will cause a panic.
func inflate(src, dst []byte, comp Compression) error {
	var reader io.ReadCloser
	var err error

	switch comp {
	case CompressionGzip:
		reader, err = gzip.NewReader(bytes.NewReader(src))
	case CompressionZlib:
		reader, err = zlib.NewReader(bytes.NewReader(src))
	case CompressionZstd:
		reader = zstd.NewReader(bytes.NewReader(src))
		err = nil
	case CompressionNone:
		// This branch isn't possible in practice, but included here Just-In-Case™
		copy(dst, src)
		return nil
	default:
		return errInvalidEnum("Compression", comp.String())
	}

	if reader != nil {
		defer reader.Close()
	}

	if err != nil {
		return err
	}

	if n, err := reader.Read(dst); err != nil && err != io.EOF {
		return err
	} else if n != len(dst) {
		return errors.New("failed to read correct number of bytes")
	}

	return nil
}

// vim: ts=4
