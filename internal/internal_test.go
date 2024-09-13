package internal

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	inputFileName    = "input.json"
	expectedFileName = "expected.go"
)

func dirs(t *testing.T, path string) []string {
	t.Helper()

	var result []string

	err := filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.Name() == inputFileName {
			result = append(result, filepath.Dir(path))
		}

		return nil
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)

	return result
}

func parseInput[T any](t *testing.T, dir string) *T {
	t.Helper()

	f, err := os.Open(filepath.Join(dir, inputFileName))
	require.NoError(t, err)

	defer f.Close()

	result := new(T)

	require.NoError(t, json.NewDecoder(f).Decode(result))

	return result
}

func toFile(t *testing.T, dir string, stmt *jen.Statement) *jen.File {
	t.Helper()

	result := jen.NewFilePath(dir)

	result.Add(stmt)

	return result
}

func render(t *testing.T, f *jen.File) string {
	t.Helper()

	sb := &strings.Builder{}

	require.NoError(t, f.Render(sb))

	return sb.String()
}

func compare(t *testing.T, dir string, stmt *jen.Statement) {
	t.Helper()

	expected, err := os.ReadFile(filepath.Join(dir, expectedFileName))
	require.NoError(t, err)

	assert.Equal(t, string(expected), render(t, toFile(t, dir, stmt)))
}
