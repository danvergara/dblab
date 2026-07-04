package client

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// sqlite  is in charge of perform all the sqlite related queries,
// without the client knowing.
type sqlite struct {
	db     *sqlx.DB
	dbName string
}

// a validation to see if sqlite is implementing databaseQuerier.
var _ databaseQuerier = (*sqlite)(nil)

// returns a pointer to a sqlite.
func newSQLite(dbName string, db *sqlx.DB) *sqlite {
	s := sqlite{
		dbName: dbName,
		db:     db,
	}

	return &s
}

// TableStructure returns a query string to retrieve all the relevant information of a given table.
func (s *sqlite) TableStructure(table TableRef) (string, []any, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s);", table.Name)
	return query, nil, nil
}

// Constraints returns all the constraints of a given table.
func (s *sqlite) Constraints(table TableRef) (string, []any, error) {
	query := sq.Select(
		"*",
	).
		From("sqlite_master").
		Where(
			sq.And{
				sq.Eq{"type": "table"},
				sq.Eq{"name": table.Name},
			})

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, nil
}

// Indexes returns a query to get all the indexes of a table.
func (s *sqlite) Indexes(table TableRef) (string, []any, error) {
	query := fmt.Sprintf(`PRAGMA index_list(%s);`, table.Name)

	return query, nil, nil
}

// Catalog returns a the pointer to a DBNode instance,
// which is the root of the current SQLite database graph.
// It starts with the database itself and a list of tables a views.
// SQLite topography:
//
//			     [Database]
//			      /     \
//			     v       v
//	 			[Tables] 	[Views]
func (s *sqlite) Catalog(ctx context.Context) (*DBNode, error) {
	rootID := fmt.Sprintf("db:%s", s.dbName)
	root := &DBNode{ID: rootID, Name: s.dbName, Type: "database"}

	queue := []*DBNode{root}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		var children []*DBNode
		var err error
		switch current.Type {
		case "database":
			tables, err := s.fetchTables(ctx, current.Name, current.ID)
			if err != nil {
				return nil, err
			}
			children = append(children, tables...)

			views, err := s.fetchViews(ctx, current.Name, current.ID)
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

// GetViewDefinition method returns the SQL definition of a given view.
func (s *sqlite) GetViewDefinition(view ViewRef) (string, []any, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Question)
	query, args, err := psql.
		Select("sql AS view_definition").
		From("sqlite_master").
		Where(sq.Eq{
			"type": "view",
			"name": view.Name,
		}).
		ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

// fetchTables method lists all the tables of the current database.
func (s *sqlite) fetchTables(ctx context.Context, parentName, parentID string) ([]*DBNode, error) {
	query := `
		SELECT
			name
		FROM
			sqlite_schema
		WHERE
			type ='table' AND
			name NOT LIKE 'sqlite_%';`

	rows, err := s.db.Query(query)
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
func (s *sqlite) fetchViews(ctx context.Context, parentName, parentID string) ([]*DBNode, error) {
	query, args, err := sq.
		Select("name AS view_name").
		From("sqlite_master").
		Where(sq.Eq{
			"type": "view",
		}).
		OrderBy("name").
		PlaceholderFormat(sq.Question).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(query, args...)
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
			ID:         fmt.Sprintf("%s.v:%s", parentID, name),
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
