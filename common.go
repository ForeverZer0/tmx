package tmx

import (
	"fmt"
	"log"
)

type ErrFormat struct{
	Message string
}

func (e *ErrFormat) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("format error: %s", e.Message)
	}
	return "forrmat error"
}

func errFormat(message string) error {
	return &ErrFormat{Message: message}
}

var ErrExpectedObject error = errFormat("expected JSON object")
var ErrExpectedArray error = errFormat("expected JSON array")

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

// vim: ts=4
