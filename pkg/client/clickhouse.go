package client

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// clickhouse struct is in charge of perform all the clickhouse related queries.
type clickhouse struct {
	db     *sqlx.DB
	dbName string
}

// a validation to see if clickhouse is implementing databaseQuerier.
var _ databaseQuerier = (*clickhouse)(nil)

// returns a pointer to a clickhouse.
func newClickHouse(dbName string, db *sqlx.DB) *clickhouse {
	m := clickhouse{
		dbName: dbName,
		db:     db,
	}

	return &m
}

// TableStructure returns a query string to retrieve all the relevant information of a given table.
func (m *clickhouse) TableStructure(table TableRef) (string, []any, error) {
	query := fmt.Sprintf("DESCRIBE TABLE %s.%s;", m.dbName, table.Name)
	return query, nil, nil
}

// Constraints returns all the constraints of a given table.
func (m *clickhouse) Constraints(table TableRef) (string, []any, error) {
	query := fmt.Sprintf("SELECT name, type, expression FROM system.constraints WHERE database = '%s' AND table = '%s';", m.dbName, table.Name)
	return query, nil, nil
}

// Indexes returns a query to get all the indexes of a table.
func (m *clickhouse) Indexes(table TableRef) (string, []any, error) {
	query := fmt.Sprintf("SHOW INDEXES FROM %s.%s;", m.dbName, table.Name)
	return query, nil, nil
}

// Catalog returns a the pointer to a DBNode instance,
// which is the root of the current ClickHouse database graph.
// It starts with the database itself and a list of tables a views.
// ClickHouse topography:
//
//			     [Database]
//			      /     \
//			     v       v
//	 			[Tables] 	[Views]
func (m *clickhouse) Catalog(ctx context.Context) (*DBNode, error) {
	rootID := fmt.Sprintf("db:%s", m.dbName)
	root := &DBNode{ID: rootID, Name: m.dbName, Type: "database"}

	queue := []*DBNode{root}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		var children []*DBNode

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

		for _, child := range children {
			current.Children = append(current.Children, child)
			queue = append(queue, child)
		}
	}

	return root, nil
}

// GetViewDefinition method returns the SQL definition of a given view.
func (m *clickhouse) GetViewDefinition(view ViewRef) (string, []any, error) {
	query := fmt.Sprintf("SHOW CREATE VIEW %s.%s;", m.dbName, view.Name)
	return query, nil, nil
}

// fetchTables method lists all the tables of the current database.
func (m *clickhouse) fetchTables(_ context.Context, parentName, parentID string) ([]*DBNode, error) {
	query := fmt.Sprintf("SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = '%s' AND TABLE_TYPE != 'VIEW';", m.dbName)

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

// fetchViews method lists all the views of the current database.
func (m *clickhouse) fetchViews(_ context.Context, parentName, parentID string) ([]*DBNode, error) {
	query := fmt.Sprintf("SELECT table_name FROM INFORMATION_SCHEMA.VIEWS WHERE table_schema = '%s';", m.dbName)

	rows, err := m.db.Query(query)
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
