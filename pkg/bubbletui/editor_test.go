package bubbletui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryAtCursor(t *testing.T) {
	content := "SELECT 1;\nSELECT 2;\nSELECT 3;"

	tests := []struct {
		name string
		row  int
		want string
	}{
		{name: "first query at start", row: 0, want: "SELECT 1;"},
		{name: "second query in middle", row: 1, want: "SELECT 2;"},
		{name: "second query on semicolon", row: 1, want: "SELECT 2;"},
		{name: "third query at semicolon", row: 2, want: "SELECT 3;"},
		{name: "third query after end of line", row: 3, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := queryAtCursor(content, tt.row)
			assert.Equal(t, tt.want, got)
		})
	}
}
