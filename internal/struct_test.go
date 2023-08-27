package internal

import (
	"testing"
)

func Test_struct(t *testing.T) {
	t.Parallel()

	str := StructMetadata{
		Name: "MyStruct",
		Fields: []FieldMetadata{
			{
				Name: "Bool",
				Type: TypeMetadata{Name: "bool"},
			},
			{
				Name: "Integers",
				Type: TypeMetadata{Name: "int", Array: true},
			},
			{
				Name: "String",
				Type: TypeMetadata{Name: "string", Ptr: true},
			},
			{
				Name: "Times",
				Type: TypeMetadata{Name: "Time", Array: true, Ptr: true, Path: "time"},
			},
		},
	}

	compare(t, "testdata/struct/nominal", str.Definition(Config{}))
}
