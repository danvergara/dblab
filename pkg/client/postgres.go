package client

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
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
func newPostgres(schema string) *postgres {
	p := postgres{
		schema: schema,
	}

	return &p
}

// TableStructure returns a query string to get all the relevant information of a given table,
// under a schema.
func (p *postgres) TableStructure(tableName string) (string, []any, error) {
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
				sq.Eq{"c.table_schema": p.schema},
				sq.Eq{"c.table_name": tableName},
			},
		).
		ToSql()

	return query, args, err
}

// Constraints returns all the constraints of a given table, under a schema.
func (p *postgres) Constraints(tableName string) (string, []any, error) {
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
		Where(sq.Eq{"tc.table_name": tableName}).
		Where(sq.Eq{"tc.table_schema": p.schema}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}
	return sql, args, err
}

// Indexes returns the indexes of a table, under a schema.
func (p *postgres) Indexes(tableName string) (string, []any, error) {
	query := sq.Select("*").
		From("pg_indexes").
		Where(sq.Eq{"tableName": tableName}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sql, args, err
}

func (p *postgres) Catalog(ctx context.Context) (*DBNode, error) {
	// root := DBNode{
	// 	ID:   "databases-postgres",
	// 	Name: "users",
	// 	Type: "database",
	// 	Children: []DBNode{
	// 		{
	// 			ID:   "schemas-public",
	// 			Name: "Public",
	// 			Type: "schema",
	// 			Children: []DBNode{
	// 				{
	// 					ID:   "schemas-public-users",
	// 					Name: "users",
	// 					Type: "table",
	// 				},
	// 				{
	// 					ID:   "schemas-public-employees",
	// 					Name: "employees",
	// 					Type: "table",
	// 				},
	// 			},
	// 		},
	// 		{
	// 			ID:   "schemas-products",
	// 			Name: "Products",
	// 			Type: "schema",
	// 			Children: []DBNode{
	// 				{
	// 					ID:   "schemas-products-products",
	// 					Name: "products",
	// 					Type: "table",
	// 				},
	// 				{
	// 					ID:   "schemas-products-prices",
	// 					Name: "prices",
	// 					Type: "table",
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	rootID := fmt.Sprintf("db:%s-%s", uuid.New().String(), "")
	root := &DBNode{ID: rootID, Name: p.dbName, Type: "database"}
	queue := []*DBNode{root}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		var children []*DBNode
		var err error
		switch current.Type {
		case "database":
			children, err = p.fetchSchemas(ctx, current.Name)
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

func (p *postgres) fetchSchemas(ctx context.Context, parentName string) ([]*DBNode, error) {
	return nil, nil
}
