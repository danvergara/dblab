package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationLifeCycle(t *testing.T) {
	// This is an example of the pagination life cycle.
	// It starts off with a limit of 50 and count from a sql query of 10000.
	// The total pages count should be equal to count / limit = 200
	// I'm calling the NextPage 3 times and PreviousPage two times.
	limit := 50
	count := 150

	m, err := New(limit, count, "users")
	assert.NoError(t, err)

	for i := 0; i < (count/m.Limit())-1; i++ {
		_ = m.NextPage()
	}

	assert.Equal(t, count/m.Limit(), m.CurrentPage())
	assert.Equal(t, (m.CurrentPage()-1)*limit, m.Offset())

	err = m.NextPage()
	assert.Error(t, err)

	for i := 0; i < (count/m.Limit())-1; i++ {
		_ = m.PreviousPage()
	}

	assert.Equal(t, 1, m.CurrentPage())
	assert.Equal(t, 0, m.Offset())

	err = m.PreviousPage()
	assert.Error(t, err)

	m.SetCurrentTable("products")
	ct := m.CurrentTable()
	assert.Equal(t, "products", ct)
}
