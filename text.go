package tmx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"strconv"
)

type FontStyle uint8

const (
	// StyleBold indicates bold font style.
	StyleBold FontStyle = 1 << iota
	// StyleItalic indicates italic font style.
	StyleItalic
	// StyleUnderline indicates underline font style.
	StyleUnderline
	// StyleStrikeout indicates strikeout font style.
	StyleStrikeout
	// StyleKerning indicates if kerning should be used when rendering the text.
	StyleKerning
)

type Text struct {
	// FontFamily is the font family used to render text. Defaults to "sans-serif".
	FontFamily string
	// Value is the text to be rendered.
	Value string
	// PixelSize is the size of the font in pixel units. Defaults to 16.
	PixelSize int
	// Color is the color of the font used to render the text.
	Color Color
	// Style is a bitfield containing flags on how the font should be styled.
	Style FontStyle
	// WordWrap indicates if word wrapping is enabled.
	WordWrap bool
	// Align describes how the alignment of the rendered text.
	Align Align
	// flags are used internally to track which fields were explicitly defined.
	flags setFlags
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (obj *Text) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	obj.FontFamily = "sans-serif"
	obj.PixelSize = 16
	obj.Color = 0xFF000000
	obj.Style = StyleKerning

	hAlign := AlignLeft
	vAlign := AlignTop

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "fontfamily":
			obj.FontFamily = attr.Value
			obj.flags |= flagFont
		case "pixelsize":
			if value, err := strconv.ParseUint(attr.Value, 10, 32); err == nil {
				obj.PixelSize = int(value)
				obj.flags |= flagFontSize
			} else {
				return err
			}
		case "wrap":
			if value, err := strconv.ParseBool(attr.Value); err == nil {
				obj.WordWrap = value
				obj.flags |= flagTextWrap
			} else {
				return err
			}
		case "color":
			if value, err := ParseColor(attr.Value); err == nil {
				obj.Color = value
				obj.flags |= flagTextColor
			} else {
				return err
			}
		case "bold":
			if value, err := strconv.ParseBool(attr.Value); err == nil {
				if value {
					obj.Style |= StyleBold
				} else {
					obj.Style &= ^StyleBold
				}
				obj.flags |= flagBold
			} else {
				return err
			}
		case "italic":
			if value, err := strconv.ParseBool(attr.Value); err == nil {
				if value {
					obj.Style |= StyleItalic
				} else {
					obj.Style &= ^StyleItalic
				}
				obj.flags |= flagItalic
			} else {
				return err
			}
		case "underline":
			if value, err := strconv.ParseBool(attr.Value); err == nil {
				if value {
					obj.Style |= StyleUnderline
				} else {
					obj.Style &= ^StyleUnderline
				}
				obj.flags |= flagUnderline
			} else {
				return err
			}
		case "strikeout":
			if value, err := strconv.ParseBool(attr.Value); err == nil {
				if value {
					obj.Style |= StyleStrikeout
				} else {
					obj.Style &= ^StyleStrikeout
				}
				obj.flags |= flagStrikeout
			} else {
				return err
			}
		case "kerning":
			if value, err := strconv.ParseBool(attr.Value); err == nil {
				if value {
					obj.Style |= StyleKerning
				} else {
					obj.Style &= ^StyleKerning
				}
				obj.flags |= flagKerning
			} else {
				return err
			}
		case "halign":
			if value, err := parseAlign(attr.Value); err != nil {
				return err
			} else {
				if value == AlignCenter {
					value = AlignCenterH
				}
				hAlign |= value
				obj.flags |= flagHAlign
			}
		case "valign":
			if value, err := parseAlign(attr.Value); err != nil {
				return err
			} else {
				if value == AlignCenter {
					value = AlignCenterV
				}
				vAlign |= value
				obj.flags |= flagVAlign
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

		if next, ok := token.(xml.StartElement); ok {
			logElem(next.Name.Local, start.Name.Local)
		} else if data, ok := token.(xml.CharData); ok {
			obj.Value = string(data)
			obj.flags |= flagText
		}
		token, err = d.Token()
	}

	obj.Align = hAlign | vAlign
	return nil
}

func (obj *Text) UnmarshalJSON(data []byte) error {
	obj.FontFamily = "sans-serif"
	obj.PixelSize = 16
	obj.Color = 0xFF000000
	obj.Style = StyleKerning

	hAlign := AlignLeft
	vAlign := AlignTop

	buffer := bytes.NewReader(data)
	d := json.NewDecoder(buffer)

	token, err := d.Token()
	if err != nil {
		return err
	}
	if token != json.Delim('{') {
		return errors.New("expected JSON object")
	}

	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == json.Delim('}') {
			break
		}

		name := token.(string)
		if token, err = d.Token(); err != nil {
			return err
		}

		switch name {
		case "pixelsize":
			obj.PixelSize = int(token.(float64))
			obj.flags |= flagFontSize
		case "text":
			obj.Value = token.(string)
			obj.flags |= flagText
		case "fontfamily":
			obj.FontFamily = token.(string)
			obj.flags |= flagFont
		case "wrap":
			obj.WordWrap = token.(bool)
			obj.flags |= flagTextWrap
		case "bold":
			if token.(bool) {
				obj.Style |= StyleBold
			} else {
				obj.Style &= ^StyleBold
			}
			obj.flags |= flagBold
		case "italic":
			if token.(bool) {
				obj.Style |= StyleItalic
			} else {
				obj.Style &= ^StyleItalic
			}
			obj.flags |= flagItalic
		case "underline":
			if token.(bool) {
				obj.Style |= StyleItalic
			} else {
				obj.Style &= ^StyleItalic
			}
			obj.flags |= flagUnderline
		case "strikeout":
			if token.(bool) {
				obj.Style |= StyleStrikeout
			} else {
				obj.Style &= ^StyleStrikeout
			}
			obj.flags |= flagStrikeout
		case "kerning":
			if token.(bool) {
				obj.Style |= StyleKerning
			} else {
				obj.Style &= ^StyleKerning
			}
			obj.flags |= flagKerning
		case "color":
			if color, err := ParseColor(token.(string)); err != nil {
				return err
			} else {
				obj.Color = color
				obj.flags |= flagTextColor
			}
		case "halign":
			if value, err := parseAlign(token.(string)); err != nil {
				return err
			} else {
				if value == AlignCenter {
					value = AlignCenterH
				}
				hAlign |= value
				obj.flags |= flagHAlign
			}
		case "valign":
			if value, err := parseAlign(token.(string)); err != nil {
				return err
			} else {
				if value == AlignCenter {
					value = AlignCenterV
				}
				vAlign |= value
				obj.flags |= flagVAlign
			}
		}
	}

	obj.Align = hAlign | vAlign
	return nil
}

// vim: ts=4
