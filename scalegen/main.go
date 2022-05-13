package main

import (
	"flag"
	"log"
	"strings"

	"github.com/spacemeshos/go-scale/gen"
)

func main() {
	var (
		pkg       string
		typesfile string
		types     string
		imports   string
	)
	flag.StringVar(&pkg, "pkg", "", "")
	flag.StringVar(&typesfile, "file", "", "")
	flag.StringVar(&types, "types", "", "")
	flag.StringVar(&imports, "imports", "", "")
	flag.Parse()
	if err := gen.RunGenerate(pkg, typesfile, strings.Split(imports, ","), strings.Split(types, ",")); err != nil {
		log.Fatalf("failed to generate: %v", err)
	}
}
