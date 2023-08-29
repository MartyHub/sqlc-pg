# sqlc-pg

`sqlc-pg` is a [sqlc](https://sqlc.dev) plugin for [PostgreSQL](https://www.postgresql.org)
using [pgx](https://github.com/jackc/pgx) Go driver.

![build](https://github.com/MartyHub/sqlc-pg/actions/workflows/go.yml/badge.svg)

## Generated Code Examples

- [authors](https://github.com/MartyHub/sqlc-pg/tree/main/internal/testdata/sqlc/examples/authors/sqlc)
- [booktest](https://github.com/MartyHub/sqlc-pg/tree/main/internal/testdata/sqlc/examples/booktest/sqlc)
- [jets](https://github.com/MartyHub/sqlc-pg/tree/main/internal/testdata/sqlc/examples/jets/sqlc)
- [ondeck](https://github.com/MartyHub/sqlc-pg/tree/main/internal/testdata/sqlc/examples/ondeck/sqlc)

## sqlc Sample Configuration

- `sqlc-pg` must be available in your path
- Sample `sqlc.yaml`:
    ```yaml
    version: "2"
    plugins:
      - name: sqlc-pg
        process:
          cmd: sqlc-pg
    sql:
      - engine: postgresql
        queries: query.sql
        schema: schema.sql
        codegen:
          - plugin: sqlc-pg
            out: sqlc
            options:
              emit_db_tags: true
              emit_exported_queries: false
              emit_params_struct_pointers: false
              emit_result_struct_pointers: false
              emit_table_names: true
              output_files_suffix: .gen
    ```
