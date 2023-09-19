package tmx

import "fmt"

// Compression provides strongly-typed constants describing compression methods.
type Compression int

const (
	// CompressionNone indicates no compression.
	CompressionNone Compression = iota
	// CompressionGzip indicates Gzip compression.
	CompressionGzip
	// CompressionGzip indicates Zlib compression.
	CompressionZlib
	// CompressionGzip indicates Z-Standard compression.
	CompressionZstd
)

const _CompressionName = "nonegzipzlibzstd"

var _CompressionMap = map[Compression]string{
	CompressionNone: _CompressionName[0:4],
	CompressionGzip: _CompressionName[4:8],
	CompressionZlib: _CompressionName[8:12],
	CompressionZstd: _CompressionName[12:16],
}

// String implements the Stringer interface.
func (x Compression) String() string {
	if str, ok := _CompressionMap[x]; ok {
		return str
	}
	return fmt.Sprintf("Compression(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x Compression) IsValid() bool {
	_, ok := _CompressionMap[x]
	return ok
}

var _CompressionValue = map[string]Compression{
	_CompressionName[0:4]:   CompressionNone,
	_CompressionName[4:8]:   CompressionGzip,
	_CompressionName[8:12]:  CompressionZlib,
	_CompressionName[12:16]: CompressionZstd,
}

// parseCompression attempts to convert a string to a Compression.
func parseCompression(name string) (Compression, error) {
	if x, ok := _CompressionValue[name]; ok {
		return x, nil
	}
	return Compression(0), errInvalidEnum("Compression", name)
}

// MarshalText implements the text marshaller method.
func (x Compression) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Compression) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseCompression(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// vim: ts=4
