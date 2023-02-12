package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationLifeCycle(t *testing.T) {
	type given struct {
		count     int
		limit     uint
		tableName string
	}

	type expected struct {
		lastPage      int
		currentOffset int
		totalPages    int
		tableName     string
	}

	var tests = []struct {
		name     string
		given    given
		expected expected
	}{
		{
			name: "count is higher than limit",
			given: given{
				limit:     50,
				count:     150,
				tableName: "products",
			},
			expected: expected{
				lastPage:      3,
				currentOffset: 0,
				totalPages:    3,
				tableName:     "products",
			},
		},
		{
			name: "limit is higher than count",
			given: given{
				limit:     150,
				count:     50,
				tableName: "products",
			},
			expected: expected{
				lastPage:      1,
				currentOffset: 0,
				totalPages:    1,
				tableName:     "products",
			},
		},
		{
			name: "count is odd and higher than limit",
			given: given{
				limit:     50,
				count:     151,
				tableName: "products",
			},
			expected: expected{
				lastPage:      4,
				currentOffset: 0,
				totalPages:    4,
				tableName:     "products",
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			m, err := New(tc.given.limit, tc.given.count, "users")
			assert.NoError(t, err)

			assert.Equal(t, tc.expected.totalPages, m.TotalPages())

			for i := 0; i < m.TotalPages(); i++ {
				_ = m.NextPage()
			}

			assert.Equal(t, tc.expected.lastPage, m.CurrentPage())
			assert.Equal(t, (m.CurrentPage()-1)*int(tc.given.limit), m.Offset())

			err = m.NextPage()
			assert.Error(t, err)

			for i := 0; i < m.TotalPages(); i++ {
				_ = m.PreviousPage()
			}
			assert.Equal(t, 1, m.CurrentPage())
			assert.Equal(t, tc.expected.currentOffset, m.Offset())

			err = m.PreviousPage()
			assert.Error(t, err)

			m.SetCurrentTable(tc.given.tableName)
			ct := m.CurrentTable()
			assert.Equal(t, tc.expected.tableName, ct)
		})
	}
}
