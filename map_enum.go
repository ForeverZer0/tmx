package tmx

import "fmt"

type (
	RenderOrder  int
	StaggerAxis  int
	StaggerIndex int
)

const (
	// RenderRightDown is a RenderOrder of type Right-Down.
	RenderRightDown RenderOrder = iota
	// RenderRightUp is a RenderOrder of type Right-Up.
	RenderRightUp
	// RenderLeftDown is a RenderOrder of type Left-Down.
	RenderLeftDown
	// RenderLeftUp is a RenderOrder of type Left-Up.
	RenderLeftUp
)

const _RenderOrderName = "right-downright-upleft-downleft-up"

var _RenderOrderMap = map[RenderOrder]string{
	RenderRightDown: _RenderOrderName[0:10],
	RenderRightUp:   _RenderOrderName[10:18],
	RenderLeftDown:  _RenderOrderName[18:27],
	RenderLeftUp:    _RenderOrderName[27:34],
}

// String implements the Stringer interface.
func (x RenderOrder) String() string {
	if str, ok := _RenderOrderMap[x]; ok {
		return str
	}
	return fmt.Sprintf("RenderOrder(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x RenderOrder) IsValid() bool {
	_, ok := _RenderOrderMap[x]
	return ok
}

var _RenderOrderValue = map[string]RenderOrder{
	_RenderOrderName[0:10]:  RenderRightDown,
	_RenderOrderName[10:18]: RenderRightUp,
	_RenderOrderName[18:27]: RenderLeftDown,
	_RenderOrderName[27:34]: RenderLeftUp,
}

// parseRenderOrder attempts to convert a string to a RenderOrder.
func parseRenderOrder(name string) (RenderOrder, error) {
	if x, ok := _RenderOrderValue[name]; ok {
		return x, nil
	}
	return RenderOrder(0), errInvalidEnum("RenderOrder", name)
}

// MarshalText implements the text marshaller method.
func (x RenderOrder) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *RenderOrder) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseRenderOrder(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

const (
	// StaggerX is a StaggerAxis of type X.
	StaggerX StaggerAxis = iota
	// StaggerY is a StaggerAxis of type Y.
	StaggerY
)

const _StaggerAxisName = "xy"

var _StaggerAxisMap = map[StaggerAxis]string{
	StaggerX: _StaggerAxisName[0:1],
	StaggerY: _StaggerAxisName[1:2],
}

// String implements the Stringer interface.
func (x StaggerAxis) String() string {
	if str, ok := _StaggerAxisMap[x]; ok {
		return str
	}
	return fmt.Sprintf("StaggerAxis(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x StaggerAxis) IsValid() bool {
	_, ok := _StaggerAxisMap[x]
	return ok
}

var _StaggerAxisValue = map[string]StaggerAxis{
	_StaggerAxisName[0:1]: StaggerX,
	_StaggerAxisName[1:2]: StaggerY,
}

// parseStaggerAxis attempts to convert a string to a StaggerAxis.
func parseStaggerAxis(name string) (StaggerAxis, error) {
	if x, ok := _StaggerAxisValue[name]; ok {
		return x, nil
	}
	return StaggerAxis(0), errInvalidEnum("StaggerAxis", name)
}

// MarshalText implements the text marshaller method.
func (x StaggerAxis) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *StaggerAxis) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseStaggerAxis(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

const (
	// StaggerEven is a StaggerIndex of type Even.
	StaggerEven StaggerIndex = iota
	// StaggerOdd is a StaggerIndex of type Odd.
	StaggerOdd
)

const _StaggerIndexName = "evenodd"

var _StaggerIndexMap = map[StaggerIndex]string{
	StaggerEven: _StaggerIndexName[0:4],
	StaggerOdd:  _StaggerIndexName[4:7],
}

// String implements the Stringer interface.
func (x StaggerIndex) String() string {
	if str, ok := _StaggerIndexMap[x]; ok {
		return str
	}
	return fmt.Sprintf("StaggerIndex(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x StaggerIndex) IsValid() bool {
	_, ok := _StaggerIndexMap[x]
	return ok
}

var _StaggerIndexValue = map[string]StaggerIndex{
	_StaggerIndexName[0:4]: StaggerEven,
	_StaggerIndexName[4:7]: StaggerOdd,
}

// parseStaggerIndex attempts to convert a string to a StaggerIndex.
func parseStaggerIndex(name string) (StaggerIndex, error) {
	if x, ok := _StaggerIndexValue[name]; ok {
		return x, nil
	}
	return StaggerIndex(0), errInvalidEnum("StaggerIndex", name)
}

// MarshalText implements the text marshaller method.
func (x StaggerIndex) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *StaggerIndex) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseStaggerIndex(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
