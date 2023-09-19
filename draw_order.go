package tmx

import "fmt"

type DrawOrder int

const (
	// DrawTopdown is a DrawOrder of type Topdown.
	DrawTopdown DrawOrder = iota
	// DrawIndex is a DrawOrder of type Index.
	DrawIndex
)

const _DrawOrderName = "topdownindex"

var _DrawOrderMap = map[DrawOrder]string{
	DrawTopdown: _DrawOrderName[0:7],
	DrawIndex:   _DrawOrderName[7:12],
}

// String implements the Stringer interface.
func (x DrawOrder) String() string {
	if str, ok := _DrawOrderMap[x]; ok {
		return str
	}
	return fmt.Sprintf("DrawOrder(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x DrawOrder) IsValid() bool {
	_, ok := _DrawOrderMap[x]
	return ok
}

var _DrawOrderValue = map[string]DrawOrder{
	_DrawOrderName[0:7]:  DrawTopdown,
	_DrawOrderName[7:12]: DrawIndex,
}

// parseDrawOrder attempts to convert a string to a DrawOrder.
func parseDrawOrder(name string) (DrawOrder, error) {
	if x, ok := _DrawOrderValue[name]; ok {
		return x, nil
	}
	return DrawOrder(0), errInvalidEnum("DrawOrder", name)
}

// MarshalText implements the text marshaller method.
func (x DrawOrder) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *DrawOrder) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseDrawOrder(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
