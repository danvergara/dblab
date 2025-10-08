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

func (p *oracle) ShowTablesPerDB(dabase string) (string, []interface{}, error) {
	return "", nil, nil
}

func (p *oracle) ShowDatabases() (string, []interface{}, error) {
	return "", nil, nil
}

// ShowTables returns a query to retrieve all the tables.
func (p *oracle) ShowTables() (string, []interface{}, error) {
	query := "SELECT OWNER || '.' || TABLE_NAME FROM ALL_TABLES"
	return query, nil, nil
}

// TableStructure returns a query string to get all the relevant information of a given table.
func (p *oracle) TableStructure(tableName string) (string, []interface{}, error) {
	owner, table := splitOracleTableName(tableName)

	query := sq.Select("*").
		From("ALL_TAB_COLUMNS").
		Where(sq.Eq{"TABLE_NAME": table}).
		OrderBy("COLUMN_ID").
		PlaceholderFormat(sq.Colon)

	if owner != "" {
		query = query.Where(sq.Eq{"OWNER": owner})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, nil
}

// Constraints returns all the constraints of a given table.
func (p *oracle) Constraints(tableName string) (string, []interface{}, error) {
	owner, table := splitOracleTableName(tableName)

	query := sq.Select(
		`CONSTRAINT_NAME`,
		`CONSTRAINT_TYPE`,
	).
		From("ALL_CONSTRAINTS").
		Where(sq.Eq{"TABLE_NAME": table}).
		PlaceholderFormat(sq.Colon)

	if owner != "" {
		query = query.Where(sq.Eq{"OWNER": owner})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}

// Indexes returns the indexes of a table.
func (p *oracle) Indexes(tableName string) (string, []interface{}, error) {
	owner, table := splitOracleTableName(tableName)

	query := sq.Select("*").
		From("ALL_INDEXES").
		Where(sq.Eq{"TABLE_NAME": table}).
		PlaceholderFormat(sq.Colon)

	if owner != "" {
		query = query.Where(sq.Eq{"OWNER": owner})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}

// splitOracleTableName splits a table name into owner (schema) and table.
func splitOracleTableName(fullName string) (string, string) {
	parts := strings.Split(fullName, ".")
	if len(parts) == 2 {
		return strings.ToUpper(parts[0]), strings.ToUpper(parts[1])
	}
	return "", strings.ToUpper(fullName)
}
