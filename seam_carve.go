package enhance

import (
  "image"
  "image/color"
  "log"
)

var _ = log.Printf

const infinity = 0xffffffff

func minIndex(row []uint32) int {
  var min int
  for i := 1; i < len(row); i++ {
    if row[i] < row[min] {
      min = i
    }
  }
  return min
}


func shrinkHorizontal(in image.Image, energy *image.Gray16) (image.Image, *image.Gray16) {
  size := in.Bounds().Size()
  if size.X < 2 { return in, energy }
  //log.Printf("shrinkHorizontal (w=%v, h=%v)\n", size.X, size.Y)
  scores := make([][]uint32, size.Y)
  for y := 0; y < size.Y; y++ {
    scores[y] = make([]uint32, size.X)
  }

  e := func(x, y int) uint32 {
    return uint32(energy.At(x, y).(color.Gray16).Y)
  }

  min2 := func(a, b uint32) uint32 {
    if a < b { return a }
    return b
  }

  // Generate erase script.
  script := make([]int, size.Y)
  // Initialize first row.
  for x := 0; x < size.X; x++ {
    scores[0][x] = e(x, 0)
  }
  // Generate scores.
  for y := 1; y < size.Y; y++ {
    last_row := scores[y - 1]
    scores[y][0] = min2(last_row[0], last_row[1]) + e(0, y)
    for x := 1; x < size.X - 1; x++ {
      scores[y][x] = min2(min2(last_row[x - 1], last_row[x]), last_row[x + 1]) + e(x, y)
    }
    scores[y][size.X - 1] = min2(last_row[size.X - 2], last_row[size.X - 1]) + e(size.X - 1, y)
  }
  // Infer least-cost path.
  deleted_x := minIndex(scores[size.Y - 1])
  for y := size.Y - 1; y > 0; y-- {
    script[y] = deleted_x
    s := scores[y][deleted_x] - e(deleted_x, y)
    if deleted_x > 0 && scores[y - 1][deleted_x - 1] == s {
      deleted_x--
    } else if deleted_x < size.X - 1 && scores[y - 1][deleted_x + 1] == s {
      deleted_x++
    } else if scores[y - 1][deleted_x] != s {
      log.Fatal("Scores matrix is inconsistent for y=%v, deleted_x=%y\n", y, deleted_x)
    }
  }
  script[0] = deleted_x



  // Execute erase script.
  new_rect := image.Rect(0, 0, size.X - 1, size.Y)
  new_image := image.NewRGBA(new_rect)
  for y, deleted_x := range script {
    for x := 0; x < new_rect.Max.X; x++ {
      src_x := x
      if x >= deleted_x { src_x++ }
      new_image.Set(x, y, in.At(src_x, y))
    }
    start := energy.PixOffset(deleted_x, y)
    end := energy.PixOffset(size.X, y)
    copy(energy.Pix[start : end - 2], energy.Pix[start + 2 : end]) // 2 b/c it's 16 bit
  }
  return new_image, energy
}

func SeamCarve(in image.Image, w int, h int) image.Image {
  energy := Energy(in)
  for in.Bounds().Size().X >= w {
    in, energy = shrinkHorizontal(in, energy)
  }
  return in
}
