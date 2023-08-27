package internal

import (
	j "github.com/dave/jennifer/jen"
)

type (
	StructMetadata struct {
		Name   string
		Fields []FieldMetadata
	}

	FieldMetadata struct {
		Name       string
		ColumnName string
		Type       TypeMetadata
	}

	TypeMetadata struct {
		Path  string
		Name  string
		Array bool
		Ptr   bool
	}
)

func (str StructMetadata) Empty() bool {
	return len(str.Fields) == 0
}

func (str StructMetadata) Match(other StructMetadata) bool {
	if len(str.Fields) != len(other.Fields) {
		return false
	}

	for i, field := range str.Fields {
		if !field.Match(other.Fields[i]) {
			return false
		}
	}

	return true
}

func (str StructMetadata) Definition(cfg Config) *j.Statement {
	return j.Type().Id(str.Name).StructFunc(func(group *j.Group) {
		for _, field := range str.Fields {
			group.Add(field.Definition(cfg))
		}
	}).Line()
}

func (str StructMetadata) Result(cfg Config) *j.Statement {
	switch len(str.Fields) {
	case 0:
		return j.Null()
	case 1:
		return str.Fields[0].Type.Statement()
	default:
		result := j.Id(str.Name)

		if cfg.EmitResultStructPointers {
			result = j.Op("*").Add(result)
		}

		return result
	}
}

func (str StructMetadata) Var(name string) *j.Statement {
	switch len(str.Fields) {
	case 0:
		return j.Null()
	case 1:
		return j.Var().Id(name).Add(str.Fields[0].Type.Statement())
	default:
		return j.Var().Id(name).Id(str.Name)
	}
}

func (field FieldMetadata) Match(other FieldMetadata) bool {
	if field.Name != other.Name {
		return false
	}

	if !field.Type.Match(other.Type) {
		return false
	}

	return true
}

func (field FieldMetadata) Definition(cfg Config) *j.Statement {
	result := j.Id(field.Name).Add(field.Type.Statement())

	if cfg.EmitDBTags && field.ColumnName != "" {
		result = result.Tag(map[string]string{
			"db": field.ColumnName,
		})
	}

	return result
}

func (t TypeMetadata) Match(other TypeMetadata) bool {
	if t.Path != other.Path {
		return false
	}

	if t.Name != other.Name {
		return false
	}

	if t.Array != other.Array {
		return false
	}

	if t.Ptr != other.Ptr {
		return false
	}

	return true
}

func (t TypeMetadata) Statement() *j.Statement {
	var result *j.Statement

	if t.Path == "" {
		result = j.Id(t.Name)
	} else {
		result = j.Qual(t.Path, t.Name)
	}

	if t.Ptr {
		result = j.Op("*").Add(result)
	}

	if t.Array {
		result = j.Index().Add(result)
	}

	return result
}
