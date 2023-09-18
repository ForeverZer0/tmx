package tmx

import (
	"fmt"
	"log"
)

// ErrInvalidEnum is an error used when an invalid enumeration value is given.
type ErrInvalidEnum struct {
	EnumType string
	Value    string
}

// Error implements the error interface.
func (e *ErrInvalidEnum) Error() string {
	return fmt.Sprintf("%s is not a valid  %s", e.Value, e.EnumType)
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
