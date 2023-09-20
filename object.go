package tmx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type setFlags uint32

const (
	flagName setFlags = 1 << iota
	flagClass
	flagX
	flagY
	flagWidth
	flagHeight
	flagRotation
	flagGID
	flagVisible
	flagKind
	flagPoints
	// Text flags
	flagFont
	flagFontSize
	flagTextWrap
	flagTextColor
	flagBold
	flagItalic
	flagUnderline
	flagStrikeout
	flagKerning
	flagHAlign
	flagVAlign
	flagText
)

// Object is an arbitray entity that can be placed on the map, or even invisible to define
// regions, events, etc.
type Object struct {
	// Properties contain arbitrary key-value pairs of data to associate with the object.
	Properties
	// ID is a unique identifier for the object. Valid IDs are greater than 0.
	ID int
	// Name is a user-defined name for the object.
	Name string
	// Class is a user-defined class for the object.
	Class string
	// Location is the coordinates of the object in pixel units.
	Location Vec2
	// Size is the dimensions of the object in pixel units.
	Size Vec2
	// Rotation is the amount of rotation of the object in degrees clockwise around its Location.
	Rotation float32
	// GID is a reference to a global tile ID (optional).
	GID TileID
	// Visible determines whenther the object is shown or not.
	Visible bool
	// Template is a reference to a template from which the object inherits its values.
	Template *Template
	// Type describes any object specialization this instance may have.
	Type ObjectType
	// Text is the definition used for Text objects, otherwise nil.
	Text *Text
	// Points is a list of vectors used for Polygon and Polyline types.
	Points []Vec2
	// flags is used internally to determine which values were explicitly set, and which are
	// simply the "zero" or default value. This is necessary for determining how to inherit
	// from a template object, as it would otherwise be impossible to determine if a value of
	// 0, false, "", etc. should be inherited, or it is merely a default.
	flags setFlags
	
	cache *Cache
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (obj *Object) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	obj.Visible = true

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			if value, err := strconv.Atoi(attr.Value); err == nil {
				obj.ID = value
			} else {
				return err
			}
		case "name":
			obj.Name = attr.Value
			obj.flags |= flagName
		case "type", "class":
			obj.Class = attr.Value
			obj.flags |= flagClass
		case "x":
			if value, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				obj.Location.X = float32(value)
				obj.flags |= flagX
			} else {
				return err
			}
		case "y":
			if value, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				obj.Location.Y = float32(value)
				obj.flags |= flagY
			} else {
				return err
			}
		case "width":
			if value, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				obj.Size.X = float32(value)
				obj.flags |= flagWidth
			} else {
				return err
			}
		case "height":
			if value, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				obj.Size.Y = float32(value)
				obj.flags |= flagHeight
			} else {
				return err
			}
		case "rotation":
			if value, err := strconv.ParseFloat(attr.Value, 32); err == nil {
				obj.Rotation = float32(value)
				obj.flags |= flagRotation
			} else {
				return err
			}
		case "gid":
			if value, err := strconv.ParseUint(attr.Value, 10, 32); err == nil {
				obj.GID = TileID(value)
				obj.flags |= flagGID
			} else {
				return err
			}
		case "visible":
			if value, err := strconv.ParseBool(attr.Value); err == nil {
				obj.Visible = value
				obj.flags |= flagVisible
			} else {
				return err
			}
		case "template":
			if tmpl, err := OpenTemplate(attr.Value, obj.cache); err == nil {
				obj.Template = tmpl
			} else {
				return err
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

		next, ok := token.(xml.StartElement)
		if ok {
			switch next.Name.Local {

			case "properties":
				obj.Properties = make(Properties)
				if err := obj.Properties.UnmarshalXML(d, next); err != nil {
					return err
				}
			case "point":
				obj.Type = ObjectPoint
				obj.flags |= flagKind
			case "ellipse":
				obj.Type = ObjectEllipse
				obj.flags |= flagKind
			case "polygon":
				if obj.Points, err = parsePoints(next); err != nil {
					return err
				}
				obj.Type = ObjectPolygon
				obj.flags |= flagPoints | flagKind
			case "polyline":
				if obj.Points, err = parsePoints(next); err != nil {
					return err
				}
				obj.Type = ObjectPolyline
				obj.flags |= flagPoints | flagKind
			case "text":
				var text Text
				if err := text.UnmarshalXML(d, next); err != nil {
					return err
				}
				obj.Type = ObjectText
				obj.Text = &text
				// Merge the flags from the text object
				obj.flags |= text.flags
			default:
				logElem(next.Name.Local, start.Name.Local)
			}
		}

		token, err = d.Token()
	}

	obj.inherit()
	return nil
}

