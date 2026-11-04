package client

import (
	"testing"

	"github.com/danvergara/dblab/pkg/drivers"
	"github.com/stretchr/testify/require"
)

func TestSanitizeIdentifierByDriver(t *testing.T) {
	tests := []struct {
		name    string
		queryer databaseQuerier
		input   string
		want    string
	}{
		{name: "sqlite", queryer: &sqlite{}, input: `odd"name`, want: `"odd""name"`},
		{name: "postgres", queryer: &postgres{}, input: `odd"name`, want: `"odd""name"`},
		{name: "mysql", queryer: &mysql{}, input: "odd`name", want: "`odd``name`"},
		{name: "mssql", queryer: &mssql{}, input: "odd]name", want: "[odd]]name]"},
		{name: "oracle", queryer: &oracle{}, input: `odd"name`, want: `"odd""name"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.queryer.SanitizeIdentifier(tt.input))
		})
	}
}

func TestContentQueryQuotesIdentifiersForEachDriver(t *testing.T) {
	tests := []struct {
		name    string
		driver  string
		queryer databaseQuerier
		want    string
	}{
		{
			name:    "postgres",
			driver:  drivers.Postgres,
			queryer: &postgres{},
			want:    `SELECT * FROM "my schema"."order details" LIMIT 7 OFFSET 3;`,
		},
		{
			name:    "mysql",
			driver:  drivers.MySQL,
			queryer: &mysql{},
			want:    "SELECT * FROM `order details` LIMIT 7 OFFSET 3;",
		},
		{
			name:    "sqlite",
			driver:  drivers.SQLite,
			queryer: &sqlite{},
			want:    `SELECT * FROM "order details" LIMIT 7 OFFSET 3;`,
		},
		{
			name:    "mssql",
			driver:  drivers.SQLServer,
			queryer: &mssql{},
			want:    "SELECT * FROM [my schema].[order details] ORDER BY (SELECT NULL) OFFSET 3 ROWS FETCH NEXT 7 ROWS ONLY",
		},
		{
			name:    "oracle",
			driver:  drivers.Oracle,
			queryer: &oracle{},
			want:    `SELECT * FROM "my schema"."order details" OFFSET 3 ROWS FETCH NEXT 7 ROWS ONLY`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contentQuery(tt.driver, tt.queryer, "my schema", "order details", 7, 3)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestMySQLMetadataQueriesQuoteIdentifiers(t *testing.T) {
	queryer := &mysql{}
	table := TableRef{Name: "order details"}

	query, args, err := queryer.TableStructure(table)
	require.NoError(t, err)
	require.Equal(t, "DESCRIBE `order details`;", query)
	require.Empty(t, args)

	query, args, err = queryer.Indexes(table)
	require.NoError(t, err)
	require.Equal(t, "SHOW INDEX FROM `order details`", query)
	require.Empty(t, args)
}
