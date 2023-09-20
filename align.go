package tmx

import "fmt"

// Align is a set of bitflags that describe the relational alignment of an object.
type Align uint8

const (
	// AlignUnspecified indicates no alignment was specified and/or an invalid value.
	AlignUnspecified Align = 0x00
	AlignLeft        Align = 0x01
	AlignRight       Align = 0x02
	AlignTop         Align = 0x04
	AlignBottom      Align = 0x08
	AlignJustify     Align = 0x10
	AlignCenterH           = AlignLeft | AlignRight
	AlignCenterV           = AlignTop | AlignBottom
	AlignCenter            = AlignCenterH | AlignCenterV
	AlignTopLeft           = AlignTop | AlignLeft
	AlignTopRight          = AlignTop | AlignRight
	AlignBottomLeft        = AlignBottom | AlignLeft
	AlignBottomRight       = AlignBottom | AlignRight

	clearHorizontal = ^AlignCenterH
	clearVertical   = ^AlignCenterV
)

const _AlignName = "unspecifiedleftrighttopbottomjustifytoplefttoprightbottomleftbottomrightcenterhcentervcenter"

var _AlignMap = map[Align]string{
	AlignUnspecified: _AlignName[0:11],
	AlignLeft:        _AlignName[11:15],
	AlignRight:       _AlignName[15:20],
	AlignTop:         _AlignName[20:23],
	AlignBottom:      _AlignName[23:29],
	AlignJustify:     _AlignName[29:36],
	AlignTopLeft:     _AlignName[36:43],
	AlignTopRight:    _AlignName[43:51],
	AlignBottomLeft:  _AlignName[51:61],
	AlignBottomRight: _AlignName[61:72],
	AlignCenterH:     _AlignName[72:79],
	AlignCenterV:     _AlignName[79:86],
	AlignCenter:      _AlignName[86:92],
}

// String implements the Stringer interface.
func (x Align) String() string {
	if str, ok := _AlignMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Align(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Align) IsValid() bool {
	_, ok := _AlignMap[x]
	return ok
}

var _AlignValue = map[string]Align{
	_AlignName[0:11]:  AlignUnspecified,
	_AlignName[11:15]: AlignLeft,
	_AlignName[15:20]: AlignRight,
	_AlignName[20:23]: AlignTop,
	_AlignName[23:29]: AlignBottom,
	_AlignName[29:36]: AlignJustify,
	_AlignName[36:43]: AlignTopLeft,
	_AlignName[43:51]: AlignTopRight,
	_AlignName[51:61]: AlignBottomLeft,
	_AlignName[61:72]: AlignBottomRight,
	_AlignName[72:79]: AlignCenterH,
	_AlignName[79:86]: AlignCenterV,
	_AlignName[86:92]: AlignCenter,
}

// parseAlign attempts to convert a string to a Align.
func parseAlign(name string) (Align, error) {
	if x, ok := _AlignValue[name]; ok {
		return x, nil
	}
	return Align(0), errInvalidEnum("Align", name)
}

// MarshalText implements the text marshaller method.
func (x Align) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Align) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseAlign(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
