package internal

import (
	"strconv"

	j "github.com/dave/jennifer/jen"

	"github.com/MartyHub/sqlc-pg/plugin"
)

const (
	queryExec = ":exec"
	queryOne  = ":one"
	queryMany = ":many"
)

type QueryMetadata struct {
	Query *plugin.Query
	Name  string

	Params           StructMetadata
	ParamsToGenerate bool

	Row           StructMetadata
	RowToGenerate bool
}

func (gen *Generator) queryFuncSig(query *plugin.Query, params StructMetadata, row StructMetadata) *j.Statement {
	switch query.Cmd {
	case queryExec:
		return gen.queryExecFuncSig(query, params)
	case queryOne:
		return gen.queryOneFuncSig(query, params, row)
	case queryMany:
		return gen.queryManyFuncSig(query, params, row)
	}

	return j.Null()
}

func (gen *Generator) queryExecFuncSig(query *plugin.Query, params StructMetadata) *j.Statement {
	return j.Id(gen.tok.ExportID(query.Name)).
		Params(
			j.Id("ctx").Qual("context", "Context"),
			gen.params(params),
		).
		Parens(
			j.Id("error"),
		)
}

func (gen *Generator) queryOneFuncSig(query *plugin.Query, params StructMetadata, row StructMetadata) *j.Statement {
	return j.Id(gen.tok.ExportID(query.Name)).
		Params(
			j.Id("ctx").Qual("context", "Context"),
			gen.params(params),
		).
		Parens(
			j.List(
				row.Result(gen.cfg),
				j.Id("error"),
			),
		)
}

func (gen *Generator) queryManyFuncSig(query *plugin.Query, params StructMetadata, row StructMetadata) *j.Statement {
	return j.Id(gen.tok.ExportID(query.Name)).
		Params(
			j.Id("ctx").Qual("context", "Context"),
			gen.params(params),
		).
		Parens(
			j.List(
				j.Index().Add(row.Result(gen.cfg)),
				j.Id("error"),
			),
		)
}

func (gen *Generator) queryFunc(repo *RepositoryMetadata, query QueryMetadata) *j.Statement {
	switch query.Query.Cmd {
	case queryExec:
		return gen.queryExecFunc(repo, query)
	case queryOne:
		return gen.queryOneFunc(repo, query)
	case queryMany:
		return gen.queryManyFunc(repo, query)
	}

	return j.Null()
}

func (gen *Generator) queryExecFunc(repo *RepositoryMetadata, query QueryMetadata) *j.Statement {
	return j.Func().
		Params(j.Id("repo").Id(gen.tok.UnexportID(repo.Name))).
		Id(gen.tok.ExportID(query.Name)).
		Params(
			j.Id("ctx").Qual("context", "Context"),
			gen.params(query.Params),
		).
		Parens(
			j.Id("error"),
		).
		Block(
			j.List(
				j.Id("_"),
				j.Id("err"),
			).Op(":=").
				Id("repo").Dot("db").Dot("Exec").Call(
				j.Id("ctx"),
				j.Id(gen.stmtName(query.Query)),
				gen.call("params", query.Params),
			),
			j.Line().Return(
				j.Id("err"),
			),
		)
}

func (gen *Generator) queryOneFunc(repo *RepositoryMetadata, query QueryMetadata) *j.Statement {
	return j.Func().
		Params(j.Id("repo").Id(gen.tok.UnexportID(repo.Name))).
		Id(gen.tok.ExportID(query.Name)).
		Params(
			j.Id("ctx").Qual("context", "Context"),
			gen.params(query.Params),
		).
		Parens(
			j.List(
				query.Row.Result(gen.cfg),
				j.Id("error"),
			),
		).
		Block(
			j.List(j.Id("rows"), j.Id("err")).Op(":=").
				Id("repo").Dot("db").Dot("Query").Call(
				j.Id("ctx"),
				j.Id(gen.stmtName(query.Query)),
				gen.call("params", query.Params),
			),
			j.If(j.Id("err").Op("!=").Nil()).
				BlockFunc(func(group *j.Group) {
					if gen.cfg.EmitResultStructPointers && len(query.Row.Fields) > 1 {
						group.Return(j.Nil(), j.Id("err"))
					} else {
						group.Add(query.Row.Var("result")).
							Line().
							Line().
							Return(
								j.Id("result"),
								j.Id("err"),
							)
					}
				}),
			j.Line().Return(
				j.Id("CollectExactlyOneRow").Call(
					j.Id("rows"),
					j.Id(gen.scanName(query.Row)),
				),
			),
		)
}

func (gen *Generator) queryManyFunc(repo *RepositoryMetadata, query QueryMetadata) *j.Statement {
	return j.Func().
		Params(j.Id("repo").Id(gen.tok.UnexportID(repo.Name))).
		Id(gen.tok.ExportID(query.Name)).
		Params(
			j.Id("ctx").Qual("context", "Context"),
			gen.params(query.Params),
		).
		Parens(
			j.List(
				j.Index().Add(query.Row.Result(gen.cfg)),
				j.Id("error"),
			),
		).
		Block(
			j.List(j.Id("rows"), j.Id("err")).Op(":=").
				Id("repo").Dot("db").Dot("Query").Call(
				j.Id("ctx"),
				j.Id(gen.stmtName(query.Query)),
				gen.call("params", query.Params),
			),
			j.If(j.Id("err").Op("!=").Nil()).
				Block(
					j.Return(
						j.Nil(),
						j.Id("err"),
					),
				),
			j.Line().Return(
				j.Id("pgx.CollectRows").Call(
					j.Id("rows"),
					j.Id(gen.scanName(query.Row)),
				),
			),
		)
}

