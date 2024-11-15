package client

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// sqlite  is in charge of perform all the sqlite related queries,
// without the client knowing.
type sqlite struct{}

// a validation to see if sqlite is implementing databaseQuerier.
var _ databaseQuerier = (*sqlite)(nil)

// returns a pointer to a sqlite.
func newSQLite() *sqlite {
	s := sqlite{}
	return &s
}

func (s *sqlite) ShowTablesPerDB(dabase string) (string, []interface{}, error) {
	return "", nil, nil
}

func (s *sqlite) ShowDatabases() (string, []interface{}, error) {
	return "", nil, nil
}

// ShowTables returns a query to retrieve all the tables.
func (s *sqlite) ShowTables() (string, []interface{}, error) {
	query := `
		SELECT
			name
		FROM
			sqlite_schema
		WHERE
			type ='table' AND
			name NOT LIKE 'sqlite_%';`

	return query, nil, nil
}

// TableStructure returns a query string to retrieve all the relevant information of a given table.
func (s *sqlite) TableStructure(tableName string) (string, []interface{}, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s);", tableName)
	return query, nil, nil
}

// Constraints returns all the constraints of a given table.
func (s *sqlite) Constraints(tableName string) (string, []interface{}, error) {
	query := sq.Select(
		"*",
	).
		From("sqlite_master").
		Where(
			sq.And{
				sq.Eq{"type": "table"},
				sq.Eq{"name": tableName},
			})

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, nil
}

// Indexes returns a query to get all the indexes of a table.
func (s *sqlite) Indexes(tableName string) (string, []interface{}, error) {
	query := fmt.Sprintf(`PRAGMA index_list(%s);`, tableName)

	return query, nil, nil
}
