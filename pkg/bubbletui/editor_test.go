package bubbletui

import "testing"

func TestQueryAtCursor(t *testing.T) {
	content := "SELECT 1;\nSELECT 2;\nSELECT 3;"

	got := queryAtCursor(content, 1, 9)
	want := "SELECT 2;"

	if got != want {
		t.Fatalf("queryAtCursor() = %q, want %q", got, want)
	}
}
