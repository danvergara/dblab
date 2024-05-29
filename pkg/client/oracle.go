package client

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/sijms/go-ora/v2"
)

// oracle struct is in charge of perform all the oracle related queries.
type oracle struct{}

// a validation to see if oracle is implementing databaseQuerier.
var _ databaseQuerier = (*oracle)(nil)

// returns a pointer to a oracle struct, it receives an schema as a parameter.
func newOracle() *oracle {
	o := oracle{}

	return &o
}

// ShowTables returns a query to retrieve all the tables.
func (p *oracle) ShowTables() (string, []interface{}, error) {
	query := "SELECT table_name FROM user_tables"
	return query, nil, nil
}

// TableStructure returns a query string to get all the relevant information of a given table.
func (p *oracle) TableStructure(tableName string) (string, []interface{}, error) {
	query := sq.Select("*").
		From("user_tab_columns").
		Where(sq.Eq{"table_name": strings.ToUpper(tableName)}).
		OrderBy("column_id").
		PlaceholderFormat(sq.Colon)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, nil
}

// Constraints returns all the constraints of a given table.
func (p *oracle) Constraints(tableName string) (string, []interface{}, error) {
	query := sq.Select(
		`constraint_name`,
		`constraint_type`,
	).
		From("user_constraints").
		Where(sq.Eq{"table_name": tableName}).
		PlaceholderFormat(sq.Colon)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}

// Indexes returns the indexes of a table.
func (p *oracle) Indexes(tableName string) (string, []interface{}, error) {
	sql, args, err := sq.Select("*").
		From("all_indexes").
		Where(sq.Eq{"table_name": tableName}).
		PlaceholderFormat(sq.Colon).
		ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}
