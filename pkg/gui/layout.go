package gui

import (
	"fmt"

	"github.com/danvergara/gocui"
	"github.com/fatih/color"
)

var (
	green   = color.New(color.FgGreen).Add(color.Bold)
	options = []string{"Rows", "Structure", "Constraints", "Indexes"}
)

func (gui *Gui) setLayout() {
	// banners.
	banner := NewBannerWidget(
		"banner",
		0,
		0,
		0.19,
		0.14,
		"dblab",
		gocui.ColorMagenta,
	)

	// table.
	tables := NewTableWidget(
		"tables",
		0,
		0.16,
		0.19,
		0.94,
		"Tables",
		gocui.ColorGreen,
		gocui.ColorBlack,
		gui,
	)

	// navigation widget.
	navigation := NewNavigationWidget(
		"navigation",
		0.2,
		0,
		-1,
		0.07,
		"Navigation",
		options,
	)

	// editor.
	editor := NewEditorWidget(
		"query",
		0.2,
		0.09,
		-1,
		0.27,
		"SQL Query",
	)

	// outputs.
	indexes := NewOutputWidget(
		"indexes",
		0.2,
		0.29,
		-1,
		0.94,
		"Indexes",
		"Please select a table!",
	)
	constraints := NewOutputWidget(
		"constraints",
		0.2,
		0.29,
		-1,
		0.94,
		"Constraints",
		"Please select a table!",
	)
	structure := NewOutputWidget(
		"structure",
		0.2,
		0.29,
		-1,
		0.94,
		"Structure",
		"Please select a table!",
	)
	rows := NewOutputWidget(
		"rows",
		0.2,
		0.29,
		-1,
		0.94,
		"Rows",
		"Type the sql query above. Press Ctrl-c to quit.",
	)

	// labels.
	currentPage := NewLabelWidget(
		"current-page",
		0.84,
		0.96,
		fmt.Sprintf("%4d", 0),
		gocui.ColorWhite,
	)
	slash := NewLabelWidget(
		"slash",
		0.87,
		0.96,
		"/",
		gocui.ColorWhite,
	)
	totalPages := NewLabelWidget(
		"total-pages",
		0.89,
		0.96,
		fmt.Sprintf("%4d", 0),
		gocui.ColorWhite,
	)

	// buttons.
	back := NewButtonWidget(
		"back",
		0.80,
		0.96,
		"< BACK",
		gocui.ColorGreen,
	)
	next := NewButtonWidget(
		"next",
		0.92,
		0.96,
		"NEXT >",
		gocui.ColorGreen,
	)

	gui.g.SetManager(
		banner,
		tables,
		navigation,
		editor,
		indexes,
		constraints,
		structure,
		rows,
		back,
		currentPage,
		slash,
		totalPages,
		next,
	)
}
