package enhance

import (
  "image"
  "image/color"
)

type RGBA struct {
  r, g, b, a uint32
}

func (l *RGBA) Minus(r *RGBA) uint32 {
  return (u32m(l.r, r.r) + u32m(l.g, r.g) + u32m(l.b, r.b)) / 3
}

func u32m(l, r uint32) uint32 {
  if l > r {
    return l - r
  }
  return r - l
}

func rgba(r, g, b, a uint32) *RGBA {
  return &RGBA{r, g, b, a}
}

func Energy(in image.Image) *image.Gray16 {
  bounds := in.Bounds()
  energy := image.NewGray16(bounds)
  for x := bounds.Min.X + 1; x < bounds.Max.X - 1; x++ {
    for y := bounds.Min.Y + 1; y < bounds.Max.Y - 1; y++ {
      dx := rgba(in.At(x - 1, y).RGBA()).Minus(rgba(in.At(x + 1, y).RGBA()))
      dy := rgba(in.At(x, y - 1).RGBA()).Minus(rgba(in.At(x, y + 1).RGBA()))
      energy.SetGray16(x, y, color.Gray16{uint16((dx + dy) / 2)})
    }
  }
  return energy
}