func (obj *Object) UnmarshalJSON(data []byte) error {
	buffer := bytes.NewReader(data)
	d := json.NewDecoder(buffer)
	obj.Visible = true

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
		switch name {
		case "properties":
			// TODO !
			var props []Property
			if err := d.Decode(&props); err != nil {
				return err
			}
			obj.Properties = make(Properties)
			for _, prop := range props {
				obj.Properties[prop.Name] = prop
			}
			continue
		case "polyline":
			if err := d.Decode(&obj.Points); err != nil {
				return err
			}
			obj.Type = ObjectPolyline
			continue
		case "polygon":
			if err := d.Decode(&obj.Points); err != nil {
				return err
			}
			obj.Type = ObjectPolygon
			continue
		case "text":
			var text Text
			if err := d.Decode(&text); err != nil {
				return err
			}
			obj.Text = &text
			obj.Type = ObjectText
			obj.flags |= text.flags
			continue
		default:
			// Set the next token for all other properties
			if token, err = d.Token(); err != nil {
				return err
			}
		}

		switch name {
		case "id":
			obj.ID = int(token.(float64))
		case "name":
			obj.Name = token.(string)
			obj.flags |= flagName
		case "gid":
			obj.GID = TileID(token.(float64))
			obj.flags |= flagGID
		case "x":
			obj.Location.X = float32(token.(float64))
			obj.flags |= flagX
		case "y":
			obj.Location.Y = float32(token.(float64))
			obj.flags |= flagY
		case "width":
			obj.Size.X = float32(token.(float64))
			obj.flags |= flagWidth
		case "height":
			obj.Size.Y = float32(token.(float64))
			obj.flags |= flagHeight
		case "rotation":
			obj.Rotation = float32(token.(float64))
			obj.flags |= flagRotation
		case "type", "class":
			obj.Class = token.(string)
			obj.flags |= flagClass
		case "visible":
			obj.Visible = token.(bool)
			obj.flags |= flagVisible
		case "template":
			if tmpl, err := OpenTemplate(token.(string), obj.cache); err == nil {
				obj.Template = tmpl
			} else {
				return err
			}
		case "point":
			obj.Type = ObjectPoint
		case "ellipse":
			obj.Type = ObjectEllipse
		}
	}

	obj.inherit()
	return nil
}