func (gen *Generator) queryParamStruct(query *plugin.Query) StructMetadata {
	result := StructMetadata{Name: gen.tok.ExportID(query.Name + "_Params")}

	for i, param := range query.Params {
		name := param.Column.Name

		if name == "" {
			name = "params" + strconv.Itoa(i+1)
		}

		result.Fields = append(result.Fields, FieldMetadata{
			Name:       gen.tok.ExportID(name),
			ColumnName: name,
			Type:       gen.mapType(param.Column),
		})
	}

	return result
}

func (gen *Generator) queryRowStruct(query *plugin.Query) StructMetadata {
	result := StructMetadata{Name: gen.tok.ExportID(query.Name)}

	for _, column := range query.Columns {
		result.Fields = append(result.Fields, FieldMetadata{
			Name:       gen.tok.ExportID(column.Name),
			ColumnName: column.Name,
			Type:       gen.mapType(column),
		})
	}

	return result
}

func (gen *Generator) queryScanFunc(rowStruct StructMetadata) *j.Statement {
	if len(rowStruct.Fields) == 1 {
		return gen.queryScanFieldFunc(rowStruct)
	}

	return gen.queryScanStructFunc(rowStruct)
}

func (gen *Generator) queryScanFieldFunc(rowStruct StructMetadata) *j.Statement {
	field := rowStruct.Fields[0]

	return j.Func().
		Id(gen.scanName(rowStruct)).
		Params(
			j.Id("row").Qual(pgxImport, "CollectableRow"),
		).
		Parens(
			j.List(field.Type.Statement(), j.Id("error")),
		).
		Block(
			j.Var().Id("result").Add(field.Type.Statement()).Line(),
			j.Id("err").Op(":=").
				Id("row").Dot("Scan").Call(j.Op("&").Id("result")),
			j.Line().
				Return(
					j.Id("result"),
					j.Id("err"),
				),
		).Line()
}

func (gen *Generator) queryScanStructFunc(rowStruct StructMetadata) *j.Statement {
	return j.Func().
		Id(gen.scanName(rowStruct)).
		Params(
			j.Id("row").Qual(pgxImport, "CollectableRow"),
		).
		Parens(
			j.List(
				j.Do(func(stmt *j.Statement) {
					if gen.cfg.EmitResultStructPointers {
						stmt.Op("*")
					}

					stmt.Id(rowStruct.Name)
				}),
				j.Id("error"),
			),
		).
		Block(
			j.Do(func(stmt *j.Statement) {
				if gen.cfg.EmitResultStructPointers {
					stmt.Id("result").Op(":=").New(j.Id(rowStruct.Name)).Line()
				} else {
					stmt.Var().Id("result").Id(rowStruct.Name).Line()
				}
			}),
			j.Do(func(stmt *j.Statement) {
				assign := j.Id("err").Op(":=").Id("row").Dot("Scan").CallFunc(func(group *j.Group) {
					for _, field := range rowStruct.Fields {
						group.Line().Op("&").Id("result").Dot(field.Name)
					}
					group.Line()
				})

				if gen.cfg.EmitResultStructPointers {
					stmt.If(assign, j.Id("err").Op("!=").Nil()).
						Block(
							j.Return(
								j.Nil(),
								j.Id("err"),
							),
						)
				} else {
					stmt.Add(assign)
				}
			}),
			j.Do(func(stmt *j.Statement) {
				if gen.cfg.EmitResultStructPointers {
					stmt.Line().Return(
						j.Id("result"),
						j.Nil(),
					)
				} else {
					stmt.Line().Return(
						j.Id("result"),
						j.Id("err"),
					)
				}
			}),
		).Line()
}

func (gen *Generator) queryStmt(query *plugin.Query) *j.Statement {
	opts := j.Options{
		Open:  "`",
		Close: "`",
	}

	return j.Const().Id(gen.stmtName(query)).Op("=").Custom(opts, j.Id(query.Text)).Line()
}

func (gen *Generator) call(name string, str StructMetadata) *j.Statement {
	switch len(str.Fields) {
	case 0:
		return j.Null()
	case 1:
		return j.Id(gen.tok.UnexportID(str.Fields[0].ColumnName))
	default:
		return j.ListFunc(func(group *j.Group) {
			for _, field := range str.Fields {
				group.Line().Id(name).Dot(field.Name)
			}

			group.Line()
		})
	}
}

func (gen *Generator) params(str StructMetadata) *j.Statement {
	switch len(str.Fields) {
	case 0:
		return j.Null()
	case 1:
		return j.Id(gen.tok.UnexportID(str.Fields[0].ColumnName)).Add(str.Fields[0].Type.Statement())
	default:
		result := j.Id("params")

		if gen.cfg.EmitParamsStructPointers {
			result = result.Op("*")
		}

		result = result.Id(str.Name)

		return result
	}
}
