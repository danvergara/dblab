package client

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// mssql struct is in charge of perform all the SQL Server related queries.
type mssql struct {
	db     *sqlx.DB
	dbName string
}

var _ databaseQuerier = (*mssql)(nil)

// returns a pointer to a mysql.
func newMSSQL(dbName string, db *sqlx.DB) *mssql {
	m := mssql{
		dbName: dbName,
		db:     db,
	}
	return &m
}

// TableStructure returns a query string to retrieve all the relevant information of a given table.
func (m *mssql) TableStructure(table TableRef) (string, []interface{}, error) {
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
		Where(sq.Eq{"c.object_id": fmt.Sprintf("OBJECT_ID('%s')", table.Name)}).
		ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

// Constraints returns all the constraints of a given table.
func (m *mssql) Constraints(table TableRef) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.AtP)

	query, args, err := psql.Select(
		`constraint_name`,
		`table_name`,
		`constraint_type`,
	).
		From("INFORMATION_SCHEMA.TABLE_CONSTRAINTS").
		Where(sq.Eq{"TABLE_NAME": table.Name}).
		ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

// Indexes returns the indexes of a table.
func (m *mssql) Indexes(table TableRef) (string, []interface{}, error) {
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
		Where(sq.Eq{"t.name": table.Name}).
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

func (m *mssql) Catalog(ctx context.Context) (*DBNode, error) {
	rootID := fmt.Sprintf("db:%s", m.dbName)
	root := &DBNode{ID: rootID, Name: m.dbName, Type: "database"}

	queue := []*DBNode{root}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		var children []*DBNode
		var err error
		switch current.Type {
		case "database":
			children, err = m.fetchTables(ctx, current.Name, current.ID)
		}
		if err != nil {
			return nil, err
		}

		for _, child := range children {
			current.Children = append(current.Children, child)
			queue = append(queue, child)
		}
	}

	return root, nil
}

func (m *mssql) fetchTables(ctx context.Context, parentName, parentID string) ([]*DBNode, error) {
	query := "SHOW TABLES;"

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []*DBNode
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		tables = append(tables, &DBNode{
			ID:         fmt.Sprintf("%s.t:%s", parentID, name),
			Name:       name,
			Type:       "table",
			ParentName: parentName,
			ParentID:   parentID,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}