// inherit applies any required changes to an object that uses a Template.
func (obj *Object) inherit() {
	if obj.Template == nil {
		return
	}

	tmp := obj.Template
	if obj.override(flagKind) {
		obj.Type = tmp.Type
	}
	if obj.override(flagName) {
		obj.Name = tmp.Name
	}
	if obj.override(flagClass) {
		obj.Class = tmp.Class
	}
	if obj.override(flagX) {
		obj.Location.X = tmp.Location.X
	}
	if obj.override(flagY) {
		obj.Location.Y = tmp.Location.Y
	}
	if obj.override(flagWidth) {
		obj.Size.X = tmp.Size.X
	}
	if obj.override(flagHeight) {
		obj.Size.Y = tmp.Size.Y
	}
	if obj.override(flagRotation) {
		obj.Rotation = tmp.Rotation
	}
	if obj.override(flagGID) {
		obj.GID = tmp.GID
	}
	if obj.override(flagVisible) {
		obj.Visible = tmp.Visible
	}
	if obj.override(flagPoints) {
		obj.Points = make([]Vec2, 0, len(tmp.Points))
		copy(obj.Points, tmp.Points)
	}

	if tmp.Text != nil {

		if obj.Text == nil {
			obj.Text = &Text{}
		}

		if obj.override(flagFont) {
			obj.Text.FontFamily = tmp.Text.FontFamily
		}
		if obj.override(flagFontSize) {
			obj.Text.PixelSize = tmp.Text.PixelSize
		}
		if obj.override(flagTextWrap) {
			obj.Text.WordWrap = tmp.Text.WordWrap
		}
		if obj.override(flagTextColor) {
			obj.Text.Color = tmp.Text.Color
		}
		if obj.override(flagBold) {
			if tmp.Text.Style&StyleBold != 0 {
				obj.Text.Style |= StyleBold
			} else {
				obj.Text.Style &= ^StyleBold
			}
		}
		if obj.override(flagItalic) {
			if tmp.Text.Style&StyleItalic != 0 {
				obj.Text.Style |= StyleItalic
			} else {
				obj.Text.Style &= ^StyleItalic
			}
		}
		if obj.override(flagUnderline) {
			if tmp.Text.Style&StyleUnderline != 0 {
				obj.Text.Style |= StyleUnderline
			} else {
				obj.Text.Style &= ^StyleUnderline
			}
		}
		if obj.override(flagStrikeout) {
			if tmp.Text.Style&StyleStrikeout != 0 {
				obj.Text.Style |= StyleStrikeout
			} else {
				obj.Text.Style &= ^StyleStrikeout
			}
		}
		if obj.override(flagKerning) {
			if tmp.Text.Style&StyleKerning != 0 {
				obj.Text.Style |= StyleKerning
			} else {
				obj.Text.Style &= ^StyleKerning
			}
		}
		if obj.override(flagHAlign) {
			obj.Text.Align &= clearHorizontal
			obj.Text.Align |= tmp.Text.Align & clearVertical
		}
		if obj.override(flagVAlign) {
			obj.Text.Align &= clearVertical
			obj.Text.Align |= tmp.Text.Align & clearHorizontal
		}
		if obj.override(flagText) {
			obj.Text = tmp.Text
		}
	}

	if tmp.Properties != nil {
		if obj.Properties == nil {
			obj.Properties = make(Properties, len(tmp.Properties))
		}
		obj.Properties.Merge(tmp.Properties, false)
	}
}

// override tests whether the field described by the given flag should be overriden by its
// Template instance.
func (obj *Object) override(flag setFlags) bool {
	// The object set its own value, do not override
	if obj.flags&flag != 0 {
		return false
	}
	// The template did explicitly set a value, override the object's value
	if obj.Template.flags&flag != 0 {
		return true
	}
	// Neither defined the value, let it be the default
	return false
}

func parsePoints(element xml.StartElement) ([]Vec2, error) {
	var points []Vec2
	for _, attr := range element.Attr {

		if attr.Name.Local != "points" {
			logAttr(attr.Name.Local, element.Name.Local)
			continue
		}

		parts := strings.Split(attr.Value, " ")
		points = make([]Vec2, len(parts))
		for i, part := range parts {

			xy := strings.Split(part, ",")
			if len(xy) != 2 {
				return nil, fmt.Errorf(`cannot parse "%s" as Vec2`, part)
			}

			var vec Vec2
			if x, err := strconv.ParseFloat(xy[0], 32); err == nil {
				vec.X = float32(x)
			} else {
				return nil, err
			}

			if y, err := strconv.ParseFloat(xy[1], 32); err == nil {
				vec.Y = float32(y)
			} else {
				return nil, err
			}

			points[i] = vec
		}
	}
	return points, nil
}

// Clone creates a deep copy of the Object.
func (obj *Object) Clone() *Object {
	dup := *obj
	if obj.Text != nil {
		text := *obj.Text
		dup.Text = &text
	}
	if len(obj.Points) > 0 {
		dup.Points = make([]Vec2, 0, len(obj.Points))
		copy(dup.Points, obj.Points)
	}

	if obj.Properties != nil {
		dup.Properties = obj.Properties.Clone()
	}
	return &dup
}

// vim: ts=4
