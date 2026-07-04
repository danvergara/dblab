package client

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/sijms/go-ora/v2"
)

// oracle struct is in charge of perform all the oracle related queries.
type oracle struct {
	db     *sqlx.DB
	dbName string
	schema string
}

// a validation to see if oracle is implementing databaseQuerier.
var _ databaseQuerier = (*oracle)(nil)

// returns a pointer to a oracle struct, it receives an schema as a parameter.
func newOracle(dbName, schema string, db *sqlx.DB) *oracle {
	// If the schema is no empty, the client queries against the ALL_* tables,
	// where the OWNER is equal to the schema/user the dblab user has access to.
	// Otherwise, the client will query against the USER_* tables,
	// meaning it only cares about what the user has direct access to.
	o := oracle{
		dbName: dbName,
		db:     db,
		schema: schema,
	}

	return &o
}

// TableStructure returns a query string to get all the relevant information of a given table.
func (o *oracle) TableStructure(table TableRef) (string, []any, error) {
	var query sq.SelectBuilder

	query = sq.Select("*").
		From("USER_TAB_COLUMNS").
		Where(sq.Eq{"TABLE_NAME": strings.ToUpper(table.Name)})

	if table.Schema != "" {
		query = sq.Select("*").
			From("ALL_TAB_COLUMNS").
			Where(sq.Eq{"TABLE_NAME": strings.ToUpper(table.Name), "OWNER": strings.ToUpper(table.Schema)})
	}

	sql, args, err := query.OrderBy("1").PlaceholderFormat(sq.Colon).ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, nil
}

// Constraints returns all the constraints of a given table.
func (o *oracle) Constraints(table TableRef) (string, []any, error) {
	var query sq.SelectBuilder

	query = sq.Select(
		`CONSTRAINT_NAME`,
		`CONSTRAINT_TYPE`,
	).
		From("USER_CONSTRAINTS").
		Where(sq.Eq{"TABLE_NAME": strings.ToUpper(table.Name)})

	if table.Schema != "" {
		query = sq.Select(
			`CONSTRAINT_NAME`,
			`CONSTRAINT_TYPE`,
		).
			From("ALL_CONSTRAINTS").
			Where(sq.Eq{"TABLE_NAME": strings.ToUpper(table.Name), "OWNER": strings.ToUpper(table.Schema)})
	}

	sql, args, err := query.PlaceholderFormat(sq.Colon).ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}

// Indexes returns the indexes of a table.
func (o *oracle) Indexes(table TableRef) (string, []any, error) {
	var query sq.SelectBuilder

	query = sq.Select("*").
		From("USER_INDEXES").
		Where(sq.Eq{"TABLE_NAME": strings.ToUpper(table.Name)})

	if table.Schema != "" {
		query = sq.Select("*").
			From("ALL_INDEXES").
			Where(sq.Eq{"TABLE_NAME": strings.ToUpper(table.Name), "OWNER": strings.ToUpper(table.Schema)})
	}

	sql, args, err := query.PlaceholderFormat(sq.Colon).ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, nil
}

// Catalog returns a the pointer to a DBNode instance,
// which is the root of the current Oracle database graph.
// It starts with the database itself,
// then the schemas and the correspondent lists of tables and views.
// Oracle topography:
//
//					 [Database]
//				       |
//				       v
//			     [Schemas]
//			      /     \
//			     v       v
//	 		 [Tables] 	[Views]
func (o *oracle) Catalog(ctx context.Context) (*DBNode, error) {
	rootID := fmt.Sprintf("db:%s", o.dbName)
	root := &DBNode{ID: rootID, Name: o.dbName, Type: "database"}
	queue := []*DBNode{root}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		var children []*DBNode
		var err error
		switch current.Type {
		case "database":
			if o.schema != "" {
				children = append(children, &DBNode{
					ID:       fmt.Sprintf("%s.s:%s", rootID, o.schema),
					Name:     o.schema,
					Type:     "schema",
					ParentID: rootID,
				})
			} else {
				children, err = o.fetchSchemas(ctx, current.Name)
			}
		case "schema":
			tables, err := o.fetchTables(ctx, current.Name, current.ID)
			if err != nil {
				return nil, err
			}
			children = append(children, tables...)

			views, err := o.fetchViews(ctx, current.Name, current.ID)
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
func (o *oracle) GetViewDefinition(view ViewRef) (string, []any, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Question)
	query, args, err := psql.
		Select("TEXT AS view_definition").
		From("ALL_VIEWS").
		Where(sq.Eq{
			"OWNER":     strings.ToUpper(view.Schema),
			"VIEW_NAME": strings.ToUpper(view.Name),
		}).
		ToSql()

	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

// fetchSchemas method lists all the schemas of the current database.
func (o *oracle) fetchSchemas(ctx context.Context, parentID string) ([]*DBNode, error) {
	query := `
		SELECT DISTINCT owner AS schema_name
		FROM all_tables
		ORDER BY owner
	`
	rows, err := o.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []*DBNode
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		schemas = append(schemas, &DBNode{

			ID:       fmt.Sprintf("%s.s:%s", parentID, name),
			Name:     name,
			Type:     "schema",
			ParentID: parentID,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schemas, nil
}

// fetchTables method returns a list of tables filtered by schema.
func (o *oracle) fetchTables(ctx context.Context, parentName, parentID string) ([]*DBNode, error) {
	query, args, err := sq.Select("TABLE_NAME").
		From("ALL_TABLES").
		Where(sq.Eq{"OWNER": strings.ToUpper(parentName)}).
		OrderBy("TABLE_NAME ASC").
		PlaceholderFormat(sq.Colon).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := o.db.Query(query, args...)
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

// fetchViews method returns a list of views filtered by schema.
func (o *oracle) fetchViews(ctx context.Context, parentName, parentID string) ([]*DBNode, error) {
	query, args, err := sq.Select("VIEW_NAME").
		From("ALL_VIEWS").
		Where(sq.Eq{"OWNER": strings.ToUpper(parentName)}).
		OrderBy("VIEW_NAME ASC").
		PlaceholderFormat(sq.Colon).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := o.db.Query(query, args...)
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
