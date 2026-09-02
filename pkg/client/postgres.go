package client

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var (
	postgresSchemaQuery = sq.Select("schema_name").
		From("information_schema.schemata").
		Where(sq.NotEq{
			"schema_name ": []string{"information_schema", "pg_catalog"},
		}).
		PlaceholderFormat(sq.Dollar)
)

// postgres struct is in charge of perform all the postgres related queries,
// without the client knowing.
type postgres struct {
	db     *sqlx.DB
	dbName string
	schema string
}

// a validation to see if postgres is implementing databaseQuerier.
var _ databaseQuerier = (*postgres)(nil)

// returns a pointer to a postgres, it receives an schema as a parameter.
func newPostgres(dbName, schema string, db *sqlx.DB) *postgres {
	p := postgres{
		dbName: dbName,
		db:     db,
		schema: schema,
	}

	return &p
}

// TableStructure returns a query string to get all the relevant information of a given table,
// under a schema.
func (p *postgres) TableStructure(table TableRef) (string, []any, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.Select(
		"a.attname AS column_name",
		"NOT a.attnotnull AS is_nullable",
		"pg_catalog.format_type(a.atttypid, a.atttypmod) AS data_type",
		"a.attnum AS ordinal_position",
		"COALESCE(i.indisprimary, false) AS pkey",
	).
		From("pg_catalog.pg_attribute a").
		Join("pg_catalog.pg_class c ON a.attrelid = c.oid").
		Join("pg_catalog.pg_namespace n ON c.relnamespace = n.oid").
		LeftJoin(`pg_catalog.pg_index i ON c.oid = i.indrelid 
                AND a.attnum = ANY(i.indkey) 
                AND i.indisprimary`).
		Where(
			sq.And{
				sq.Eq{"n.nspname": table.Schema},
				sq.Eq{"c.relname": table.Name},
				sq.Gt{"a.attnum": 0},
				sq.Expr("NOT a.attisdropped"),
			},
		).
		ToSql()

	return query, args, err
}

// Constraints returns all the constraints of a given table, under a schema.
func (p *postgres) Constraints(table TableRef) (string, []any, error) {
	sql, args, err := sq.Select(
		"con.conname AS constraint_name",
		"c.relname AS table_name",
		"con.contype AS constraint_type",
		"pg_catalog.pg_get_constraintdef(con.oid) AS constraint_definition",
	).
		From("pg_catalog.pg_constraint con").
		Join("pg_catalog.pg_class c ON con.conrelid = c.oid").
		Join("pg_catalog.pg_namespace n ON c.relnamespace = n.oid").
		Where(
			sq.And{
				sq.Eq{"n.nspname": table.Schema},
				sq.Eq{"c.relname": table.Name},
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, err
}

// Indexes returns the indexes of a table, under a schema.
func (p *postgres) Indexes(table TableRef) (string, []any, error) {
	sql, args, err := sq.Select(
		"n.nspname as schema_name",
		"ix.relname AS index_name",
		"i.indisunique AS is_unique",
		"i.indisprimary AS is_primary",
		"a.attname AS column_name",
		"pg_catalog.pg_get_indexdef(i.indexrelid) AS index_definition",
	).
		From("pg_catalog.pg_class t").
		Join("pg_catalog.pg_index i ON t.oid = i.indrelid").
		Join("pg_catalog.pg_class ix ON i.indexrelid = ix.oid").
		Join("pg_catalog.pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(i.indkey)").
		Join("pg_catalog.pg_namespace n ON t.relnamespace = n.oid").
		Where(
			sq.And{
				sq.Eq{"n.nspname": table.Schema},
				sq.Eq{"t.relname": table.Name},
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	return sql, args, err
}

// Catalog returns a the pointer to a DBNode instance,
// which is the root of the current PostgreSQL database graph.
// It starts with the database itself,
// then the schemas and the correspondent lists of tables and views.
// PostgreSQL topography:
//
//					 [Database]
//				       |
//				       v
//			     [Schemas]
//			      /     \
//			     v       v
//	 		 [Tables] 	[Views]
func (p *postgres) Catalog(ctx context.Context) (*DBNode, error) {
	rootID := fmt.Sprintf("db:%s", p.dbName)
	root := &DBNode{ID: rootID, Name: p.dbName, Type: "database"}
	queue := []*DBNode{root}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		children := make([]*DBNode, 0)
		var err error
		switch current.Type {
		case "database":
			if p.schema != "" {
				children = append(children, &DBNode{
					ID:         fmt.Sprintf("%s.s:%s", rootID, p.schema),
					Name:       p.schema,
					EntityName: p.schema,
					Type:       "schema",
					ParentID:   rootID,
				})
			} else {
				children, err = p.fetchSchemas(ctx, current.Name)
			}
		case "schema":
			tables, err := p.fetchTables(ctx, current.Name, current.ID)
			if err != nil {
				return nil, err
			}
			children = append(children, tables...)

			views, err := p.fetchViews(ctx, current.Name, current.ID)
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
func (p *postgres) GetViewDefinition(view ViewRef) (string, []any, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.
		Select().
		Column(sq.Expr("pg_get_viewdef(?::text::regclass, true) AS view_definition", fmt.Sprintf("%s.%s", view.Schema, view.Name))).
		ToSql()
	if err != nil {
		return "", nil, err
	}

	return query, args, nil
}

// fetchSchemas method lists all the schemas of the current database.
func (p *postgres) fetchSchemas(_ context.Context, parentID string) ([]*DBNode, error) {
	query, args, err := postgresSchemaQuery.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(query, args...)
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
			ID:         fmt.Sprintf("%s.s:%s", parentID, name),
			Name:       name,
			EntityName: name,
			Type:       "schema",
			ParentID:   parentID,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schemas, nil
}

// fetchTables method returns a list of tables filtered by schema.
func (p *postgres) fetchTables(_ context.Context, parentName, parentID string) ([]*DBNode, error) {
	query, args, err := sq.Select("table_name").
		From("information_schema.tables").
		Where(sq.Eq{"table_schema": parentName}).
		OrderBy("table_name").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(query, args...)
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
func (p *postgres) fetchViews(_ context.Context, parentName, parentID string) ([]*DBNode, error) {
	// 'v' is View, 'm' is Materialized View.
	query, args, err := sq.Select("c.relname AS view_name").
		From("pg_class c").
		Join("pg_namespace n ON n.oid = c.relnamespace").
		Where(sq.Eq{
			"n.nspname": parentName,
			"c.relkind": []string{"v", "m"},
		}).
		OrderBy("c.relname ASC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.db.Query(query, args...)
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

func (p *postgres) Schemas() (string, []any, error) {
	return postgresSchemaQuery.ToSql()
}

func (p *postgres) SetActiveSchema(schema string) (string, []any, error) {
	return fmt.Sprintf("set search_path = '%s'", schema), nil, nil
}
