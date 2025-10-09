package client

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/sijms/go-ora/v2"
)

// oracle struct is in charge of perform all the oracle related queries.
type oracle struct {
	schema string
}

// a validation to see if oracle is implementing databaseQuerier.
var _ databaseQuerier = (*oracle)(nil)

// returns a pointer to a oracle struct, it receives an schema as a parameter.
func newOracle(schema string) *oracle {
	// If the schema is no empty, the client queries against the ALL_* tables,
	// where the OWNER is equal to the schema/user the dblab user has access to.
	// Otherwise, the client will query against the USER_* tables,
	// meaning it only cares about what the user has direct access to.
	o := oracle{schema: schema}

	return &o
}

func (o *oracle) ShowTablesPerDB(dabase string) (string, []interface{}, error) {
	return "", nil, nil
}

func (o *oracle) ShowDatabases() (string, []interface{}, error) {
	return "", nil, nil
}

// ShowTables returns a query to retrieve all the tables.
func (o *oracle) ShowTables() (string, []interface{}, error) {
	var query sq.SelectBuilder

	if o.schema != "" {
		query = sq.Select("TABLE_NAME").
			From("ALL_TABLES").
			Where(sq.Eq{"OWNER": strings.ToUpper(o.schema)})
	} else {
		query = sq.Select("TABLE_NAME").
			From("USER_TABLES")
	}

	sql, args, err := query.OrderBy("1").PlaceholderFormat(sq.Colon).ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, nil
}

// TableStructure returns a query string to get all the relevant information of a given table.
func (o *oracle) TableStructure(tableName string) (string, []interface{}, error) {
	var query sq.SelectBuilder

	if o.schema != "" {
		query = sq.Select("*").
			From("ALL_TAB_COLUMNS").
			Where(sq.Eq{"TABLE_NAME": strings.ToUpper(tableName), "OWNER": strings.ToUpper(o.schema)})
	} else {
		query = sq.Select("*").
			From("USER_TAB_COLUMNS").
			Where(sq.Eq{"TABLE_NAME": strings.ToUpper(tableName)})
	}

	sql, args, err := query.OrderBy("1").PlaceholderFormat(sq.Colon).ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, nil
}

// Constraints returns all the constraints of a given table.
func (o *oracle) Constraints(tableName string) (string, []interface{}, error) {
	var query sq.SelectBuilder

	if o.schema != "" {
		query = sq.Select(
			`CONSTRAINT_NAME`,
			`CONSTRAINT_TYPE`,
		).
			From("ALL_CONSTRAINTS").
			Where(sq.Eq{"TABLE_NAME": strings.ToUpper(tableName), "OWNER": strings.ToUpper(o.schema)})
	} else {
		query = sq.Select(
			`CONSTRAINT_NAME`,
			`CONSTRAINT_TYPE`,
		).
			From("USER_CONSTRAINTS").
			Where(sq.Eq{"TABLE_NAME": strings.ToUpper(tableName)})
	}

	sql, args, err := query.PlaceholderFormat(sq.Colon).ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}

// Indexes returns the indexes of a table.
func (o *oracle) Indexes(tableName string) (string, []interface{}, error) {
	var query sq.SelectBuilder

	if o.schema != "" {
		query = sq.Select("*").
			From("ALL_INDEXES").
			Where(sq.Eq{"TABLE_NAME": strings.ToUpper(tableName), "OWNER": strings.ToUpper(o.schema)})
	} else {
		query = sq.Select("*").
			From("USER_INDEXES").
			Where(sq.Eq{"TABLE_NAME": strings.ToUpper(tableName)})
	}

	sql, args, err := query.PlaceholderFormat(sq.Colon).ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, nil
}
