package tmx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// CustomTypes maintains references to all known user-defined types.
//
// As the [JSON format] does not include this information in its output, the
// class definition is required to have been added/loaded prior to implementing
// properties being parsed. They can still be used, but will be missing type
// information, as there is no way to determine some types apart (i.e. a "0" could
// be an integer, a float, or an object ID, or a string could be a string, a file,
// or a color).
//
// XML does not suffer this this same problem, as it includes the property type
// within the normal property definition.
//
// [JSON format]: https://github.com/mapeditor/tiled/issues/3820
var CustomTypes map[string]*CustomClass

// CustomClass is a type used to define custom property types.
type CustomClass struct {
	// Name is the user-defined name of the type.
	Name string
	// Members contain a collection properties that described the name, type, and
	// default value for the members that make up the class.
	Members Properties
}

func LoadTypes(path string) error {
	if CustomTypes == nil {
		CustomTypes = make(map[string]*CustomClass)
	}

	abs, err := FindPath(path)
	if err != nil {
		return err
	} 

	file, err := os.Open(abs);
	if err != nil {
		return err
	} 
	defer file.Close()

	format := detectFileExt(abs)
	switch format {
	case FormatXML:
		type x struct {
			Types []*CustomClass `xml:"objecttype"`
		}
		var result x
		d := xml.NewDecoder(file)
		if err := d.Decode(&result); err != nil {
			return err
		}
		for _, class := range result.Types {
			CustomTypes[class.Name] = class
		}
	case FormatJSON:
		d := json.NewDecoder(file)
		if token, err := d.Token(); err != nil {
			return err
		} else if token != json.Delim('[') {
			return ErrExpectedArray
		}	

		for d.More() {
			var prop CustomClass
			if err = d.Decode(&prop); err != nil {
				return err
			}
			if len(prop.Members) > 0 {
				CustomTypes[prop.Name] = &prop
			}
		}
	default:
		return errInvalidEnum("Format", fmt.Sprintf("Format(%d)", format))
	}

	return nil
}

// NewClass initializes and returns a new custom class with the specified name.
//
// Custom classes define the structure of user-defined Property types, providing
// their data type, default value, etc.
func NewClass(name string) *CustomClass {
	class := &CustomClass{
		Name:    name,
		Members: make(Properties),
	}

	if CustomTypes == nil {
		CustomTypes = make(map[string]*CustomClass)
	}
	CustomTypes[name] = class
	return class
}

// String implements the Stringer interface.
func (c *CustomClass) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s:{", strconv.Quote(c.Name))
	var count int
	for _, value := range c.Members {
		if count > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(value.String())
		count++
	}
	sb.WriteRune('}')
	return sb.String()
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (c *CustomClass) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if c.Members == nil {
		c.Members = make(Properties)
	}

	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "name":
			c.Name = attr.Value
		case "color":
			// Skip
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
			case "property":
				var prop Property
				if err = d.DecodeElement(&prop, &child); err != nil {
					return err
				}
				c.Members[prop.Name] = prop
			default:
				logElem(child.Name.Local, start.Name.Local)
			}
		}
		token, err = d.Token()
	}

	if CustomTypes == nil {
		CustomTypes = make(map[string]*CustomClass)
	}
	CustomTypes[c.Name] = c
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CustomClass) UnmarshalJSON(data []byte) error {
	if c.Members == nil {
		c.Members = make(Properties)
	}

	d := json.NewDecoder(bytes.NewBuffer(data))
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
			if c.Name, err = jsonProp[string](d); err != nil {
				return err
			}
		case "members":
			if token, err = d.Token(); err != nil {
				return err
			} else if token != json.Delim('[') {
				return ErrExpectedArray
			}
			for d.More() {
				var prop Property
				if err = d.Decode(&prop); err != nil {
					return err
				}
				c.Members[prop.Name] = prop
			}
			// Consume the closing ']' token
			if token, err = d.Token(); err != nil {
				return err
			}
		default:
			// Do not log unhandled values
			jsonSkip(d)
		}
	}

	if CustomTypes == nil {
		CustomTypes = make(map[string]*CustomClass)
	}
	CustomTypes[c.Name] = c	
	return nil
}

// vim: ts=4
