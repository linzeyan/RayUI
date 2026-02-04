// zip_artifact creates a zip archive from a file or directory.
// For .app bundles it preserves the directory structure.
// Usage: go run ./scripts/zip_artifact <output.zip> <source>
package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage: zip_artifact <output.zip> <source>")
	}

	outPath := os.Args[1]
	srcPath := os.Args[2]

	fi, err := os.Stat(srcPath)
	if err != nil {
		log.Fatalf("stat %s: %v", srcPath, err)
	}

	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("create %s: %v", outPath, err)
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	defer w.Close()

	if fi.IsDir() {
		base := filepath.Dir(srcPath)
		err = filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(base, path)
			if err != nil {
				return err
			}

			// Use forward slashes in zip entries.
			rel = filepath.ToSlash(rel)

			if info.IsDir() {
				// Add directory entry.
				_, err := w.Create(rel + "/")
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			header.Name = rel
			header.Method = zip.Deflate

			// Preserve executable permission.
			if info.Mode()&0111 != 0 {
				header.SetMode(info.Mode())
			}

			writer, err := w.CreateHeader(header)
			if err != nil {
				return err
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(writer, f)
			return err
		})
	} else {
		header, err := zip.FileInfoHeader(fi)
		if err != nil {
			log.Fatalf("file info header: %v", err)
		}
		header.Name = fi.Name()
		header.Method = zip.Deflate

		writer, err := w.CreateHeader(header)
		if err != nil {
			log.Fatalf("create header: %v", err)
		}

		f, err := os.Open(srcPath)
		if err != nil {
			log.Fatalf("open %s: %v", srcPath, err)
		}
		defer f.Close()

		_, err = io.Copy(writer, f)
		if err != nil {
			log.Fatalf("copy: %v", err)
		}
	}

	if err != nil {
		log.Fatalf("walk %s: %v", srcPath, err)
	}

	fmt.Printf("Created %s\n", outPath)
}
