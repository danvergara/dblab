package bubbletui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type focusState int

const (
	focusInput focusState = iota
	focusList
)

type model struct {
	input textinput.Model
	list  list.Model
}
