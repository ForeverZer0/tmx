# Features

* Seamlessly reads both XML/JSON formats into a single interface
* Full-coverage of the [TMX specification](https://doc.mapeditor.org/en/stable/reference/tmx-map-format/#) (v1.10)
* Supports Gzip, Zlib, and Zstd compression formats
* Intuitive API surface
* Optional "callback" for providing realtime image/texture processing, including embedded images
* Hook for providing user-defined resolution of file paths (useful for embeddeded/binary documents)
* Convenience functions for basic tasks, including retrieving tiles from map locations (including infinite maps), UV texture coordinates, resource caching, and more.

# Usage

The main entry point is the `LoadMap` function:

```go
// Initialize a cache to maintain references for reusable tilesets, templates, etc.
// This is optional, pass nil to have objects get garbage-collected normally.
cache := tmx.NewCache()


tilemap, err := tmx.LoadMap("path/to/map.tmx", FormatXML, cache)
if err != nil {
    // handle error
}

// To access layers of given type, access the slices that contain them directly.

for _, tileLayer := range tilemap.TileLayers {
    // process tile layers
    for _, chunk := range tileLayer.Chunks {
        ...
    }
} 

for _, imageLayer := range tilemap.ImageLayers {
    // process image layers

    if imageLayer.Image != nil {
        
    }
}

// Alternatively, layers are also implemented as a doubly-linked list, which allows
// all layers to be accessed in order, regardless of type.

var layer Layer // Common interface for different layer types

for layer = tilemap.Head(); layer != nil; layer = layer.Next() {

    switch value := layer.(type) {
    case *TileLayer:
        // Do things with tile-based layers
    case *ImageLayer:
        // Do things with image layers
    case *ObjectLayer:
        // Do things with layers of map objects
    case *GroupLayer:
        // Do things with groups of layers
    }
}

for layer := tilemap.Tail(); layer != nil; layer = layer.Prev() {

}

```
