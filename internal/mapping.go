package internal

import (
	"strconv"

	"github.com/MartyHub/sqlc-pg/plugin"
)

const (
	_32  = 32
	_64  = 64
	_int = "int"
)

func (gen *Generator) mapType(column *plugin.Column) TypeMetadata { //nolint:cyclop
	result := TypeMetadata{
		Name:  gen.tok.ExportID(column.Type.Name),
		Array: column.IsArray,
		Ptr:   !column.IsArray && !column.NotNull,
	}

	switch column.Type.Name {
	case "blob",
		"bytea", "pg_catalog.bytea",
		"json", "jsonb":
		result.Name = "byte"
		result.Array = true
	case "boolean", "bool", "pg_catalog.bool":
		result.Name = "bool"
	case "date",
		"timestamp", "pg_catalog.timestamp",
		"timestamptz", "pg_catalog.timestamptz",
		"pg_catalog.time", "pg_catalog.timetz":
		result.Name = "Time"
		result.Path = "time"
	case "float4", "pg_catalog.float4", "real":
		result.Name = "float32"
	case "float8", "pg_catalog.float8", "float", "double precision":
		result.Name = "float64"
	case "int2", "pg_catalog.int2", "smallint", "smallserial", "serial2", "pg_catalog.serial2":
		result.Name = "int16"
	case "int4", "pg_catalog.int4", _int, "integer", "serial", "serial4", "pg_catalog.serial4":
		result.Name = _int32()
	case "int8", "pg_catalog.int8", "bigint", "bigserial", "serial8", "pg_catalog.serial8":
		result.Name = _int64()
	case "varchar", "pg_catalog.varchar", "text", "string", "uuid":
		result.Name = "string"
	}

	return result
}

func _int32() string {
	if strconv.IntSize == _32 {
		return _int
	}

	return "int32"
}

func _int64() string {
	if strconv.IntSize == _64 {
		return _int
	}

	return "int64"
}
