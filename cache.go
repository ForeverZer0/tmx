package tmx

type Cache struct {
	tilesets  map[string]*Tileset
	templates map[string]*Template
}

func NewCache() *Cache {
	return &Cache{
		tilesets:  make(map[string]*Tileset),
		templates: make(map[string]*Template),
	}
}

func (c *Cache) Tileset(key string) (*Tileset, bool) {
	if value, ok := c.tilesets[key]; ok {
		return value, true
	}
	return nil, false
}

func (c *Cache) Template(key string) (*Template, bool) {
	if value, ok := c.templates[key]; ok {
		return value, true
	}
	return nil, false
}

func (c *Cache) AddTileset(key string, tileset *Tileset) bool {
	if tileset == nil {
		return false
	}
	if _, ok := c.tilesets[key]; ok {
		return false
	}
	c.tilesets[key] = tileset
	return true
}

func (c *Cache) AddTemplate(key string, template *Template) bool {
	if template == nil {
		return false
	}
	if _, ok := c.templates[key]; ok {
		return false
	}
	c.templates[key] = template
	return true
}

func (c *Cache) Clear() {
	c.tilesets = make(map[string]*Tileset)
	c.templates = make(map[string]*Template)
}

// vim: ts=4
