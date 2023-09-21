package tmx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
)

// Property describes an arbitrary named value associated with an object.
type Property struct {
	// Name is the user-defined name of the Property. Names are used as the key within
	// the parent map, and are therefore unique within any set of Properties.
	Name string
	// Type describes the data type of the property Value.
	Type DataType
	// Class is the name of the user-defined class of the property (optional).
	Class string
	// Value is the untyped value of the property.
	Value interface{}
}

// String implements the Stringer interface.
func (p Property) String() string {
	var value string
	if str, ok := p.Value.(string); ok {
		value = strconv.Quote(str)
	} else {
		value = fmt.Sprint(p.Value)
	}

	if p.Name != "" {
		return fmt.Sprintf(`"%s":%v`, p.Name, value)
	}
	return fmt.Sprint(p.Value)
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (p *Property) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "name":
			p.Name = attr.Value
		case "type":
			if value, err := parseDataType(attr.Value); err != nil {
				return err
			} else {
				p.Type = value
			}
		case "propertytype":
			p.Class = attr.Value
		case "value", "default":
			switch p.Type {
			case TypeBool:
				if value, err := strconv.ParseBool(attr.Value); err != nil {
					return err
				} else {
					p.Value = value
				}
			case TypeInt, TypeObject:
				if value, err := strconv.Atoi(attr.Value); err != nil {
					return err
				} else {
					p.Value = value
				}
			case TypeFloat:
				if value, err := strconv.ParseFloat(attr.Value, 64); err != nil {
					return err
				} else {
					p.Value = value
				}
			case TypeColor:
				if value, err := ParseColor(attr.Value); err != nil {
					return err
				} else {
					p.Value = value
				}
			case TypeClass:
				// Value is a child element, and processed below
			default:
				p.Value = attr.Value
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
			if child.Name.Local != "properties" {
				logElem(child.Name.Local, start.Name.Local)
			} else {
				var props Properties
				// Initialize to default class if defined...
				if CustomTypes != nil && p.Class != "" {
					if base, ok := CustomTypes[p.Class]; ok {
						props = base.Members.Clone()
					}
				}
				// ...else just start with a blank map
				if props == nil {
					props = make(Properties)
				}
				// Default values will get overwritten if defined
				if err = props.UnmarshalXML(d, child); err != nil {
					return err
				}
				p.Value = props
			}
		}

		token, err = d.Token()
	}

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (p *Property) UnmarshalJSON(data []byte) error {
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
		case "name":
			if p.Name, err = jsonProp[string](d); err != nil {
				return err
			}
		case "type":
			if str, err := jsonProp[string](d); err != nil {
				return err
			} else if value, err := parseDataType(str); err != nil {
				return err
			} else {
				p.Type = value
			}
		case "propertytype", "propertyType":
			if p.Class, err = jsonProp[string](d); err != nil {
				return err
			}
		case "value", "default":
			if value, err := p.jsonValue(d); err != nil {
				return err
			} else {
				p.Value = value
			}
		default:
			logProp(name, "property")
			jsonSkip(d)
		}
	}

	return nil
}

// jsonValue decodes a value from the current position in the stream.
func (p Property) jsonValue(d *json.Decoder) (interface{}, error) {
	if p.Type == TypeClass {
		return p.jsonClass(d, p.Class)
	}

	token, err := d.Token()
	if err != nil {
		return nil, err
	}

	switch value := token.(type) {
	case float64:
		switch p.Type {
		case TypeInt, TypeObject:
			return int(value), nil
		default:
			return value, nil
		}
	case bool:
		return value, nil
	case string:
		if p.Type == TypeColor {
			if color, err := ParseColor(value); err != nil {
				return nil, err
			} else {
				return color, nil
			}
		}
		return value, nil
	default:
		return nil, errors.New("unknown property type")
	}
}

// jsonClass creates a set of Properties for a custom class type.
//
// Unfortunately the JSON format does not indicate the data-types or class-names as
// the XML format does, and it is given as a simple list of key-value pairs. While
// every effort is made to set the correct type, this makes it impossible to
// differentiate between certain values. For example, a float value may be written
// as "0". Because there is no indicator for what type it is, it will be assumed to
// be an integer, although that was not the case.
//
// https://github.com/mapeditor/tiled/issues/3820
func (p Property) jsonClass(d *json.Decoder, class string) (Properties, error) {
	var props Properties
	// Initialize to default class if defined...
	if CustomTypes != nil && class != "" {
		if base, ok := CustomTypes[class]; ok {
			props = base.Members.Clone()
		}
	}
	// ...else just start with a blank map
	if props == nil {
		props = make(Properties)
	}

	token, err := d.Token()
	if err != nil {
		return nil, err
	} else if token != json.Delim('{') {
		return nil, ErrExpectedObject
	}

	for {
		token, err = d.Token()
		if err != nil {
			return nil, err
		} else if token == json.Delim('}') {
			break
		}

		name := token.(string)
		var prop Property

		// Get value based on existing property (CustomClass)
		if base, ok := props[name]; ok {
			prop = base
			if value, err := prop.jsonValue(d); err != nil {
				return nil, err
			} else {
				prop.Value = value
			}
		} else {
			// Get value without a known type
			prop.Type = -1
			prop.Value, err = prop.jsonValue(d)
			if err != nil {
				return nil, err
			}

			// Reasonable attempt at guessing the correct DataType
			// based on heursitics. 
			switch v := prop.Value.(type) {
			case string:
				if color, err := ParseColor(v); err == nil {
					prop.Type = TypeColor
					prop.Value = color
				} else {
					prop.Type = TypeString
				}
			case int:
				prop.Type = TypeInt
			case Color:
				prop.Type = TypeColor
			case float64:
				// Truncate float value and see if still equal to original value.
				// If so, just assume it is an integer. This might be the correct
				// choice for a bit over 50% of the time, maybe not, it's all guessing
				// at this point... Hopefully this gets fixed upstream in the format.
				if float64(int(v)) == v {
					prop.Type = TypeInt
					prop.Value = int(v)
				} else {
					prop.Type = TypeFloat
				}
			case bool:
				prop.Type = TypeBool
			case Properties:
				prop.Type = TypeClass
			}
		}
		
		props[prop.Name] = prop
	}

	return props, nil
}

// Clone implements the Cloner interface.
func (p Property) Clone() Property {
	dup := p
	if class, ok := p.Value.(Properties); ok {
		dup.Value = class.Clone()
	}
	return dup
}

// vim: ts=4
