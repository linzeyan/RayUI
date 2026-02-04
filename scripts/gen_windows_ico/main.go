// gen_windows_ico converts a PNG to a minimal ICO file containing 256x256, 48x48, 32x32, 16x16 sizes.
// Usage: go run ./scripts/gen_windows_ico <input.png> <output.ico>
package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/png"
	"log"
	"os"

	_ "image/gif"
	_ "image/jpeg"

	"golang.org/x/image/draw"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage: gen_windows_ico <input.png> <output.ico>")
	}

	in, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	src, _, err := image.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	sizes := []int{256, 48, 32, 16}
	images := make([][]byte, len(sizes))

	for i, sz := range sizes {
		dst := image.NewRGBA(image.Rect(0, 0, sz, sz))
		draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

		var buf bytes.Buffer
		if err := png.Encode(&buf, dst); err != nil {
			log.Fatal(err)
		}
		images[i] = buf.Bytes()
	}

	out, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// ICO header
	binary.Write(out, binary.LittleEndian, uint16(0))          // reserved
	binary.Write(out, binary.LittleEndian, uint16(1))          // type: ICO
	binary.Write(out, binary.LittleEndian, uint16(len(sizes))) // count

	// Calculate offsets: header(6) + entries(16 * count)
	offset := 6 + 16*len(sizes)

	// Directory entries
	for i, sz := range sizes {
		w := uint8(sz)
		h := uint8(sz)
		if sz == 256 {
			w, h = 0, 0 // 0 means 256 in ICO format
		}
		out.Write([]byte{w, h, 0, 0})                                    // width, height, palette, reserved
		binary.Write(out, binary.LittleEndian, uint16(1))                // color planes
		binary.Write(out, binary.LittleEndian, uint16(32))               // bits per pixel
		binary.Write(out, binary.LittleEndian, uint32(len(images[i])))   // image size
		binary.Write(out, binary.LittleEndian, uint32(offset))           // offset
		offset += len(images[i])
	}

	// Image data
	for _, img := range images {
		out.Write(img)
	}
}
