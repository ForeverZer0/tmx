package tmx

// Container describes a type that implements a doubly linked-list of map layers.
type Container interface {
	// Head returns the first layer in a doubly linked-list of layers, or nil when the
	// container is empty.
	Head() Layer
	// Tail returns the last layer in a double linked-list of layers, or nil when the
	// container is empty.
	Tail() Layer
	// Len returns the total number of layers within the container.
	Len() int
	// AddLayer appends a new layer to the container.
	AddLayer(layer Layer)
}

// container is a concrete implementation of the Container interface to be used as a composite
// type for implementing types.
type container struct {
	head Layer
	tail Layer

	// TileLayers is a slice of all tile layers in the container.
	//
	// This is field is exported for convenience, but should not be modified (i.e. append/delete).
	TileLayers []*TileLayer
	// ImageLayers is a slice of all image layers in the container.
	//
	// This is field is exported for convenience, but should not be modified (i.e. append/delete).
	ImageLayers []*ImageLayer
	// ObjectLayers is a slice of all object layers in the container.
	//
	// This is field is exported for convenience, but should not be modified (i.e. append/delete).
	ObjectLayers []*ObjectLayer
	// GroupLayers is a slice of all group layers in the container.
	//
	// This is field is exported for convenience, but should not be modified (i.e. append/delete).
	GroupLayers []*GroupLayer
}

// Head returns the first layer in a doubly linked-list of layers, or nil when empty.
func (c *container) Head() Layer {
	return c.head
}

// Tail returns the last layer in a doubly linked-list of layers, or nil when the empty.
func (c *container) Tail() Layer {
	return c.tail
}

// Len returns the number of layers in the container.
func (c *container) Len() int {
	return len(c.TileLayers) + len(c.ImageLayers) + len(c.ObjectLayers) + len(c.GroupLayers)
}

// AddLayer appends a new layer to the container.
func (g *container) AddLayer(layer Layer) {
	switch v := layer.(type) {
	case *TileLayer:
		g.TileLayers = append(g.TileLayers, v)
	case *ImageLayer:
		g.ImageLayers = append(g.ImageLayers, v)
	case *ObjectLayer:
		g.ObjectLayers = append(g.ObjectLayers, v)
	case *GroupLayer:
		g.GroupLayers = append(g.GroupLayers, v)
	}

	if g.head == nil {
		g.head = layer
	}

	if g.tail != nil {
		g.tail.setNext(layer)
		layer.setPrev(g.tail)
	}
	g.tail = layer
	g.head.setParent(g.parent)
	g.head.setContainer(g)
}

// vim: ts=4
