package tmx

import "fmt"

// DataType describes the value type of a Property.
type DataType int

const (
	// TypeString indicates the value is a string type.
	TypeString DataType = iota
	// TypeInt indicates the value is an int type.
	TypeInt
	// TypeFloat indicates the value is a float64 type.
	TypeFloat
	// TypeBool indicates the value is a bool type.
	TypeBool
	// TypeColor indicates the value is a Color type.
	TypeColor
	// TypeFile indicates the value is a string type that
	// represents a filesystem path.
	TypeFile
	// TypeObject indicates the value is an int type that
	// is the ID of an Object on the map.
	TypeObject
	// TypeClass indicates the value is a set of Properties
	// describing a user-defined type.
	TypeClass
)

const _DataTypeName = "stringintfloatboolcolorfileobjectclass"

var _DataTypeMap = map[DataType]string{
	TypeString: _DataTypeName[0:6],
	TypeInt:    _DataTypeName[6:9],
	TypeFloat:  _DataTypeName[9:14],
	TypeBool:   _DataTypeName[14:18],
	TypeColor:  _DataTypeName[18:23],
	TypeFile:   _DataTypeName[23:27],
	TypeObject: _DataTypeName[27:33],
	TypeClass:  _DataTypeName[33:38],
}

// String implements the Stringer interface.
func (x DataType) String() string {
	if str, ok := _DataTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("DataType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x DataType) IsValid() bool {
	_, ok := _DataTypeMap[x]
	return ok
}

var _DataTypeValue = map[string]DataType{
	_DataTypeName[0:6]:   TypeString,
	_DataTypeName[6:9]:   TypeInt,
	_DataTypeName[9:14]:  TypeFloat,
	_DataTypeName[14:18]: TypeBool,
	_DataTypeName[18:23]: TypeColor,
	_DataTypeName[23:27]: TypeFile,
	_DataTypeName[27:33]: TypeObject,
	_DataTypeName[33:38]: TypeClass,
}

// parseDataType attempts to convert a string to a DataType.
func parseDataType(name string) (DataType, error) {
	if x, ok := _DataTypeValue[name]; ok {
		return x, nil
	}
	return DataType(0), errInvalidEnum("DataType", name)
}

// MarshalText implements the text marshaller method.
func (x DataType) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *DataType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseDataType(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
