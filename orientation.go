package tmx

import "fmt"

// TODO: Documentation

type Orientation int

const (
	// Orthogonal is a Orientation of type Orthogonal.
	Orthogonal Orientation = iota
	// Isometric is a Orientation of type Isometric.
	Isometric
	// Staggered is a Orientation of type Staggered.
	Staggered
	// Hexagonal is a Orientation of type Hexagonal.
	Hexagonal
)

const _OrientationName = "orthogonalisometricstaggeredhexagonal"

var _OrientationMap = map[Orientation]string{
	Orthogonal: _OrientationName[0:10],
	Isometric:  _OrientationName[10:19],
	Staggered:  _OrientationName[19:28],
	Hexagonal:  _OrientationName[28:37],
}

// String implements the Stringer interface.
func (e Orientation) String() string {
	if str, ok := _OrientationMap[e]; ok {
		return str
	}
	return fmt.Sprintf("Orientation(%d)", e)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (e Orientation) IsValid() bool {
	_, ok := _OrientationMap[e]
	return ok
}

var _OrientationValue = map[string]Orientation{
	_OrientationName[0:10]:  Orthogonal,
	_OrientationName[10:19]: Isometric,
	_OrientationName[19:28]: Staggered,
	_OrientationName[28:37]: Hexagonal,
}

// parseOrientation attempts to convert a string to a Orientation.
func parseOrientation(name string) (Orientation, error) {
	if x, ok := _OrientationValue[name]; ok {
		return x, nil
	}
	return Orientation(0), errInvalidEnum("Orientation", name)
}

// MarshalText implements the text marshaller method.
func (e Orientation) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (e *Orientation) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseOrientation(name)
	if err != nil {
		return err
	}
	*e = tmp
	return nil
}
