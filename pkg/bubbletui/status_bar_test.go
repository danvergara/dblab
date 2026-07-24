package bubbletui

import (
	"testing"

	"github.com/danvergara/dblab/pkg/client"
	"github.com/danvergara/dblab/pkg/command"
)

func TestModelStatusBarTracksEditorModeChanges(t *testing.T) {
	c, _ := client.New(command.Options{Driver: "sqlite", Host: "/tmp/file.bd"})
	kb := &command.TUIKeyMap{}
	editor := NewEditor(kb)
	m := &Model{editor: editor, statusBar: NewStatusBar(editor.mode, kb, c)}

	updated, _ := m.Update(modeChangeMsg{mode: InsertMode})
	model, ok := updated.(Model)
	if !ok {
		t.Fatalf("expected Model, got %T", updated)
	}

	if got := model.statusBar.mode.String(); got != InsertMode.String() {
		t.Fatalf("expected status bar mode %q, got %q", InsertMode.String(), got)
	}
}
