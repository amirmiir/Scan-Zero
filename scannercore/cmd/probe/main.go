package main

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/amirmiir/scannercore/internal/imaging"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: probe <image.jpg|png>")
		os.Exit(1)
	}

	imageInPath := os.Args[1]

	src, err := os.Open(imageInPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open: %v\n", err)
		os.Exit(1)
	}
	defer src.Close()

	img, format, err := image.Decode(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decode: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("decoded %s (%s %dx%d)\n", filepath.Base(imageInPath), format, img.Bounds().Dx(), img.Bounds().Dy())

	gray := imaging.ToGray(img)
	if err := writeJPEG(stagePath(imageInPath, "gray"), gray); err != nil {
		fmt.Fprintf(os.Stderr, "write gray: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("gray -> %s\n", stagePath(imageInPath, "gray"))
}

func stagePath(input, stage string) string {
	ext := filepath.Ext(input)
	base := strings.TrimSuffix(input, ext)
	return base + "_" + stage + ".jpg"
}

func writeJPEG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if err := jpeg.Encode(w, img, &jpeg.Options{Quality: 95}); err != nil {
		return err
	}
	return w.Flush()
}
