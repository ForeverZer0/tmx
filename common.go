package tmx

import (
	"encoding/json"
	"fmt"
	"log"
)

// ErrFormat is an error type used for format-related errors.
type ErrFormat struct {
	// Message describes the cause of the error.
	Message string
}

// Cloner represents a type that can create a deep-clone.
type Cloner[T any] interface {
	// Clone creates a deep-clone and returns it.
	Clone() T
}

// ErrInvalidEnum is an error used when an invalid enumeration value is given.
type ErrInvalidEnum struct {
	EnumType string
	Value    string
}

var (
	ErrExpectedObject error = errFormat("expected JSON object")
	ErrExpectedArray  error = errFormat("expected JSON array")
)

// Error implements the error interface.
func (e *ErrFormat) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("format error: %s", e.Message)
	}
	return "forrmat error"
}

func errFormat(format string, args ...any) error {
	return &ErrFormat{Message: fmt.Sprintf(format, args...)}
}

// Error implements the error interface.
func (e *ErrInvalidEnum) Error() string {
	return fmt.Sprintf("%s is not a valid %s", e.Value, e.EnumType)
}

// errInvalidEnum is a helper function to create a new ErrInvalidEnum error.
func errInvalidEnum(enum, value string) error {
	return &ErrInvalidEnum{EnumType: enum, Value: value}
}

// logElem is used to log an unhandled/unrecognized element in TMX document.
func logElem(name, parent string) {
	log.Printf(`skipped unrecognized child element in <%s> in <%s>`, name, parent)
}

// logAttr is used to log an unhandled/unrecognized attribute in TMX document.
func logAttr(name, parent string) {
	log.Printf(`skipped unrecognized child attribute in "%s" in <%s>`, name, parent)
}

// logProp is used to log an unhandled/unrecognized property in TMJ document.
func logProp(name, parent string) {
	log.Printf(`skipped unrecognized child property in "%s" in "%s"`, name, parent)
}

var terrainWarn bool

func logTerrain() {
	if !terrainWarn {
		log.Println("use of terrains is deprecated, and has been replaced by wangsets")
		terrainWarn = true
	}
}

// jsonProp reads JSON value of the given type.
func jsonProp[T any](d *json.Decoder) (value T, err error) {
	var token json.Token
	token, err = d.Token()
	if err != nil {
		return
	}

	var ok bool
	if value, ok = token.(T); !ok {
		err = errFormat("expected type of %T", value)
	}
	return
}

// jsonSkip consumes the current JSON value without processing.
func jsonSkip(d *json.Decoder) error {
	var d1, d2 int

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch token {
		case json.Delim('}'):
			d1--
		case json.Delim(']'):
			d2--
		case json.Delim('{'):
			d1++
		case json.Delim('['):
			d2++
		}
		if d1 == 0 && d2 == 0 {
			break
		}
	}

	return nil
}

// vim: ts=4
