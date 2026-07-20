package bubbletui

import "testing"

func TestQueryAtCursor(t *testing.T) {
	content := "SELECT 1;\nSELECT 2;\nSELECT 3;"

	tests := []struct {
		name string
		row  int
		col  int
		want string
	}{
		{name: "first query at start", row: 0, col: 0, want: "SELECT 1"},
		{name: "second query in middle", row: 1, col: 4, want: "SELECT 2"},
		{name: "second query on semicolon", row: 1, col: 8, want: "SELECT 2"},
		{name: "third query at semicolon", row: 2, col: 8, want: "SELECT 3"},
		{name: "third query after end of line", row: 2, col: 9, want: content},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := queryAtCursor(content, tt.row, tt.col)
			if got != tt.want {
				t.Fatalf("queryAtCursor(%d, %d) = %q, want %q", tt.row, tt.col, got, tt.want)
			}
		})
	}
}
