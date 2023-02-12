package pagination

import (
	"fmt"
	"math"
)

// Manager handles the pagination.
type Manager struct {
	totalPages   int
	currentPage  int
	limit        uint
	offset       int
	currentTable string
}

// New returns a pointer to a Manager instance.
func New(limit uint, count int, currentTable string) (*Manager, error) {
	m := Manager{
		limit:        limit,
		currentPage:  1,
		currentTable: currentTable,
	}

	m.setOffset()

	if err := m.setTotalPages(count); err != nil {
		return nil, err
	}

	return &m, nil
}

// NextPage increases the value currentPage.
func (m *Manager) NextPage() error {
	if m.currentPage+1 > m.totalPages {
		return fmt.Errorf("current page should not be greater than the total pages count")
	}

	m.currentPage++
	m.setOffset()

	return nil
}

// PreviousPage decreases the value of currentPage.
func (m *Manager) PreviousPage() error {
	if m.currentPage-1 <= 0 {
		return fmt.Errorf("current page should not be less than 0")
	}

	m.currentPage--
	m.setOffset()

	return nil
}

// Offset returns the limit.
func (m *Manager) Offset() int {
	return m.offset
}

// Limit returns the limit.
func (m *Manager) Limit() uint {
	return m.limit
}

// TotalPages returns the total pages count.
func (m *Manager) TotalPages() int {
	return m.totalPages
}

// CurrentPage returns the currentPage value.
func (m *Manager) CurrentPage() int {
	return m.currentPage
}

// SetCurrentTable sets the current table name.
func (m *Manager) SetCurrentTable(tableName string) {
	m.currentTable = tableName
}

// CurrentTable sets the current table name.
func (m *Manager) CurrentTable() string {
	return m.currentTable
}

// setOffset calculates the offset based of the current page and the limit.
func (m *Manager) setOffset() {
	m.offset = (m.currentPage - 1) * int(m.limit)
}

// setTotalPages total pages = count / limit, if the limit is greater than 0.
func (m *Manager) setTotalPages(count int) error {
	m.totalPages = int(math.Ceil(float64(count) / float64(m.limit)))

	return nil
}
