package main

import (
	"image"
	"image/color"

	"github.com/hexaflex/wireworld-gpu/math"
)

// Internal cell state values used by the simulation shader.
// These need to stay in sync with the constants in the `shared`
// shader source in shader_shared.go
const (
	CellEmpty = 0
	CellWire  = 50
	CellTail  = 100
	CellHead  = 255
)

// Palette defines the color pallette to use for cell states.
type Palette struct {
	Empty color.RGBA
	Wire  color.RGBA
	Head  color.RGBA
	Tail  color.RGBA
}

// LoadDefault sets the palette to its default values.
func (p *Palette) LoadDefault() {
	p.Empty = color.RGBA{0x00, 0x00, 0x00, 0xff}
	p.Wire = color.RGBA{0x01, 0x5b, 0x96, 0xff}
	p.Head = color.RGBA{0xff, 0xff, 0xff, 0xff}
	p.Tail = color.RGBA{0x99, 0xff, 0x00, 0xff}
}

// fromInternalFormat converts the given 8bpp pixel buffer into an RGBA image
// with colors from the pallette.
func (p *Palette) fromInternalFormat(pix []byte, size math.Vec2) image.Image {
	w, h := int(size[0]), int(size[1])
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			switch pix[y*w+x] {
			case CellWire:
				img.Set(x, y, p.Wire)
			case CellHead:
				img.Set(x, y, p.Head)
			case CellTail:
				img.Set(x, y, p.Tail)
			default:
				img.Set(x, y, p.Empty)
			}
		}
	}

	return img
}

// toInternalFormat converts the given image to its 8bpp internal equivalent.
// Returns the pixel data and dimensions.
func (p *Palette) toInternalFormat(img image.Image) ([]byte, math.Vec2) {
	b := img.Bounds()
	out := image.NewAlpha(b)

	for y := b.Min.Y; y <= b.Max.Y; y++ {
		for x := b.Min.X; x <= b.Max.X; x++ {
			c := p.toCellState(img.At(x, y))
			out.Set(x, y, c)
		}
	}

	return out.Pix, math.Vec2{float32(b.Dx()), float32(b.Dy())}
}

// toCellState translates color c to its internal simulation representation.
func (p *Palette) toCellState(c color.Color) color.Color {
	switch {
	case colorEquals(p.Wire, c):
		return color.Alpha{CellWire}
	case colorEquals(p.Head, c):
		return color.Alpha{CellHead}
	case colorEquals(p.Tail, c):
		return color.Alpha{CellTail}
	default: // All other colors are treated as empty cells.
		return color.Alpha{CellEmpty}
	}
}

// colorEquals returns true if the two colors have the same component values.
func colorEquals(a, b color.Color) bool {
	ar, ag, ab, _ := a.RGBA()
	br, bg, bb, _ := b.RGBA()
	return ar == br && ag == bg && ab == bb
}
