package internal

import (
	j "github.com/dave/jennifer/jen"

	"github.com/MartyHub/sqlc-pg/plugin"
)

func (gen *Generator) enum(e *plugin.Enum) *j.Statement {
	name := gen.tok.ExportID(e.Name)

	return j.Const().
		DefsFunc(func(group *j.Group) {
			for _, value := range e.Vals {
				group.Id(name + gen.tok.ExportID(value)).Op("=").Lit(value)
			}
		}).
		Line().
		Type().Id(name).Id("string")
}
