package tmx

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

// tmxFrame is an internal type used to marshal the data, and is the composite type of the
// public Frame struct. It merely provides an easy way to marshal, and adjustment of the
// duration to an idiomatic Go duration.
type tmxFrame struct {
	// ID is the local tile ID to display during this frame.
	ID TileID `json:"tileid" xml:"tileid,attr"`
	// Duration is the length of time this frame should be displayed before incrementing
	// to the next frame in the animation.
	Duration time.Duration `json:"duration" xml:"duration,attr"`
}

// Frame describes a single frame within an animation.
type Frame struct {
	tmxFrame // unexported
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (f *Frame) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Simply delegate to composite type, then adjust the duration
	if err := d.DecodeElement(&f.tmxFrame, &start); err != nil {
		return err
	}
	f.Duration *= time.Millisecond
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *Frame) UnmarshalJSON(data []byte) error {
	// Simply delegate to composite type, then adjust the duration
	if err := json.Unmarshal(data, &f.tmxFrame); err != nil {
		return err
	}
	f.Duration *= time.Millisecond
	return nil
}

// vim: ts=4
