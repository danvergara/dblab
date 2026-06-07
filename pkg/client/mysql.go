package client

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// mysql struct is in charge of perform all the mysql related queries.
type mysql struct {
	db     *sqlx.DB
	dbName string
}

// a validation to see if mysql is implementing databaseQuerier.
var _ databaseQuerier = (*mysql)(nil)

// returns a pointer to a mysql.
func newMySQL(dbName string, db *sqlx.DB) *mysql {
	m := mysql{
		dbName: dbName,
		db:     db,
	}

	return &m
}

// TableStructure returns a query string to retrieve all the relevant information of a given table.
func (m *mysql) TableStructure(table TableRef) (string, []any, error) {
	query := fmt.Sprintf("DESCRIBE %s;", table.Name)
	return query, nil, nil
}

// Constraints returns all the constraints of a given table.
func (m *mysql) Constraints(table TableRef) (string, []any, error) {
	query := sq.Select(
		`tc.constraint_name`,
		`tc.table_name`,
		`tc.constraint_type`,
	).
		From("information_schema.table_constraints AS tc").
		Where("tc.table_name = ?", table.Name)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}

// Indexes returns a query to get all the indexes of a table.
func (m *mysql) Indexes(table TableRef) (string, []any, error) {
	query := fmt.Sprintf("SHOW INDEX FROM %s", table.Name)
	return query, nil, nil
}

func (m *mysql) Catalog(ctx context.Context) (*DBNode, error) {
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
			tables, err := m.fetchTables(ctx, current.Name, current.ID)
			if err != nil {
				return nil, err
			}
			children = append(children, tables...)

			views, err := m.fetchViews(ctx, current.Name, current.ID)
			if err != nil {
				return nil, err
			}
			children = append(children, views...)
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

func (m *mysql) GetViewDefinition(view ViewRef) (string, []any, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Question)

	query, args, err := psql.
		Select("VIEW_DEFINITION AS view_definition").
		From("information_schema.VIEWS").
		Where(sq.Eq{
			"TABLE_SCHEMA": m.dbName,
			"TABLE_NAME":   view.Name,
		}).
		ToSql()

	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

func (m *mysql) fetchTables(ctx context.Context, parentName, parentID string) ([]*DBNode, error) {
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
			Name:       name + " - " + "t",
			EntityName: name,
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

func (m *mysql) fetchViews(ctx context.Context, parentName, parentID string) ([]*DBNode, error) {
	query, args, err := sq.
		Select("TABLE_NAME AS view_name").
		From("information_schema.VIEWS").
		Where(sq.Eq{
			"TABLE_SCHEMA": m.dbName,
		}).
		PlaceholderFormat(sq.Question).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := m.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	views := make([]*DBNode, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		views = append(views, &DBNode{
			ID:         fmt.Sprintf("%s.v:%s", parentName, name),
			Name:       name + " - " + "v",
			EntityName: name,
			Type:       "view",
			ParentName: parentName,
			ParentID:   parentID,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return views, nil
}
