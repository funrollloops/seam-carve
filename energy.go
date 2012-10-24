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

func Energy(in image.Image) *image.Gray16 {
  bounds := in.Bounds()
  energyAt := func(minx, maxx, miny, maxy int) float32 {
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

  var max_energy float32
  as_float := make([][]float32, bounds.Size().Y)
  for y := bounds.Min.Y + 1; y < bounds.Max.Y - 1; y++ {
    row := make([]float32, bounds.Size().X)
    as_float[y - bounds.Min.Y] = row
    for x := bounds.Min.X + 1; x < bounds.Max.X - 1; x++ {
      e := energyAt(x - 1, x + 1, y - 1, y + 1)
      if e > max_energy { max_energy = e }
      row[x - bounds.Min.X] = e
    }
  }

  energy := image.NewGray16(bounds)
  for y, row := range as_float {
    for x, e := range row {
      energy.SetGray16(x, y, color.Gray16{uint16(float32((1<<16) - 1) * e / max_energy)})
    }
  }
  return energy
}
