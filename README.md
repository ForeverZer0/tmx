# tmx

A Go library for reading Tiled TMX format, both XML and JSON formats.

## Features

* Seamlessly reads both XML/JSON formats into a single interface
* Full-coverage of the [TMX specification](https://doc.mapeditor.org/en/stable/reference/tmx-map-format/#) (v1.10)
* Supports Gzip, Zlib, and Zstd compression formats
* Intuitive API surface with idiomatic Go style
* Optional "callback" for providing realtime image/texture processing, including embedded images
* Hook for providing user-defined resolution of file paths and/or "virtual" filesystem (useful for embeddeded/binary documents)
* Convenience functions for basic tasks, including retrieving tiles from map locations (including infinite maps), UV texture coordinates, resource caching, and more.

## Usage


### Reading Files

The main entry point to loading a map is simply the `ReadMap` function:

```go
tilemap, err := tmx.ReadMap("path/to/map.tmx", FormatUnknown, nil)
if err != nil {
    // handle error
}
```

An optional `Cache` object can be created and passed into the read functions. This is useful to
automaatically maintain references to shared resources, such as tilesets and templates.

For example, a region in your game with multiple map "screens" that all use the same tileset. The
cache can be used and passed in for each map being transitioned to, and it will not require
re-processing the tileset each time, and will use the existing instance in the cache.

```go
cache := tmx.NewCache()

tilemap, err := tmx.ReadMap("../maps/lofty_mountains_east.tmj", FormatJSON, cache)

// ...the player moves to the next area...

// This time, if the both maps use the same tileset, the instance in the cache
// will be reused instead of reloading it from disk.
tilemap, err = tmx.ReadMap("../maps/lofty_mountains_summit.tmj", FormatJSON, cache)

```

There are additionally the `ReadTileset` and `ReadTemplate` functions for reading those types if
needed.

### Layers

There are multiple ways to iterate through the layers, allowing you to choose the best method
for your project's specific needs.

#### By Layer Type

Each type of layer is stored within a slice of that type, which are exported fields and readily
available for iteration. This makes "filtering by type" quite trivial, as they are already stored
internally in such a way.


```go
for _, tileLayer := range tilemap.TileLayers {
    // Perform TileLayer-specific operations
}
```

#### Ordered

The most common approach to iterating through layers in 2D is of course to do so in an ordered
fashion. The `Container` type provides an interfaces which is implmented by the both the `Map` and
`GroupLayer` types, which allows for doubly-linked list functionality.

To iterate from the "bottom-most" layer to the "top-most":

```go
var layer Layer // Common interface for all layer types

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
```

Alternatively, you can enumerate in the reverse, and move top to bottom:

```go
for layer := tilemap.Tail(); layer != nil; layer = layer.Prev() {

}
```
