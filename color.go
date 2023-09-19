package tmx

import (
	"fmt"
	"strconv"
	"strings"
)

// Color describes a 32-bit RGBA color (0xAABBGGRR).
type Color uint32

// A is the value of the alpha channel.
func (c Color) A() uint8 {
	return uint8((c >> 24) & 0xFF)
}

// R is the value of the red channel.
func (c Color) B() uint8 {
	return uint8((c >> 16) & 0xFF)
}

// G is the value of the green channel.
func (c Color) G() uint8 {
	return uint8((c >> 8) & 0xFF)
}

// B is the value of the blue channel.
func (c Color) R() uint8 {
	return uint8(c & 0xFF)
}

// String returns the string representation of the color.
func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x%02x", c.A(), c.R(), c.G(), c.B())
}

// NewRGB creates a new fully opaque color from the specified values.
func NewRGB(r, g, b uint8) Color {
	return 0xFF000000 | (Color(b) << 16) | (Color(g) << 8) | Color(r)
}

// NewRGBA creates a new color from the specified values.
func NewRGBA(r, g, b, a uint8) Color {
	return (Color(a) << 24) | (Color(b) << 16) | (Color(g) << 8) | Color(r)
}

// ParseColor parses a string in the form of "#AARRGGBB" or "#RRGGBB" to a Color value.
func ParseColor(str string) (color Color, err error) {
	var result uint64
	if strings.HasPrefix(str, "#") {
		result, err = strconv.ParseUint(str[1:], 16, 32)
	} else {
		result, err = strconv.ParseUint(str, 16, 32)
	}

	if err == nil {
		if len(str) < 8 {
			result |= 0xFF000000
		}
		// we do a little bit-shifting to convert AARRGGBB to AABBGGRR
		// just need to swap the red and blue positions
		color = Color(result)
		color = (color & 0xFF00FF00) | ((color & 0xFF) << 16) | ((color >> 16) & 0xFF)
	}
	return
}

// Implements the encoding.TextMarshaler interface.
func (c Color) MarshalText() (text []byte, err error) {
	return []byte(c.String()), nil
}

// Implements the encoding.TextUnmarshaler interface.
func (c *Color) UnmarshalText(text []byte) error {
	if color, err := ParseColor(string(text)); err == nil {
		*c = color
		return nil
	} else {
		return err
	}
}

// vim: ts=4
