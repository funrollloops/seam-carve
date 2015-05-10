package enhance

import (
  "image"
  "image/color"
  "math"
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

type ImageF32 struct {
  Data []float32
  Stride int
  Width int
  Height int
}

func (i *ImageF32) Offset(x, y int) int {
  return y * i.Stride + x
}

func (i *ImageF32) At(x, y int) float32 {
  return i.Data[i.Offset(x,y)]
}

func (i *ImageF32) Set(x, y int, val float32) {
  i.Data[i.Offset(x, y)] = val
}

func (i *ImageF32) ToGray8() *image.Gray {
  gray := image.NewGray(image.Rect(0, 0, i.Width, i.Height))
  var max_energy float32
  for y := 0; y < i.Height; y++ {
    for x := 0; x < i.Width; x++ {
      if i.At(x, y) > max_energy {
        max_energy = i.At(x, y)
      }
    }
  }
  for y := 0; y < i.Height; y++ {
    for x := 0; x < i.Width; x++ {
      gray.SetGray(x, y, color.Gray{Y: uint8(i.At(x, y) * 255 / max_energy)})
    }
  }
  return gray
}

func NewImageF32(w, h int) *ImageF32 {
  return &ImageF32{Data: make([]float32, w * h), Stride: w, Width: w, Height: h}
}

func EnergyAt(in image.Image, x, y int) float32 {
  minx := x - 1
  maxx := x + 1
  miny := y - 1
  maxy := y + 1
  if minx < in.Bounds().Min.X { minx = in.Bounds().Min.X }
  if maxx >= in.Bounds().Max.X { maxx = in.Bounds().Max.X }
  if miny < in.Bounds().Min.X { miny = in.Bounds().Min.X }
  if maxy >= in.Bounds().Max.X { maxy = in.Bounds().Max.X }
  if maxx - minx < 1 || maxy - miny < 1 { return 0. }

  var avgR, avgG, avgB, avgA uint32
  var count uint32
  for x := minx; x <= maxx; x++ {
    for y := miny; y <= maxy; y++ {
      r, g, b, a := in.At(x, y).RGBA()
      count++
      avgR += r
      avgG += g
      avgB += b
      avgA += a
    }
  }
  avgR /= count
  avgG /= count
  avgB /= count
  avgA /= count

  var variance float64
  sqdiff := func(a, b uint32) float64 { d:= a - b; return float64(d * d); }
  for x := minx; x <= maxx; x++ {
    for y := miny; y <= maxy; y++ {
      r, g, b, a := in.At(x, y).RGBA()
      variance += sqdiff(avgR, r) + sqdiff(avgG, g) + sqdiff(avgB, b) + sqdiff(avgA, a)
    }
  }
  return float32(math.Sqrt(variance))
}

func Energy(in image.Image) *ImageF32 {
  bounds := in.Bounds()
  energy := NewImageF32(bounds.Size().X, bounds.Size().Y)
  for y := bounds.Min.Y + 1; y < bounds.Max.Y - 1; y++ {
    row := energy.Data[energy.Offset(0, y):]
    for x := bounds.Min.X + 1; x < bounds.Max.X - 1; x++ {
      row[x - bounds.Min.X] = EnergyAt(in, x, y)
    }
  }
  return energy
}
