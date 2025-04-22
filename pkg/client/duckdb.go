package client

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type duckdb struct{}

func newDuckDB() *duckdb {
	d := duckdb{}
	return &d
}

var _ databaseQuerier = (*duckdb)(nil)

func (d *duckdb) ShowTablesPerDB(dabase string) (string, []interface{}, error) {
	return "", nil, nil
}

func (d *duckdb) ShowDatabases() (string, []interface{}, error) {
	return "", nil, nil
}

// ShowTables returns a query to retrieve all the tables.
func (d *duckdb) ShowTables() (string, []interface{}, error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'main'`

	return query, nil, nil
}

func (d *duckdb) TableStructure(tableName string) (string, []interface{}, error) {
	query := fmt.Sprintf("DESCRIBE %s;", tableName)
	return query, nil, nil
}

func (d *duckdb) Constraints(tableName string) (string, []interface{}, error) {
	query := sq.Select("*").
		From("information_schema.table_constraints").
		Where(sq.Eq{"tableName": tableName})

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}

func (d *duckdb) Indexes(tableName string) (string, []interface{}, error) {
	query := fmt.Sprintf(`PRAGMA index_list(%s);`, tableName)

	return query, nil, nil
}
