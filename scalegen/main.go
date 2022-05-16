package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spacemeshos/go-scale/scalegen/gen"
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

	for _, ev := range []string{"GOARCH", "GOOS", "GOFILE", "GOLINE", "GOPACKAGE", "DOLLAR"} {
		fmt.Println("  ", ev, "=", os.Getenv(ev))
	}

	if err := gen.RunGenerate(os.Getenv("GOPACKAGE"), typesfile, strings.Split(imports, ","), strings.Split(types, ",")); err != nil {
		log.Fatalf("failed to generate: %v", err)
	}
}
