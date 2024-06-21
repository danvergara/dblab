package client

import (
	sq "github.com/Masterminds/squirrel"
)

// postgres struct is in charge of perform all the postgres related queries,
// without the client knowing.
type postgres struct {
	schema string
}

// a validation to see if postgres is implementing databaseQuerier.
var _ databaseQuerier = (*postgres)(nil)

// returns a pointer to a postgres, it receives an schema as a parameter.
func newPostgres(schema string) *postgres {
	p := postgres{
		schema: schema,
	}

	return &p
}

// ShowTables returns a query to retrieve all the tables under a specific schema.
func (p *postgres) ShowTables() (string, []interface{}, error) {
	var (
		query string
		err   error
		args  []interface{}
	)
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err = psql.Select("table_name").
		From("information_schema.tables").
		Where(sq.Eq{"table_schema": p.schema}).
		OrderBy("table_name").
		ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

// TableStructure returns a query string to get all the relevant information of a given table,
// under a schema.
func (p *postgres) TableStructure(tableName string) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.Select(
		"c.column_name",
		"c.is_nullable",
		"c.data_type",
		"c.character_maximum_length",
		"c.numeric_precision",
		"c.numeric_scale",
		"c.ordinal_position",
		"tc.constraint_type AS pkey",
	).
		From("information_schema.columns AS c").
		LeftJoin(
			`information_schema.constraint_column_usage AS ccu
					ON c.table_schema = ccu.table_schema
						AND c.table_name = ccu.table_name
						AND c.column_name = ccu.column_name`,
		).
		LeftJoin(
			`information_schema.table_constraints AS tc
					ON ccu.constraint_schema = tc.constraint_schema
						AND ccu.constraint_name = tc.constraint_name`,
		).
		Where(
			sq.And{
				sq.Eq{"c.table_schema": p.schema},
				sq.Eq{"c.table_name": tableName},
			},
		).
		ToSql()

	return query, args, err
}

// Constraints returns all the constraints of a given table, under a schema.
func (p *postgres) Constraints(tableName string) (string, []interface{}, error) {
	var (
		query sq.SelectBuilder
		sql   string
	)

	query = sq.Select(
		`tc.constraint_name`,
		`tc.table_name`,
		`tc.constraint_type`,
	).
		From("information_schema.table_constraints AS tc").
		Where(sq.Eq{"tc.table_name": tableName}).
		Where(sq.Eq{"tc.table_schema": p.schema}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, err
}

// Indexes returns the indexes of a table, under a schema.
func (p *postgres) Indexes(tableName string) (string, []interface{}, error) {
	query := sq.Select("*").
		From("pg_indexes").
		Where(sq.Eq{"tableName": tableName}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}
