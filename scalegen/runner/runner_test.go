package runner

import (
	"io/ioutil"
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
	files, err := ioutil.ReadDir(dir)
	require.NoError(t, err)

	for _, file := range files {
		file := file
		if strings.Contains(file.Name(), scaleSuffix) || strings.Contains(file.Name(), "test.go") || !strings.Contains(file.Name(), ".go") {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			in := filepath.Join(dir, file.Name())
			out := filepath.Join(t.TempDir(), "scale.go")
			require.NoError(t, RunGenerate(in, out, nil))
			golden := filepath.Join(dir, ScaleFile(file.Name()))

			outdata, err := ioutil.ReadFile(out)
			require.NoError(t, err)
			goldendata, err := ioutil.ReadFile(golden)
			require.NoError(t, err)
			require.Equal(t, string(goldendata), string(outdata))
		})
	}
}
