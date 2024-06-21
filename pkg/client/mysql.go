package client

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// mysql struct is in charge of perform all the mysql related queries.
type mysql struct{}

// a validation to see if mysql is implementing databaseQuerier.
var _ databaseQuerier = (*mysql)(nil)

// returns a pointer to a mysql.
func newMySQL() *mysql {
	m := mysql{}
	return &m
}

// ShowTables returns a query to retrieve all the tables.
func (m *mysql) ShowTables() (string, []interface{}, error) {
	query := "SHOW TABLES;"
	return query, nil, nil
}

// TableStructure returns a query string to retrieve all the relevant information of a given table.
func (m *mysql) TableStructure(tableName string) (string, []interface{}, error) {
	query := fmt.Sprintf("DESCRIBE %s;", tableName)
	return query, nil, nil
}

// Constraints returns all the constraints of a given table.
func (m *mysql) Constraints(tableName string) (string, []interface{}, error) {
	query := sq.Select(
		`tc.constraint_name`,
		`tc.table_name`,
		`tc.constraint_type`,
	).
		From("information_schema.table_constraints AS tc").
		Where("tc.table_name = ?", tableName)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}

// Indexes returns a query to get all the indexes of a table.
func (m *mysql) Indexes(tableName string) (string, []interface{}, error) {
	query := fmt.Sprintf("SHOW INDEX FROM %s", tableName)
	return query, nil, nil
}
