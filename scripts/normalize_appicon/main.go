// normalize_appicon re-encodes an image file as a canonical PNG.
// Usage: go run ./scripts/normalize_appicon <input> <output>
package main

import (
	"image"
	"image/png"
	"log"
	"os"

	_ "image/gif"
	_ "image/jpeg"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage: normalize_appicon <input> <output>")
	}

	in, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	img, _, err := image.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		log.Fatal(err)
	}
}
