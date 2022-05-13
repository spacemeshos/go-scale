package gen

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const program = `package main

import (
	"log"

	"github.com/spacemeshos/go-scale/gen"

	{{ range $pkg := .Imports }}"{{ $pkg }}"
    {{ end }}
)

func main() {
	if err := gen.Generate("{{ .Package }}", "{{ .File }}", {{ .Objects }}); err != nil {
		log.Fatalf("Generate failed with %v", err)
	}
}
`

type context struct {
	Package string
	File    string
	Imports []string
	Objects string
}

func RunGenerate(pkg string, typesfile string, imports []string, objects []string) error {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get wd %v", err)
	}

	initialized := []string{}
	for _, obj := range objects {
		initialized = append(initialized, fmt.Sprintf("%v.%v{}", pkg, obj))
	}
	ctx := context{
		Package: pkg,
		File:    filepath.Join(wd, typesfile),
		Objects: strings.Join(initialized, ", "),
		Imports: imports,
	}
	tpl, err := template.New("").Parse(program)
	if err != nil {
		return err
	}
	now := time.Now()
	programfile := filepath.Join("/tmp/", fmt.Sprintf("scale_gen_%v.go", now.Unix()))
	f, err := os.Create(programfile)
	if err != nil {
		return err
	}
	defer f.Close()
	log.Printf("program file: %v", programfile)
	defer os.Remove(f.Name())

	if err := tpl.Execute(f, ctx); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	cmd := exec.Command("go", "run", programfile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = wd
	return cmd.Run()
}
