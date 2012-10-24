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

func main() {
  flag.Parse()
  file, err := os.Open(*filename)
  if err != nil { log.Fatal(err) }
  img, _, err := image.Decode(file)
  if err != nil { log.Fatal(err) }
  w, err := os.OpenFile(*output_filename, os.O_RDWR | os.O_CREATE, 0664)
  if err != nil { log.Fatal(err) }
  png.Encode(w, enhance.Energy(img))
  w.Close()
}
