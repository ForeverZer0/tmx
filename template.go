package tmx

type Template struct {
	Tileset *Tileset `json:"tileset" xml:"tileset"`
	Object  `json:"object" xml:"object"`
}

func OpenTemplate(path string, cache *Cache) (*Template, error) {
	// TODO
	return nil, nil
}

// vim: ts=4
