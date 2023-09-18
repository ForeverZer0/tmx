package tmx

import (
	"encoding/json"
	"encoding/xml"
	"strings"
)

// Properties describes a map with string keys and Property values. Order of values is not
// guaranteed.
type Properties map[string]Property

// UnmarshalXML implements the xml.Unmarshaler interface.
func (p *Properties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		}
		if child, ok := token.(xml.StartElement); ok {
			if child.Name.Local != "property" {
				logElem(child.Name.Local, start.Name.Local)
			} else {
				var prop Property
				if err = prop.UnmarshalXML(d, child); err != nil {
					return err
				}
				(*p)[prop.Name] = prop
			}
		}

		token, err = d.Token()
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (p *Properties) UnmarshalJSON(data []byte) error {
	var props []Property
	if err := json.Unmarshal(data, &props); err != nil {
		return err
	}

	for _, prop := range props {
		(*p)[prop.Name] = prop
	}

	return nil
}

// String implements the Stringer interface
func (p Properties) String() string {
	var sb strings.Builder
	sb.WriteRune('{')

	var n int
	for _, v := range p {
		if n > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(v.String())
		n++
	}

	sb.WriteRune('}')
	return sb.String()
}

func propValue[T any](p Properties, name string) (value T, ok bool) {
	if p == nil {
		return
	}
	if prop, found := p[name]; found {
		value, ok = prop.Value.(T)
	}
	return
}

func mustValue[T any](p Properties, name string, def T) T {
	if value, ok := propValue[T](p, name); ok {
		return value
	}
	return def
}

// GetBool retrieves a boolean property with the given name, including a flag if the property
// was found and returned successfully.
func (p Properties) GetBool(name string) (bool, bool) {
	return propValue[bool](p, name)
}

// GetInt retrieves an integer property with the given name, including a flag if the property
// was found and returned successfully.
//
// If a property with the specified name is found, but contains a float value, it will be
// automatically converted to an integer and return successfully.
func (p Properties) GetInt(name string) (int, bool) {
	if n, ok := propValue[int](p, name); ok {
		return n, true
	}
	if n, ok := propValue[float64](p, name); ok {
		return int(n), true
	}
	return 0, false
}

// GetFloat retrieves a float property with the given name, including a flag if the property
// was found and returned successfully.
//
// If a property with the specified name is found, but contains an integer value, it will be
// automatically converted to a float and return successfully.
func (p Properties) GetFloat(name string) (float64, bool) {
	if n, ok := propValue[float64](p, name); ok {
		return n, true
	}
	if n, ok := propValue[int](p, name); ok {
		return float64(n), true
	}
	return 0, false
}

// GetColor retrieves a color property with the given name, including a flag if the property
// was found and returned successfully.
func (p Properties) GetColor(name string) (Color, bool) {
	return propValue[Color](p, name)
}

// GetClass retrieves a custom class property with the given name, including a flag if the
// property was found and returned successfully.
func (p Properties) GetClass(name string) (Properties, bool) {
	return propValue[Properties](p, name)
}

// MustBool retrieves a boolean property with the given name, or the given default
// value upon failure.
func (p Properties) MustBool(name string, def bool) bool {
	return mustValue(p, name, def)
}

// MustInt retrieves an integer property with the given name, or the given default
// value upon failure.
//
// If a property with the specified name is found, but contains a float value, it will be
// automatically converted to an integer and return successfully.
func (p Properties) MustInt(name string, def int) int {
	return mustValue(p, name, def)
}

// MustFloat retrieves a float property with the given name, or the given default
// value upon failure.
//
// If a property with the specified name is found, but contains an integer value, it will be
// automatically converted to a float.
func (p Properties) MustFloat(name string, def float64) float64 {
	return mustValue(p, name, def)
}

// MustColor retrieves a color property with the given name, or the given default
// value upon failure.
func (p Properties) MustColor(name string, def Color) Color {
	return mustValue(p, name, def)
}

// Clone implements the Cloner interface.
func (p Properties) Clone() Properties {
	dup := make(Properties, len(p))
	for k, v := range p {
		dup[k] = v.Clone()
	}
	return dup
}

// vim: ts=4
