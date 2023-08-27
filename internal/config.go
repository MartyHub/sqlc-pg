package internal

import (
	"encoding/json"
	"path/filepath"

	"github.com/MartyHub/sqlc-pg/plugin"
)

type Config struct {
	DumpInput bool `json:"dump_input,omitempty"`

	EmitDBTags               bool `json:"emit_db_tags,omitempty"`
	EmitExportedQueries      bool `json:"emit_exported_queries,omitempty"`
	EmitResultStructPointers bool `json:"emit_result_struct_pointers,omitempty"`
	EmitParamsStructPointers bool `json:"emit_params_struct_pointers,omitempty"`
	EmitTableNames           bool `json:"emit_table_names,omitempty"`

	OutputDBFileName  string `json:"output_db_file_name,omitempty"`
	OutputFilesSuffix string `json:"output_files_suffix,omitempty"`
	Package           string `json:"package,omitempty"`
}

func newConfig(req *plugin.CodeGenRequest) (Config, error) {
	var (
		err    error
		result Config
	)

	if len(req.PluginOptions) > 0 {
		if err = json.Unmarshal(req.PluginOptions, &result); err != nil {
			return result, err
		}
	}

	if result.Package == "" {
		if result.Package, err = pkg(req.Settings.Codegen.Out); err != nil {
			return result, err
		}
	}

	if result.OutputDBFileName == "" {
		result.OutputDBFileName = "db" + result.OutputFilesSuffix + ".go"
	}

	return result, nil
}

func pkg(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.Base(absPath), nil
}
