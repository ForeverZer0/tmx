package tmx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
)

type Property struct {
	Name  string
	Type  DataType
	Class string
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
		case "value":
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
				props := make(Properties)
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
		return errors.New("expected JSON object")
	}

	for {
		if token, err = d.Token(); err != nil {
			return err
		} else if token == json.Delim('}') {
			break
		}

		name, ok := token.(string)
		if !ok {
			continue
		}

		// Depending on value type, we don't always want to consume the token yet
		if name != "value" {
			if token, err = d.Token(); err != nil {
				return err
			}
		}

		switch name {
		case "name":
			if value, ok := token.(string); ok {
				p.Name = value
			} else {
				// TODO
			}
		case "type":
			if value, ok := token.(string); ok {
				if t, err := parseDataType(value); err != nil {
					return err
				} else {
					p.Type = t
				}
			} else {
				// TODO
			}
		case "propertytype":
			if value, ok := token.(string); ok {
				p.Class = value
			} else {
				// TODO
			}
		case "value":
			if value, err := p.jsonValue(d, p.Type); err != nil {
				return err
			} else {
				p.Value = value
			}
		}
	}

	return nil
}

func (p Property) jsonValue(d *json.Decoder, dt DataType) (interface{}, error) {
	if dt == TypeClass {
		return p.jsonClass(d)
	}

	token, err := d.Token()
	if err != nil {
		return nil, err
	}

	switch value := token.(type) {
	case float64:
		switch dt {
		case TypeInt, TypeObject:
			return int(value), nil
		default:
			return value, nil
		}
	case bool:
		return value, nil
	case string:
		if dt == TypeColor {
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

func (p Property) jsonClass(d *json.Decoder) (Properties, error) {
	props := make(Properties)

	if token, err := d.Token(); err != nil {
		return nil, err
	} else if token != json.Delim('{') {
		return nil, errors.New("expected JSON object")
	}

	for {
		token, err := d.Token()
		if err != nil {
			return nil, err
		} else if token == json.Delim('}') {
			break
		}

		name := token.(string)
		if value, err := p.jsonValue(d, -1); err != nil {
			return nil, err
		} else {
			///////////////////////////////////////////////////
			// TODO: Load the class definition
			var dt DataType
			switch v := value.(type) {
			case string:
				if color, err := ParseColor(v); err == nil {
					dt = TypeColor
					value = color
				} else {
					dt = TypeString
				}
			case int:
				dt = TypeInt
			case Color:
				dt = TypeColor
			case float64:
				if float64(int(v)) == v {
					dt = TypeInt
					value = int(v)
				} else {
					dt = TypeFloat
				}
			case bool:
				dt = TypeBool
			case Properties:
				dt = TypeClass
			default:
				dt = -1
			}
			//////////////////////////////////////////////////
			props[name] = Property{Name: name, Value: value, Type: dt}
		}
	}

	return props, nil
}

// vim: ts=4
