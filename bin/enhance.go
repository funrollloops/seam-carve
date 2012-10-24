package main

import (
  "sagarmittal.com/enhance"

  "flag"
  "image"
  "image/png"
  "log"
  "os"

  _ "image/jpeg"
  _ "image/gif"
)

var (
  filename = flag.String("image", "", "Image to process")
  energy_filename = flag.String("energy", "", "Filename for energy image")
  output_filename = flag.String("output", "", "Image to process")
  xscale = flag.Float64("xscale", 1., "Horizontal scale factor")
  yscale = flag.Float64("yscale", 1., "Vertical scale factor")
)

func writeImage(in image.Image, path string) {
  w, err := os.OpenFile(*output_filename, os.O_RDWR | os.O_CREATE, 0664)
  if err != nil { log.Fatal(err) }
  err = png.Encode(w, in)
  if err != nil { log.Fatal(err) }
  err = w.Close()
  if err != nil { log.Fatal(err) }
}

func readImage(path string) image.Image {
  file, err := os.Open(*filename)
  if err != nil { log.Fatal(err) }
  img, _, err := image.Decode(file)
  if err != nil { log.Fatal(err) }
  return img
}

func main() {
  flag.Parse()
  img := readImage(*filename)
  if *energy_filename != "" {
    log.Printf("Writing energy image to \"%v\"\n", *energy_filename)
    writeImage(enhance.Energy(img), *energy_filename)
  }
  if *xscale < 1 && *xscale > 0 {
    new_width := int(float64(img.Bounds().Size().X) * *xscale)
    new_height := int(float64(img.Bounds().Size().Y) * *yscale)
    log.Printf("Rescaling to w=%v h=%v", new_width, new_height)
    img = enhance.SeamCarve(img, new_width, new_height)
  }
  writeImage(img, *output_filename)
}
