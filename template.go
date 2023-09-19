package tmx

type Template struct {
	Tileset *Tileset `json:"tileset" xml:"tileset"`
	Object  `json:"object" xml:"object"`
}

type Tileset struct{}

func OpenTX(path string) (*Template, error) {
	return nil, nil
}

// vim: ts=4
