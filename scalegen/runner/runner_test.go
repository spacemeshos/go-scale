package runner

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func examplesDir(tb testing.TB) string {
	tb.Helper()
	rel, err := filepath.Abs("../../examples")
	require.NoError(tb, err)
	return rel
}

func TestGoldenExamples(t *testing.T) {
	dir := examplesDir(t)
	files, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, file := range files {
		file := file
		if strings.Contains(file.Name(), scaleSuffix) || strings.Contains(file.Name(), "test.go") {
			continue
		}
		if file.IsDir() {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			in := filepath.Join(dir, file.Name())
			out := filepath.Join(t.TempDir(), "scale.go")
			require.NoError(t, RunGenerate(in, out, nil))
			golden := filepath.Join(dir, ScaleFile(file.Name()))

			outdata, err := os.ReadFile(out)
			require.NoError(t, err)
			goldendata, err := os.ReadFile(golden)
			require.NoError(t, err)
			require.Equal(t, string(goldendata), string(outdata))
		})
	}
}

func testDataDir(tb testing.TB) string {
	tb.Helper()
	rel, err := filepath.Abs("./testdata")
	require.NoError(tb, err)
	return rel
}

func TestCleanupScaleFile(t *testing.T) {
	dir := testDataDir(t)
	files, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, file := range files {
		file := file
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), scaleSuffix) {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			dataFile := filepath.Join(dir, file.Name())
			scaleFile := filepath.Join(dir, ScaleFile(file.Name()))
			scaleEmptyFile := scaleFile + ".empty"

			scaleFileData, err := os.ReadFile(scaleFile)
			require.NoError(t, err)
			scaleEmptyFileData, err := os.ReadFile(scaleEmptyFile)
			require.NoError(t, err)

			scaleFileCopy := scaleFile + ".copy"
			require.NoError(t, err)
			err = os.WriteFile(scaleFileCopy, scaleFileData, 0o644)
			require.NoError(t, err)
			t.Cleanup(func() { os.Remove(scaleFileCopy) })

			err = cleanupScaleFile(dataFile, scaleFileCopy)
			require.NoError(t, err)

			scaleFileCopyData, err := os.ReadFile(scaleFileCopy)
			require.NoError(t, err)

			require.Equal(t, string(scaleFileCopyData), string(scaleEmptyFileData))
		})
	}
}

func TestScaleFileNoErrorOnGenerateWhenFieldRemoved(t *testing.T) {
	dir := testDataDir(t)
	files, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, file := range files {
		file := file
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), scaleSuffix) {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			typeFile := filepath.Join(dir, file.Name())
			scaleFile := filepath.Join(dir, ScaleFile(file.Name()))

			typeFileData, err := os.ReadFile(typeFile)
			require.NoError(t, err)
			t.Cleanup(func() { restoreFile(typeFile, typeFileData) })

			scaleFileData, err := os.ReadFile(scaleFile)
			require.NoError(t, err)
			t.Cleanup(func() { restoreFile(scaleFile, scaleFileData) })

			err = removeOneFieldInEveryStructType(typeFile)
			require.NoError(t, err)

			require.NoError(t, RunGenerate(typeFile, scaleFile, nil))
		})
	}
}

func restoreFile(file string, body []byte) error {
	return os.WriteFile(file, body, 0o644)
}

func removeOneFieldInEveryStructType(file string) error {
	fIn, err := os.Open(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to open file '%s': %w", file, err)
	}
	defer fIn.Close()

	fset := token.NewFileSet()

	parsed, err := parser.ParseFile(fset, file, fIn, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("failed parsing file '%s': %w", file, err)
	}

	// remove the first field in every struct type
	ast.Inspect(parsed, func(n ast.Node) bool {
		switch typ := n.(type) {
		case *ast.TypeSpec:
			structType, ok := typ.Type.(*ast.StructType)
			if !ok {
				return true
			}
			if len(structType.Fields.List) < 2 {
				return true
			}
			structType.Fields.List = structType.Fields.List[1:]
		}
		return true
	})

	// write modified syntax tree back to the file
	fOut, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to truncate file '%s': %w", file, err)
	}
	defer fOut.Close()

	err = printer.Fprint(fOut, fset, parsed)
	if err != nil {
		return fmt.Errorf("failed writing changes back to file '%s': %w", file, err)
	}

	return nil
}

func TestScaleFileNoErrorOnGenerateWhenTypeRemoved(t *testing.T) {
	dir := testDataDir(t)
	files, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, file := range files {
		file := file
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), scaleSuffix) {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			typeFile := filepath.Join(dir, file.Name())
			scaleFile := filepath.Join(dir, ScaleFile(file.Name()))

			typeFileData, err := os.ReadFile(typeFile)
			require.NoError(t, err)
			t.Cleanup(func() { restoreFile(typeFile, typeFileData) })

			scaleFileData, err := os.ReadFile(scaleFile)
			require.NoError(t, err)
			t.Cleanup(func() { restoreFile(scaleFile, scaleFileData) })

			err = removeFirstTypeDeclaration(typeFile)
			require.NoError(t, err)

			require.NoError(t, RunGenerate(typeFile, scaleFile, nil))
		})
	}
}

func removeFirstTypeDeclaration(file string) error {
	fIn, err := os.Open(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("failed to open file '%s': %w", file, err)
	}
	defer fIn.Close()

	fset := token.NewFileSet()

	parsed, err := parser.ParseFile(fset, file, fIn, parser.AllErrors)
	if err != nil {
		return fmt.Errorf("failed parsing file '%s': %w", file, err)
	}

	filteredDecls := parsed.Decls[:0]
	firstTypeDeclRemoved := false
	for _, decl := range parsed.Decls {
		switch declType := decl.(type) {
		case *ast.GenDecl:
			if !firstTypeDeclRemoved && declType.Tok == token.TYPE {
				firstTypeDeclRemoved = true
				continue
			}
			filteredDecls = append(filteredDecls, declType)
		}
	}

	parsed.Decls = filteredDecls

	// write modified syntax tree back to the file
	fOut, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to truncate file '%s': %w", file, err)
	}
	defer fOut.Close()

	err = printer.Fprint(fOut, fset, parsed)
	if err != nil {
		return fmt.Errorf("failed writing changes back to file '%s': %w", file, err)
	}

	return nil
}

func FuzzScaleFile(f *testing.F) {
	f.Fuzz(func(t *testing.T, pattern string) {
		in := pattern + ".go"
		out := pattern + "_scale.go"

		actual := ScaleFile(in)
		require.Equal(t, out, actual)
	})
}
