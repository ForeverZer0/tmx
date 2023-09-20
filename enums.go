package tmx

import "fmt"

type (
	// Align is a set of bitflags that describe the relational alignment of an object.
	Align uint8
	// Compression provides strongly-typed constants describing compression methods.
	Compression int
	// Encoding provides strongly-typed constants describing data encoding used in the TMX format.
	Encoding int

	FillMode int
	// LayerType provides strongly-typed constants describing the type of a layer.
	LayerType byte
	// ObjectType provides strongly-typed constants describing types of map objects.
	ObjectType int
	// Orientation describes methods the method/perspective a map is rendered.
	Orientation int
	TileRender  int
	// WangType describes the behavior of terrain generation.
	WangType int
	// DataType describes the value type of a Property.
	DataType    int
	DrawOrder   int
	RenderOrder int

	StaggerAxis  int
	StaggerIndex int
	// Format describes the format of a TMX document.
	Format int
	// setFlags are used internally for objects to track which fields were explicity
	// to determine how a template is inherited.
	setFlags uint32
	// FontStyle provides strongly-typed constants for describing styles of a rendered
	// font, implemented as a set of bit flags.
	FontStyle uint8
)

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

const (
	// EncodingNone indicates unencoded text.
	EncodingNone Encoding = iota
	// EncodingCSV indicates a comma-separated list of values.
	EncodingCSV
	// EncodingBase64 indicates base64-encoded text.
	EncodingBase64
)

const (
	// FillStretch is a FillMode of type Stretch.
	FillStretch FillMode = iota
	// FillPreserveAspect is a FillMode of type Preserve-Aspect-Fit.
	FillPreserveAspect
)

const (
	// LayerNone is a LayerType of type None.
	LayerNone LayerType = 0x00
	// LayerTile is a LayerType of type Tile.
	LayerTile LayerType = 0x01
	// LayerImage is a LayerType of type Image.
	LayerImage LayerType = 0x02
	// LayerObject is a LayerType of type Object.
	LayerObject LayerType = 0x04
	// LayerGroup is a LayerType of type Group.
	LayerGroup LayerType = 0x08
	// LayerAll is a LayerType of type All.
	LayerAll LayerType = 0xFF
)

