package enhance

import (
  "image"
  "log"
)

func shrinkHorizontal(in image.Image, energy *image.Gray16) (out image.Image, energe *image.Gray16) {
  size := in.Bounds().Size()
  log.Printf("shrinkHorizontal (w=%v, h=%v)\n", size.X, size.Y)
  scores := make([][]uint32, size.Y)
  for y := 0; y < size.Y; y++ {
    scores[y] = make([]uint32, size.X)
  }

  for y := 0; y < size.Y; y++ {
  }

  script := make([]int, size.Y)
  // Erase script
  new_rect := image.Rect(0, 0, size.X - 1, size.Y)
  new_image := image.NewRGBA(new_rect)
  new_energy := image.NewGray16(new_rect)
  for y, deleted_x := range script {
    for x := 0; x < new_rect.Max.X; x++ {
      src_x := x
      if x >= deleted_x { src_x++ }
      new_energy.Set(x, y, energy.At(src_x, y))
      new_image.Set(x, y, in.At(src_x, y))
    }
  }
  return new_image, new_energy
}

func SeamCarve(in image.Image, goal image.Rectangle) image.Image {
  energy := Energy(in)
  for in.Bounds().Size().X >= goal.Size().X {
    in, energy = shrinkHorizontal(in, energy)
  }
  return in
}
