package client

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// mssql struct is in charge of perform all the SQL Server related queries.
type mssql struct{}

var _ databaseQuerier = (*mssql)(nil)

// returns a pointer to a mysql.
func newMSSQL() *mssql {
	m := mssql{}
	return &m
}

func (m *mssql) ShowTables() (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.AtP)

	query, args, err := psql.Select("table_name").
		From("INFORMATION_SCHEMA.TABLES").
		Where(sq.Eq{"table_type": "BASE TABLE"}).
		ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

// TableStructure returns a query string to retrieve all the relevant information of a given table.
func (m *mssql) TableStructure(tableName string) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.AtP)
	query, args, err := psql.Select(
		"c.name AS ColumnName",
		"t.Name AS DataType",
		"c.max_length AS MaxLength",
		"c.precision",
		"c.scale",
		"c.is_nullable",
	).
		From("sys.columns c").
		InnerJoin("sys.types t ON c.user_type_id = t.user_type_id").
		Where(sq.Eq{"c.object_id": fmt.Sprintf("OBJECT_ID('%s')", tableName)}).
		ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

// Constraints returns all the constraints of a given table.
func (m *mssql) Constraints(tableName string) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.AtP)

	query, args, err := psql.Select(
		`constraint_name`,
		`table_name`,
		`constraint_type`,
	).
		From("INFORMATION_SCHEMA.TABLE_CONSTRAINTS").
		Where(sq.Eq{"TABLE_NAME": tableName}).
		ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

// Indexes returns the indexes of a table.
func (m *mssql) Indexes(tableName string) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.AtP)
	query, args, err := psql.Select(
		"ind.name AS IndexName",
		"ind.index_id AS IndexID",
		"col.name AS ColumnName",
		"ind.type_desc AS IndexType",
		"ind.is_unique AS IsUnique",
		"ind.is_primary_key AS IsPrimaryKey",
		"ind.is_unique_constraint AS IsUniqueConstraint",
		"ic.index_column_id AS ColumnID",
		"ic.is_descending_key AS IsDescendingKey",
		"ic.is_included_column AS IsIncludedColumn",
	).
		From("sys.indexes ind").
		InnerJoin(
			`sys.index_columns ic
          ON ind.object_id = ic.object_id
            AND ind.index_id = ic.index_id`,
		).
		InnerJoin(
			`sys.columns col
        ON ic.object_id = col.object_id
          AND ic.column_id = col.column_id`,
		).
		InnerJoin("sys.tables t ON ind.object_id = t.object_id").
		Where(sq.Eq{"t.name": tableName}).
		OrderBy(
			"ind.name",
			"ic.index_column_id",
		).
		ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}