const (
	// ObjectNone describes a standard object with no specialized type (e.g. a rectangle).
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

const (
	// Orthogonal is a Orientation of type Orthogonal.
	Orthogonal Orientation = iota
	// Isometric is a Orientation of type Isometric.
	Isometric
	// Staggered is a Orientation of type Staggered.
	Staggered
	// Hexagonal is a Orientation of type Hexagonal.
	Hexagonal
)

const (
	// RenderTile is a TileRender of type Tile.
	RenderTile TileRender = iota
	// RenderGrid is a TileRender of type Grid.
	RenderGrid
)

const (
	// WangTypeCorner is a WangType of type Corner.
	WangTypeCorner WangType = iota
	// WangTypeEdge is a WangType of type Edge.
	WangTypeEdge
	// WangTypeMixed is a WangType of type Mixed.
	WangTypeMixed
)

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

const (
	// DrawTopDown is a DrawOrder of type Topdown.
	DrawTopDown DrawOrder = iota
	// DrawIndex is a DrawOrder of type Index.
	DrawIndex
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

const (
	// StaggerX is a StaggerAxis of type X.
	StaggerX StaggerAxis = iota
	// StaggerY is a StaggerAxis of type Y.
	StaggerY
)

const (
	// StaggerEven is a StaggerIndex of type Even.
	StaggerEven StaggerIndex = iota
	// StaggerOdd is a StaggerIndex of type Odd.
	StaggerOdd
)

const (
	// FormatUnknown indicates an unknown/undefined TMX format.
	FormatUnknown Format = iota
	// FormatXML indicates the standard XML-based TMX format.
	FormatXML
	// FormatJSON indicates the standard JSON-based TMX format.
	FormatJSON
)

const (
	// StyleBold indicates bold font style.
	StyleBold FontStyle = 1 << iota
	// StyleItalic indicates italic font style.
	StyleItalic
	// StyleUnderline indicates underline font style.
	StyleUnderline
	// StyleStrikeout indicates strikeout font style.
	StyleStrikeout
	// StyleKerning indicates if kerning should be used when rendering the text.
	StyleKerning
)

const (
	flagName setFlags = 1 << iota
	flagClass
	flagX
	flagY
	flagWidth
	flagHeight
	flagRotation
	flagGID
	flagVisible
	flagKind
	flagPoints
	// Text flags
	flagFont
	flagFontSize
	flagTextWrap
	flagTextColor
	flagBold
	flagItalic
	flagUnderline
	flagStrikeout
	flagKerning
	flagHAlign
	flagVAlign
	flagText
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

const _FillModeName = "stretchpreserve-aspect-fit"

var _FillModeMap = map[FillMode]string{
	FillStretch:        _FillModeName[0:7],
	FillPreserveAspect: _FillModeName[7:26],
}

// String implements the Stringer interface.
func (x FillMode) String() string {
	if str, ok := _FillModeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("FillMode(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x FillMode) IsValid() bool {
	_, ok := _FillModeMap[x]
	return ok
}

var _FillModeValue = map[string]FillMode{
	_FillModeName[0:7]:  FillStretch,
	_FillModeName[7:26]: FillPreserveAspect,
}

// parseFillMode attempts to convert a string to a FillMode.
func parseFillMode(name string) (FillMode, error) {
	if x, ok := _FillModeValue[name]; ok {
		return x, nil
	}
	return FillMode(0), errInvalidEnum("FillMode", name)
}

// MarshalText implements the text marshaller method.
func (x FillMode) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *FillMode) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseFillMode(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

const _LayerTypeName = "nonetilelayerobjectgroupimagelayergroupall"

var _LayerTypeMap = map[LayerType]string{
	LayerNone:   _LayerTypeName[0:4],
	LayerTile:   _LayerTypeName[4:13],
	LayerObject: _LayerTypeName[13:24],
	LayerImage:  _LayerTypeName[24:34],
	LayerGroup:  _LayerTypeName[34:39],
	LayerAll:    _LayerTypeName[39:42],
}

// String implements the Stringer interface.
func (x LayerType) String() string {
	if str, ok := _LayerTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("LayerType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x LayerType) IsValid() bool {
	_, ok := _LayerTypeMap[x]
	return ok
}

var _LayerTypeValue = map[string]LayerType{
	_LayerTypeName[0:4]:   LayerNone,
	_LayerTypeName[4:13]:  LayerTile,
	_LayerTypeName[13:24]: LayerObject,
	_LayerTypeName[24:34]: LayerImage,
	_LayerTypeName[34:39]: LayerGroup,
	_LayerTypeName[39:42]: LayerAll,
}

// parseLayerType attempts to convert a string to a LayerType.
func parseLayerType(name string) (LayerType, error) {
	if x, ok := _LayerTypeValue[name]; ok {
		return x, nil
	}
	return LayerType(0), errInvalidEnum("LayerType", name)
}

// MarshalText implements the text marshaller method.
func (x LayerType) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *LayerType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseLayerType(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

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

const _OrientationName = "orthogonalisometricstaggeredhexagonal"

var _OrientationMap = map[Orientation]string{
	Orthogonal: _OrientationName[0:10],
	Isometric:  _OrientationName[10:19],
	Staggered:  _OrientationName[19:28],
	Hexagonal:  _OrientationName[28:37],
}

// String implements the Stringer interface.
func (e Orientation) String() string {
	if str, ok := _OrientationMap[e]; ok {
		return str
	}
	return fmt.Sprintf("Orientation(%d)", e)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (e Orientation) IsValid() bool {
	_, ok := _OrientationMap[e]
	return ok
}

var _OrientationValue = map[string]Orientation{
	_OrientationName[0:10]:  Orthogonal,
	_OrientationName[10:19]: Isometric,
	_OrientationName[19:28]: Staggered,
	_OrientationName[28:37]: Hexagonal,
}

// parseOrientation attempts to convert a string to a Orientation.
func parseOrientation(name string) (Orientation, error) {
	if x, ok := _OrientationValue[name]; ok {
		return x, nil
	}
	return Orientation(0), errInvalidEnum("Orientation", name)
}

// MarshalText implements the text marshaller method.
func (e Orientation) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (e *Orientation) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseOrientation(name)
	if err != nil {
		return err
	}
	*e = tmp
	return nil
}

const _TileRenderName = "tilegrid"

var _TileRenderMap = map[TileRender]string{
	RenderTile: _TileRenderName[0:4],
	RenderGrid: _TileRenderName[4:8],
}

// String implements the Stringer interface.
func (x TileRender) String() string {
	if str, ok := _TileRenderMap[x]; ok {
		return str
	}
	return fmt.Sprintf("TileRender(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x TileRender) IsValid() bool {
	_, ok := _TileRenderMap[x]
	return ok
}

var _TileRenderValue = map[string]TileRender{
	_TileRenderName[0:4]: RenderTile,
	_TileRenderName[4:8]: RenderGrid,
}

// parseTileRender attempts to convert a string to a TileRender.
func parseTileRender(name string) (TileRender, error) {
	if x, ok := _TileRenderValue[name]; ok {
		return x, nil
	}
	return TileRender(0), errInvalidEnum("TileRender", name)
}

// MarshalText implements the text marshaller method.
func (x TileRender) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *TileRender) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseTileRender(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

const _WangTypeName = "corneredgemixed"

var _WangTypeMap = map[WangType]string{
	WangTypeCorner: _WangTypeName[0:6],
	WangTypeEdge:   _WangTypeName[6:10],
	WangTypeMixed:  _WangTypeName[10:15],
}

// String implements the Stringer interface.
func (e WangType) String() string {
	if str, ok := _WangTypeMap[e]; ok {
		return str
	}
	return fmt.Sprintf("WangType(%d)", e)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (e WangType) IsValid() bool {
	_, ok := _WangTypeMap[e]
	return ok
}

var _WangTypeValue = map[string]WangType{
	_WangTypeName[0:6]:   WangTypeCorner,
	_WangTypeName[6:10]:  WangTypeEdge,
	_WangTypeName[10:15]: WangTypeMixed,
}

// parseWangType attempts to convert a string to a WangType.
func parseWangType(name string) (WangType, error) {
	if x, ok := _WangTypeValue[name]; ok {
		return x, nil
	}
	return WangType(0), errInvalidEnum("WangType", name)
}

// MarshalText implements the text marshaller method.
func (e WangType) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (e *WangType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseWangType(name)
	if err != nil {
		return err
	}
	*e = tmp
	return nil
}

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

const _DrawOrderName = "topdownindex"

var _DrawOrderMap = map[DrawOrder]string{
	DrawTopDown: _DrawOrderName[0:7],
	DrawIndex:   _DrawOrderName[7:12],
}

// String implements the Stringer interface.
func (x DrawOrder) String() string {
	if str, ok := _DrawOrderMap[x]; ok {
		return str
	}
	return fmt.Sprintf("DrawOrder(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x DrawOrder) IsValid() bool {
	_, ok := _DrawOrderMap[x]
	return ok
}

var _DrawOrderValue = map[string]DrawOrder{
	_DrawOrderName[0:7]:  DrawTopDown,
	_DrawOrderName[7:12]: DrawIndex,
}

// parseDrawOrder attempts to convert a string to a DrawOrder.
func parseDrawOrder(name string) (DrawOrder, error) {
	if x, ok := _DrawOrderValue[name]; ok {
		return x, nil
	}
	return DrawOrder(0), errInvalidEnum("DrawOrder", name)
}

// MarshalText implements the text marshaller method.
func (x DrawOrder) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *DrawOrder) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := parseDrawOrder(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}

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

// vim: ts=4

