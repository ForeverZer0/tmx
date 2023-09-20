package tmx

import "fmt"


type FillMode int

const (
	// FillStretch is a FillMode of type Stretch.
	FillStretch FillMode = iota
	// FillPreserveAspect is a FillMode of type Preserve-Aspect-Fit.
	FillPreserveAspect
)

const _FillModeName = "stretchpreserve-aspect-fit"

var _FillModeMap = map[FillMode]string{
	FillStretch:           _FillModeName[0:7],
	FillPreserveAspect: _FillModeName[7:26],
}

// String implements the Stringer interface.
func (x FillMode) String() string {
	if str, ok := _FillModeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("FillMode(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x FillMode) IsValid() bool {
	_, ok := _FillModeMap[x]
	return ok
}

var _FillModeValue = map[string]FillMode{
	_FillModeName[0:7]:  FillStretch,
	_FillModeName[7:26]: FillPreserveAspect,
}

// parseFillMode attempts to convert a string to a FillMode.
func parseFillMode(name string) (FillMode, error) {
	if x, ok := _FillModeValue[name]; ok {
		return x, nil
	}
	return FillMode(0), errInvalidEnum("FillMode", name)
}

// MarshalText implements the text marshaller method.
func (x FillMode) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *FillMode) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseFillMode(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
