package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spacemeshos/go-scale/scalegen/gen"
)

func main() {
	var (
		types    string
		original = os.Getenv("GOFILE")
		split    []string
	)
	flag.StringVar(&types, "types", "", "autodiscovers types if not provided")
	flag.Parse()

	if len(types) > 0 {
		split = strings.Split(types, ",")
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get wd %v", err)
	}
	if err := gen.RunGenerate(
		filepath.Join(wd, original),
		filepath.Join(wd, gen.ScaleFile(original)),
		split,
	); err != nil {
		log.Fatalf("failed to generate: %v", err)
	}
}
