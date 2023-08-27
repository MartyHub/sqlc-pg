package internal

import (
	"path/filepath"
	"strings"

	"github.com/MartyHub/sqlc-pg/plugin"
)

func fileNameWithoutExt(path string) string {
	fileName := filepath.Base(path)
	name, _ := strings.CutSuffix(fileName, filepath.Ext(fileName))

	return name
}

func (gen *Generator) scanName(str StructMetadata) string {
	return gen.tok.ExportID("Scan_" + str.Name)
}

func (gen *Generator) stmtName(query *plugin.Query) string {
	return gen.tok.ToCamel(query.Name+"_Stmt", gen.cfg.EmitExportedQueries)
}

func (gen *Generator) fileName(s string) string {
	sb := strings.Builder{}

	for _, tok := range gen.tok.Tokens(s) {
		if !tok.Valid() {
			continue
		}

		if sb.Len() > 0 {
			sb.WriteRune('_')
		}

		sb.WriteString(string(tok.Runes))
	}

	return sb.String()
}
