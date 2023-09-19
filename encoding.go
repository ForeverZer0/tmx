package tmx

import "fmt"

// Encoding provides strongly-typed constants describing data encoding used in the TMX format.
type Encoding int

const (
	// EncodingNone indicates unencoded text.
	EncodingNone Encoding = iota
	// EncodingCSV indicates a comma-separated list of values.
	EncodingCSV
	// EncodingBase64 indicates base64-encoded text.
	EncodingBase64
)

const _EncodingName = "nonecsvbase64"

var _EncodingMap = map[Encoding]string{
	EncodingNone:   _EncodingName[0:4],
	EncodingCSV:    _EncodingName[4:7],
	EncodingBase64: _EncodingName[7:13],
}

// String implements the Stringer interface.
func (x Encoding) String() string {
	if str, ok := _EncodingMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Encoding(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Encoding) IsValid() bool {
	_, ok := _EncodingMap[x]
	return ok
}

var _EncodingValue = map[string]Encoding{
	_EncodingName[0:4]:  EncodingNone,
	_EncodingName[4:7]:  EncodingCSV,
	_EncodingName[7:13]: EncodingBase64,
}

// parseEncoding attempts to convert a string to a Encoding.
func parseEncoding(name string) (Encoding, error) {
	if x, ok := _EncodingValue[name]; ok {
		return x, nil
	}
	return Encoding(0), errInvalidEnum("Encoding", name)
}

// MarshalText implements the text marshaller method.
func (x Encoding) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Encoding) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseEncoding(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
