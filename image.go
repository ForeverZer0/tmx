package tmx

import (
	"encoding/xml"
	"image"
	"strconv"
)

// ImageCallback is a function that can be assigned a callback function, which will be called each
// time an Image is parsed from a TMX document. This provides an opportunity to prepare resources,
// such as image loading, decoding of embedded data, texture creation, etc.
//
// The parsed Image instance is supplied to the callback.
var ImageCallback func(image *Image)

// Image describes a graphic resource used by tilesets, objects, etc.
//
// Images typically point to an external file where the resource is located. Although the Tiled
// editor does not directly support it, the TMX format does support image data being embedded into
// the file as raw bytes to be decoded by the application.
//
// For handling image loading, decoding, caching, etc. on the fly, see ImageCallback.
type Image struct {
	// Format describes the image type for embedded images.
	// Valid values are file extensions like png, gif, jpg, bmp, etc.
	Format string
	// Source is the reference to the tileset image file. Only used if the image is not embedded.
	Source string
	// Transparency defines a specific color that is treated as transparent.
	Transparency Color
	// Size describes the dimensions of the image in pixel units (optional).
	Size Size
	// Data contains the payload of an embedded image. This is not supported by the Tiled editor,
	// but is by the TMX specification.
	Data *Data
	// UserID provides a field that can be used to store a value within the Image instance, such
	// as an OpenGL texture.
	//
	// This library does not touch the field, and it is the responsiblity of the user to manage
	// any required freeing of unmanaged resources before the Image is garbage collected and the
	// reference lost.
	UserID uint32
	// UserImage provides a field to that can be used to store a decoded image with the Image
	// instance.
	UserImage image.Image
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (img *Image) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "format":
			img.Format = attr.Value
		case "source":
			img.Source = attr.Value
		case "trans":
			if color, err := ParseColor(attr.Value); err != nil {
				return err
			} else {
				img.Transparency = color
			}
		case "width":
			if value, err := strconv.ParseUint(attr.Value, 10, 32); err != nil {
				return err
			} else {
				img.Size.Width = int(value)
			}
		case "height":
			if value, err := strconv.ParseUint(attr.Value, 10, 32); err != nil {
				return err
			} else {
				img.Size.Height = int(value)
			}
		case "id": // Ignore, deprecated legacy Java filth
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
			if child.Name.Local != "data" {
				logElem(child.Name.Local, start.Name.Local)
			} else {
				var data Data
				if err := data.UnmarshalXML(d, child); err != nil {
					return err
				}
				img.Data = &data
			}
		}

		token, err = d.Token()
	}

	if ImageCallback != nil {
		ImageCallback(img)
	}

	return nil
}

// vim: ts=4
