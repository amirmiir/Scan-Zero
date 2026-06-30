// Package scannercore is the JNI bridge surface for the Scan-Zero Go core.
//
// Serialized Primitive Protocol:
//   - Inputs:    a plain string, or a JSON-encoded array of absolute paths.
//   - Responses: a single primitive (string), or a thrown error.
//     (An exported func whose last return is error becomes a Java method
//     that throws.)
//   - All file IO is confined to the caller-supplied cache directory.
package scannercore

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"

	"github.com/amirmiir/scannercore/internal/imaging"
)

func Ping(input string) string {

	return "pong: " + input
}

func unpackPaths(input string) ([]string, error) {
	//string input goes to json.Unmarshal
	//json.Unmarshal goes into a []string, and handle error cases.
	var a []string
	var err error
	return a, err
}

func ProcessImage(inputPath string, outputPath string) (string, error) {
	inputImage, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("Failed to open file: %w", err)
	}
	defer inputImage.Close()

	img, format, err := image.Decode(inputImage)
	if err != nil {
		return "", fmt.Errorf("Failed to decode image: %w", err)
	}
	if format != "png" && format != "jpeg" {
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	grayImg := imaging.ToGray(img)
	t := imaging.Threshold(grayImg)
	resultImg := imaging.ApplyContrastAnchor(grayImg, t, 10, 10)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("create output: %w", err)
	}
	defer outFile.Close()

	w := bufio.NewWriter(outFile)
	if err := jpeg.Encode(w, resultImg, &jpeg.Options{Quality: 95}); err != nil {
		return "", fmt.Errorf("encode: %w", err)
	}
	if err := w.Flush(); err != nil {
		return "", fmt.Errorf("flush: %w", err)
	}
	return outputPath, nil
}
