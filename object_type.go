package tmx

import "fmt"

// ENUM(unspecified, left, right, top, bottom, justify, topleft, topright, bottomleft, bottomright, centerh, centerv, center)
type Align uint8

// nspecified, topleft, top, topright, left, center, right, bottomleft, bottom and bottomright. The default value is unspecified,

// ObjectType provides strongly-typed constants describing types of map objects.
type ObjectType int

const (
	// ObjectNone describes a standard object with no specialized type.
	ObjectNone ObjectType = iota
	// ObjectEllipse describes a elliptical shape, using the existing fields to determine the size
	// of the ellipse.
	ObjectEllipse
	// ObjectPoint describes a single location, using the existing fields to determine the location
	// of the point.
	ObjectPoint
	// ObjectPolygon describes a polygon shape. The Points field will be populated with points that
	// defines a closed shape.
	ObjectPolygon
	// ObjectPolyline describes a polyline shape. The Points field will be populated with points
	// that defines an open shape.
	ObjectPolyline
	// ObjectText describes a text object. The Text field will be populated with an object that
	// defines the font/text to be rendered.
	ObjectText
)

const _ObjectTypeName = "noneellipsepointpolygonpolylinetext"

var _ObjectTypeMap = map[ObjectType]string{
	ObjectNone:     _ObjectTypeName[0:4],
	ObjectEllipse:  _ObjectTypeName[4:11],
	ObjectPoint:    _ObjectTypeName[11:16],
	ObjectPolygon:  _ObjectTypeName[16:23],
	ObjectPolyline: _ObjectTypeName[23:31],
	ObjectText:     _ObjectTypeName[31:35],
}

// String implements the Stringer interface.
func (x ObjectType) String() string {
	if str, ok := _ObjectTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("ObjectType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ObjectType) IsValid() bool {
	_, ok := _ObjectTypeMap[x]
	return ok
}

var _ObjectKindValue = map[string]ObjectType{
	_ObjectTypeName[0:4]:   ObjectNone,
	_ObjectTypeName[4:11]:  ObjectEllipse,
	_ObjectTypeName[11:16]: ObjectPoint,
	_ObjectTypeName[16:23]: ObjectPolygon,
	_ObjectTypeName[23:31]: ObjectPolyline,
	_ObjectTypeName[31:35]: ObjectText,
}

// parseObjectType attempts to convert a string to a ObjectKind.
func parseObjectType(name string) (ObjectType, error) {
	if x, ok := _ObjectKindValue[name]; ok {
		return x, nil
	}
	return ObjectType(0), errInvalidEnum("ObjectType", name)
}

// MarshalText implements the text marshaller method.
func (x ObjectType) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *ObjectType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseObjectType(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

// vim: ts=4
