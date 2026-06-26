# Scan-Zero
cam scanning application for documentation written in golang - no ads or bloat, simple to use

---

## Testing the Go core

Requires Go 1.23+. All commands run from `scannercore/`.

### Unit tests + benchmark

```bash
# all tests
go test ./internal/imaging/ -v

# benchmark — allocs/op must stay constant regardless of image size
go test ./internal/imaging/ -run XXX -bench BenchmarkToGray -benchmem

# escape analysis — no per-pixel interface boxing should appear
go build -gcflags='-m' ./internal/imaging/ 2>&1
```

### Visual probe — inspect pipeline output on a real image

```bash
go run ./cmd/probe/ path/to/photo.jpg
```

Output is written next to the input file:

```
photo.jpg        ← original (untouched)
photo_gray.jpg   ← grayscale output
```

Accepts JPEG and PNG. More output stages are added as the pipeline grows.
**Note:** probe output files are not committed — add them to `.gitignore` or keep them outside the repo.
