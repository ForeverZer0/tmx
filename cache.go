package tmx

// Cache provides a mechanism for maintaining references that are shared among multiple
// objects or that will be used frequently.
type Cache struct {
	tilesets  map[string]*Tileset
	templates map[string]*Template
}

// NewCache initializes and returns a new Cache.
func NewCache() *Cache {
	return &Cache{
		tilesets:  make(map[string]*Tileset),
		templates: make(map[string]*Template),
	}
}

// Tileset retrieves a Tileset from the cache, or nil if it was not found.
func (c *Cache) Tileset(key string) (*Tileset, bool) {
	if value, ok := c.tilesets[key]; ok {
		return value, true
	}
	return nil, false
}

// Tileset retrieves a Template from the cache, or nil if it was not found.
func (c *Cache) Template(key string) (*Template, bool) {
	if value, ok := c.templates[key]; ok {
		return value, true
	}
	return nil, false
}

// AddTileset adds a new Tileset to the cache with the given key, returning
// a value if it was successfully added.
//
// If the key already exists in the cache, the operation will fail and
// return false.
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

// AddTemplate adds a new Template to the cache with the given key, returning
// a value if it was successfully added.
//
// If the key already exists in the cache, the operation will fail and
// return false.
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

// Clear removes all values from the Cache, allowing them to be
// garbage collected.
func (c *Cache) Clear() {
	c.tilesets = make(map[string]*Tileset)
	c.templates = make(map[string]*Template)
}

// vim: ts=4
