package runner

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

const program = `package main

import (
	"log"

	"github.com/spacemeshos/go-scale"

	{{ range $pkg := .Imports }}"{{ $pkg }}"
    {{ end }}
)

func main() {
	if err := scale.Generate("{{ .Package }}",` + " `{{ .Output }}`" + `, {{ .Objects }}); err != nil {
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

func getModule(in string, parts []string) (string, error) {
	if in == "/" {
		return "", errors.New("not a module")
	}
	dir := filepath.Dir(in)
	modf := filepath.Join(dir, "go.mod")
	if f, err := os.Open(modf); err == nil {
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}
		parsed, err := modfile.Parse(modf, data, nil)
		if err != nil {
			return "", err
		}
		parts = append(parts, parsed.Module.Mod.Path)
		for i := 0; i < len(parts)/2; i++ {
			j := len(parts) - 1 - i
			parts[i], parts[j] = parts[j], parts[i]
		}
		parts = parts[:len(parts)-1]
		return strings.Join(parts, "/"), nil
	}
	return getModule(dir, append(parts, filepath.Base(dir)))
}

const scaleSuffix = "_scale.go"

func ScaleFile(original string) string {
	ext := filepath.Ext(original)
	base := strings.TrimSuffix(original, ext)
	return base + scaleSuffix
}

// cleanupScaleFile removes all function bodies in provided scale file leaving the last
// (usually "return ...") statement only. It also removes scale methods for types missing in dataFilePath.
func cleanupScaleFile(dataFilePath, scaleFilePath string) error {
	// get types defained in data file
	dataFile, err := os.Open(dataFilePath)
	if err != nil {
		return fmt.Errorf("failed to open data file '%s': %w", dataFilePath, err)
	}
	defer dataFile.Close()

	dataFileSet := token.NewFileSet()
	dataFileParsed, err := parser.ParseFile(dataFileSet, dataFilePath, dataFile, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("failed parsing data file '%s': %w", dataFilePath, err)
	}

	dataFileTypes := getTypes(dataFileParsed)

	scaleFile, err := os.Open(scaleFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to open scale file '%s': %w", scaleFilePath, err)
	}
	defer scaleFile.Close()

	fset := token.NewFileSet()

	parsed, err := parser.ParseFile(fset, scaleFilePath, scaleFile, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("failed parsing scale file '%s': %w", scaleFilePath, err)
	}

	parsed.Imports = filterImports(parsed.Imports)
	parsed.Decls = filterDecls(parsed.Decls, dataFileTypes)

	// for every method in a scale file leave the last ("return ...") statement only
	ast.Inspect(parsed, func(n ast.Node) bool {
		switch typ := n.(type) {
		case *ast.FuncDecl:
			typ.Body.List = []ast.Stmt{typ.Body.List[len(typ.Body.List)-1]}
		}
		return true
	})

	// write modified syntax tree back to the file
	fOut, err := os.Create(scaleFilePath)
	if err != nil {
		return fmt.Errorf("failed to truncate scale file '%s': %w", scaleFilePath, err)
	}
	defer fOut.Close()

	err = printer.Fprint(fOut, fset, parsed)
	if err != nil {
		return fmt.Errorf("failed writing changes back to scale file '%s': %w", scaleFilePath, err)
	}

	return nil
}

// getReceiver returns receiver type name for a function declaration.
func getReceiver(f *ast.FuncDecl) string {
	if f == nil || f.Recv == nil {
		return ""
	}
	for _, field := range f.Recv.List {
		switch typ := field.Type.(type) {
		case *ast.StarExpr:
			if si, ok := typ.X.(*ast.Ident); ok {
				return si.Name
			}
		case *ast.Ident:
			return typ.Name
		}
	}
	return ""
}

const goScaleImport = `"github.com/spacemeshos/go-scale"`

func filterImports(imports []*ast.ImportSpec) []*ast.ImportSpec {
	newImports := imports[:0]
	for _, imp := range imports {
		if imp.Path.Value != goScaleImport {
			continue
		}
		newImports = append(newImports, imp)
	}

	return newImports
}

// filterDecls removes scale methods for deleted types as well as all imports but go-scale.
func filterDecls(decls []ast.Decl, dataFileTypes []string) []ast.Decl {
	typesIndex := make(map[string]struct{}, len(dataFileTypes))
	for _, t := range dataFileTypes {
		typesIndex[t] = struct{}{}
	}

	newDecls := decls[:0]
	for _, decl := range decls {
		switch declType := decl.(type) {
		case *ast.FuncDecl:
			// skip scale method if receiver type is not defined in data file
			receiverTypeName := getReceiver(declType)
			if receiverTypeName == "" {
				panic("receiver can't be empty")
			}
			if _, exists := typesIndex[receiverTypeName]; exists {
				newDecls = append(newDecls, declType)
			}

		case *ast.GenDecl:
			if len(declType.Specs) == 0 {
				continue
			}
			newSpecs := declType.Specs[:0]
			for _, spec := range declType.Specs {
				switch specType := spec.(type) {
				case *ast.ImportSpec:
					if specType.Path.Value != goScaleImport {
						continue
					}
					newSpecs = append(newSpecs, spec)
				default:
					newSpecs = append(newSpecs, spec)
				}
			}
			declType.Specs = newSpecs
			newDecls = append(newDecls, declType)
		default:
			newDecls = append(newDecls, declType)
		}
	}

	return newDecls
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
	}
	pkg := getPkg(parsed)
	module, err := getModule(in, nil)
	if err != nil {
		return err
	}

	list := []string{}
	for _, obj := range types {
		list = append(list, fmt.Sprintf("%v.%v{}", pkg, obj))
	}

	// replace all scale methods with empty ones to be sure it has no compile errors after receiver type changed
	err = cleanupScaleFile(in, out)
	if err != nil {
		return err
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
	f, err = os.CreateTemp("", "scale_gen_*.go")
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())

	if err := tpl.Execute(f, ctx); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	cmd := exec.Command("go", "run", f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
