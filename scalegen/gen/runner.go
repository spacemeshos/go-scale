package gen

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
)

const program = `package main

import (
	"log"

	"github.com/spacemeshos/go-scale/scalegen/gen"

	{{ range $pkg := .Imports }}"{{ $pkg }}"
    {{ end }}
)

func main() {
	if err := gen.Generate("{{ .Package }}", "{{ .Output }}", {{ .Objects }}); err != nil {
		log.Fatalf("Generate failed with %v", err)
	}
}
`

type context struct {
	Package string
	Output  string
	Imports []string
	Objects string
}

func getTypes(parsed *ast.File) []string {
	var rst []string
	ast.Inspect(parsed, func(n ast.Node) bool {
		switch typ := n.(type) {
		case *ast.TypeSpec:
			_, ok := typ.Type.(*ast.StructType)
			if !ok {
				return true
			}
			if typ.Name != nil {
				rst = append(rst, typ.Name.String())
			}
		}
		return true
	})
	return rst
}

func getPkg(parsed *ast.File) string {
	return parsed.Name.Name
}

func getModule(in string) (string, error) {
	if in == "/" {
		return "", errors.New("not a module")
	}
	dir := filepath.Dir(in)
	log.Printf("looking for go.mod in %s", dir)
	modf := filepath.Join(dir, "go.mod")
	if f, err := os.Open(modf); err == nil {
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		parsed, err := modfile.Parse(modf, data, nil)
		if err != nil {
			return "", err
		}
		return parsed.Module.Mod.Path, nil
	}
	return getModule(dir)
}

const scaleSuffix = "_scale.go"

func ScaleFile(original string) string {
	ext := filepath.Ext(original)
	base := strings.TrimRight(original, ext)
	return base + scaleSuffix
}

func RunGenerate(in, out string, types []string) error {
	f, err := os.Open(in)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", in, err)
	}
	defer f.Close()

	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, in, f, parser.AllErrors)
	if err != nil {
		return err
	}

	if types == nil {
		types = getTypes(parsed)
		log.Printf("discovered types %+v", types)
	}
	pkg := getPkg(parsed)
	log.Printf("discovered package '%s'", pkg)
	module, err := getModule(in)
	if err != nil {
		return err
	}
	log.Printf("discovered module '%s'", module)

	list := []string{}
	for _, obj := range types {
		list = append(list, fmt.Sprintf("%v.%v{}", pkg, obj))
	}

	ctx := context{
		Package: pkg,
		Output:  out,
		Objects: strings.Join(list, ", "),
		Imports: []string{module + "/" + pkg},
	}
	tpl, err := template.New("").Parse(program)
	if err != nil {
		return err
	}
	now := time.Now()
	programfile := filepath.Join("/tmp/", fmt.Sprintf("scale_gen_%v.go", now.Unix()))
	f, err = os.Create(programfile)
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
	return cmd.Run()
}
