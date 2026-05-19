package client

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
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
		"c.column_name",
		"c.is_nullable",
		"c.data_type",
		"c.character_maximum_length",
		"c.numeric_precision",
		"c.numeric_scale",
		"c.ordinal_position",
		"tc.constraint_type AS pkey",
	).
		From("information_schema.columns AS c").
		LeftJoin(
			`information_schema.constraint_column_usage AS ccu
					ON c.table_schema = ccu.table_schema
						AND c.table_name = ccu.table_name
						AND c.column_name = ccu.column_name`,
		).
		LeftJoin(
			`information_schema.table_constraints AS tc
					ON ccu.constraint_schema = tc.constraint_schema
						AND ccu.constraint_name = tc.constraint_name`,
		).
		Where(
			sq.And{
				sq.Eq{"c.table_schema": table.Schema},
				sq.Eq{"c.table_name": table.Name},
			},
		).
		ToSql()

	return query, args, err
}

// Constraints returns all the constraints of a given table, under a schema.
func (p *postgres) Constraints(table TableRef) (string, []any, error) {
	var (
		query sq.SelectBuilder
		sql   string
	)

	query = sq.Select(
		`tc.constraint_name`,
		`tc.table_name`,
		`tc.constraint_type`,
	).
		From("information_schema.table_constraints AS tc").
		Where(sq.Eq{"tc.table_name": table.Name}).
		Where(sq.Eq{"tc.table_schema": table.Schema}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, err
}

// Indexes returns the indexes of a table, under a schema.
func (p *postgres) Indexes(table TableRef) (string, []any, error) {
	query := sq.Select("*").
		From("pg_indexes").
		Where(sq.And{
			sq.Eq{"schemaname": table.Schema},
			sq.Eq{"tablename": table.Name},
		}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}

func (p *postgres) Catalog(ctx context.Context) (*DBNode, error) {
	rootID := fmt.Sprintf("db:%s", p.dbName)
	root := &DBNode{ID: rootID, Name: p.dbName, Type: "database"}
	queue := []*DBNode{root}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		var children []*DBNode
		var err error
		switch current.Type {
		case "database":
			if p.schema != "" {
				children = append(children, &DBNode{
					ID:       fmt.Sprintf("%s.s:%s", rootID, p.schema),
					Name:     p.schema,
					Type:     "schema",
					ParentID: rootID,
				})
			} else {
				children, err = p.fetchSchemas(ctx, current.Name)
			}
		case "schema":
			children, err = p.fetchTables(ctx, current.Name, current.ID)
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

func (p *postgres) fetchSchemas(ctx context.Context, parentID string) ([]*DBNode, error) {
	query, args, err := sq.Select("schema_name").
		From("information_schema.schemata").
		Where(sq.NotEq{
			"schema_name ": []string{"information_schema", "pg_catalog"},
		}).
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

func (p *postgres) fetchTables(ctx context.Context, parentName, parentID string) ([]*DBNode, error) {
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
