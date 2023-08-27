package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MartyHub/sqlc-pg/plugin"
)

func Test_sqlc_examples(t *testing.T) {
	t.Parallel()

	for _, dir := range dirs(t, filepath.Join("testdata", "sqlc", "examples")) {
		dir := dir

		t.Run(dir, func(t *testing.T) {
			t.Parallel()

			req := parseInput[plugin.CodeGenRequest](t, dir)

			generator, err := New(req)
			require.NoError(t, err)

			rep, err := generator.Generate()
			require.NoError(t, err)
			require.Greater(t, len(rep.Files), 0)

			for _, file := range rep.Files {
				expected, err := os.ReadFile(filepath.Join(dir, file.Name))
				require.NoError(t, err)

				assert.Equal(t, string(expected), string(file.Contents))
			}
		})
	}
}
