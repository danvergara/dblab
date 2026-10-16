package client

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danvergara/dblab/pkg/command"
	"github.com/danvergara/dblab/pkg/drivers"
	"github.com/stretchr/testify/require"
)

func TestSQLiteCatalogNamesMetadata(t *testing.T) {
	c, err := New(command.Options{
		Driver: drivers.SQLite,
		DBName: filepath.Join(t.TempDir(), "catalog.db"),
		Limit:  5,
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, c.DB().Close()) })

	for _, name := range []string{"widgets", "order-items", "order details", "select", `item"labels`} {
		t.Run(name, func(t *testing.T) {
			quoted := `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
			_, err := c.DB().Exec(fmt.Sprintf("CREATE TABLE %s (id INTEGER PRIMARY KEY, label TEXT)", quoted))
			require.NoError(t, err)
			_, err = c.DB().Exec(fmt.Sprintf("INSERT INTO %s VALUES (1, 'dummy')", quoted))
			require.NoError(t, err)
			_, err = c.DB().Exec(fmt.Sprintf("CREATE VIEW %s AS SELECT * FROM %s", `"view-`+strings.ReplaceAll(name, `"`, `""`)+`"`, quoted))
			require.NoError(t, err)

			catalog, err := c.Catalog(context.Background())
			require.NoError(t, err)
			var tableFound, viewFound bool
			for _, node := range catalog.Children {
				switch node.EntityName {
				case name:
					tableFound = true
					metadata, err := c.Metadata(TableRef{Name: node.EntityName})
					require.NoError(t, err)
					require.Equal(t, [][]string{{"1", "dummy"}}, metadata.TableContent.Rows)
					require.Len(t, metadata.Structure.Rows, 2)
				case "view-" + name:
					viewFound = true
					metadata, err := c.ViewMetadata(ViewRef{Name: node.EntityName})
					require.NoError(t, err)
					require.Equal(t, [][]string{{"1", "dummy"}}, metadata.TableContent.Rows)
					require.Len(t, metadata.ViewDef.Rows, 1)
				}
			}
			require.True(t, tableFound)
			require.True(t, viewFound)
		})
	}
}
