package enhance

import (
  "image"
  "image/draw"
  "log"
)

var _ = log.Printf

const infinity = 0xffffffff

func shrinkHorizontal(in *image.RGBA, energy *ImageF32) (*image.RGBA, *ImageF32) {
  size := in.Bounds().Size()
  if energy.Width != size.X || energy.Height != size.Y {
    log.Fatalf("Energy.size = %v,%v != in.size = %v", energy.Width, energy.Height, size)
  }

  if size.X < 2 { return in, energy }
  type Cell struct {
    Score float32
    Prev int
  }
  scores := make([][]Cell, size.Y)
  for y := 0; y < size.Y; y++ {
    scores[y] = make([]Cell, size.X)
  }

  // Initialize first row.
  for x := 0; x < size.X; x++ {
    scores[0][x].Score = energy.At(x, 0)
  }
  // Generate scores.
  for y := 1; y < size.Y; y++ {
    last := scores[y - 1]
    current := scores[y]
    if last[0].Score < last[1].Score {
      current[0] = Cell{last[0].Score, 0}
    } else {
      current[0] = Cell{last[1].Score, 1}
    }
    for x := 1; x < size.X - 1; x++ {
      switch {
      case last[x-1].Score <= last[x].Score && last[x-1].Score <= last[x + 1].Score:
        current[x].Prev = x - 1
      case last[x].Score <= last[x + 1].Score:
        current[x].Prev = x
      default:
        current[x].Prev = x + 1
      }
      current[x].Score = last[current[x].Prev].Score + energy.At(x, y)
    }
    if last[size.X - 2].Score < last[size.X - 1].Score {
      current[size.X - 1] = Cell{last[size.X - 2].Score, size.X - 2}
    } else {
      current[size.X - 1] = Cell{last[size.X - 1].Score, size.X - 1}
    }
  }
  // Generate erase script.
  script := make([]int, size.Y)
  // Infer least-cost path.
  var deleted_x int
  for x := 1; x < size.X; x++ {
    if scores[size.Y - 1][x].Score < scores[size.Y - 1][deleted_x].Score {
      deleted_x = x
    }
  }
  for y := size.Y - 1; y > 0; y-- {
    script[y] = deleted_x
    deleted_x = scores[y][deleted_x].Prev
  }
  script[0] = deleted_x

  // Execute erase script.
  for y, deleted_x := range script {
    start := energy.Offset(deleted_x, y)
    end := energy.Offset(size.X, y)
    copy(energy.Data[start : end - 1], energy.Data[start + 1 : end])

    start = in.PixOffset(deleted_x, y)
    end = in.PixOffset(size.X, y)
    copy(in.Pix[start : end - 4], in.Pix[start + 4 : end]) // 4 b/c it's 32 bit
  }
  energy.Width -= 1
  in.Rect.Max.X -= 1
  return in, energy
}

func SeamCarve(in image.Image, w int, h int) image.Image {
  energy := Energy(in)
  out := image.NewRGBA(in.Bounds())
  draw.Draw(out, in.Bounds(), in, in.Bounds().Min, draw.Src)
  for out.Bounds().Size().X >= w {
    out, energy = shrinkHorizontal(out, energy)
  }
  return out
}
