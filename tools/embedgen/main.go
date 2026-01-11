package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: embedgen <out.go> <src1> [src2 ...]\n")
	os.Exit(2)
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	outPath := os.Args[1]
	srcs := os.Args[2:]

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create output dir: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	w := func(s string) {
		if _, err := io.WriteString(f, s); err != nil {
			fmt.Fprintf(os.Stderr, "write failed: %v\n", err)
			os.Exit(1)
		}
	}

	w("package main\n\n")
	w("import (\n\t\"encoding/base64\"\n\t\"bytes\"\n\t\"io/fs\"\n\t\"time\"\n)\n\n")

	w(`// In-memory file implementation
type memFile struct {
	name string
	data []byte
	r *bytes.Reader
}

func (f *memFile) Read(p []byte) (int, error) { return f.r.Read(p) }
func (f *memFile) Close() error { return nil }
func (f *memFile) Stat() (fs.FileInfo, error) { return memFileInfo{name: f.name, size: int64(len(f.data))}, nil }

// Minimal FileInfo implementation
type memFileInfo struct { name string; size int64 }
func (fi memFileInfo) Name() string { return fi.name }
func (fi memFileInfo) Size() int64  { return fi.size }
func (fi memFileInfo) Mode() fs.FileMode { return 0444 }
func (fi memFileInfo) ModTime() time.Time { return time.Time{} }
func (fi memFileInfo) IsDir() bool { return false }
func (fi memFileInfo) Sys() interface{} { return nil }

// In-memory FS
type memFS struct { files map[string][]byte }

func (m memFS) Open(name string) (fs.File, error) {
	if data, ok := m.files[name]; ok {
		mf := &memFile{name: name, data: data, r: bytes.NewReader(data)}
		return mf, nil
	}
	return nil, fs.ErrNotExist
}

func decode(s string) []byte {
	b, _ := base64.StdEncoding.DecodeString(s)
	return b
}

var staticFiles = memFS{files: map[string][]byte{
`)

	for _, src := range srcs {
		bname := filepath.Base(src)
		data, err := os.ReadFile(src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read %s: %v\n", src, err)
			os.Exit(1)
		}
		enc := base64.StdEncoding.EncodeToString(data)
		line := fmt.Sprintf("  \"dist/%s\": decode(\"%s\"),\n", bname, enc)
		w(line)
	}

	w("}}\n")
}
