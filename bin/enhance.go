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
  output_filename = flag.String("output", "", "Image to process")
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
  img, _, err := image.Decode(file)
  if err != nil { log.Fatal(err) }
  return img
}

func main() {
  flag.Parse()
  img := readImage(*filename)
  new_size := image.Rect(0, 0, img.Bounds().Size().X - 20, img.Bounds().Size().Y)
  carved := enhance.SeamCarve(img, new_size)
  writeImage(carved, *output_filename)
}
