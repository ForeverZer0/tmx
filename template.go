package tmx

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
)

type Template struct {
	Source string `json:"-" xml:"-"`
	Tileset *MapTileset `json:"tileset" xml:"tileset"`
	Object  `json:"object" xml:"object"`
	cache *Cache
}

func (t *Template) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	token, err := d.Token()
	for token != start.End() {
		if err != nil {
			return err
		}

		if child, ok := token.(xml.StartElement); ok {
			switch child.Name.Local {
			case "object":
				var obj Object
				obj.cache = t.cache
				if err = d.DecodeElement(&obj, &child); err != nil {
					return err
				}
				t.Object = obj
			case "tileset":
				var ts MapTileset
				ts.cache = t.cache
				if err = d.DecodeElement(&ts, &child); err != nil {
					return err
				}
				t.Tileset = &ts
			default:
				logElem(child.Name.Local, start.Name.Local)
			}
		}

		token, err = d.Token()
	}
	return nil
}

// OpenTileset reads a tileset from a file, using the specified format.
//
// An optional cache can be supplied that maintains references to tilesets and
// templates to prevent frequent re-processing of them.
func OpenTemplate(path string, format Format, cache *Cache) (*Template, error) {
	var abs string
	var err error
	if abs, err = FindPath(path); err != nil {
		return nil, err
	}

	// Check cache
	if cache != nil {
		if template, ok := cache.Template(abs); ok {
			return template, nil
		}
	}

	reader, _, err := getStream(abs)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	IncludePaths = append(IncludePaths, filepath.Dir(abs))
	defer func() { IncludePaths = IncludePaths[:len(IncludePaths)-1] }()

	var template Template
	template.Source = abs
	template.cache = cache

	if err := Decode(reader, format, &template); err != nil {
		return nil, err
	}

	if cache != nil {
		cache.AddTemplate(abs, &template)
	}
	return &template, nil
}

// Decode reads a TMX object from the current position in the reader using
// the specified format, storing the result to the given pointer.
func Decode(r io.Reader, format Format, obj any) error {
	switch format {
	case FormatXML:
		d := xml.NewDecoder(r)
		if err := d.Decode(obj); err != nil {
			return err
		}
	case FormatJSON:
		d := json.NewDecoder(r)
		if err := d.Decode(obj); err != nil {
			return err
		}
	default:
		return errInvalidEnum("Format", fmt.Sprintf("Format(%d)", format))
	}

	return nil
}

// vim: ts=4
