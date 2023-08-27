package internal

import (
	j "github.com/dave/jennifer/jen"
)

type RepositoryMetadata struct {
	Name    string
	Queries []QueryMetadata
}

func (gen *Generator) repository(repo *RepositoryMetadata) *j.Statement {
	db := j.Id("db")
	dbParam := j.Id("db").Id("Database")
	exportedName := gen.tok.ExportID(repo.Name)
	unexportedName := gen.tok.UnexportID(repo.Name)

	return j.Type().Id(exportedName).
		InterfaceFunc(func(group *j.Group) {
			for _, query := range repo.Queries {
				group.Add(gen.queryFuncSig(query.Query, query.Params, query.Row))
			}
		}).
		Line().
		Line().
		Type().Id(unexportedName).Struct(dbParam).
		Line().
		Func().
		Id("New" + exportedName).
		Params(dbParam).
		Id(exportedName).
		Block(
			j.Return(
				j.Id(unexportedName).Values(j.Dict{db: db}),
			),
		)
}
