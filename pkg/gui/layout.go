package gui

import (
	"github.com/danvergara/gocui"
	"github.com/fatih/color"
)

var (
	green   = color.New(color.FgGreen).Add(color.Bold)
	options = []string{"Rows", "Structure", "Constraints", "Indexes"}
)

func (gui *Gui) setLayout() {
	maxX, maxY := gui.g.Size()

	banner := NewBannerWidget("banner", 0, 0, int(0.19*float32(maxX)), int(0.14*float32(maxY)), "dblab", gocui.ColorMagenta)

	tables := NewTableWidget("tables", 0, int(0.16*float32(maxY)), int(0.19*float32(maxX)), int(0.94*float32(maxY)), "Tables", gocui.ColorGreen, gocui.ColorBlack, gui)

	navigation := NewNavigationWidget("navigation", int(0.2*float32(maxX)), 0, maxX-1, int(0.07*float32(maxY)), "Navigation", options)

	editor := NewEditorWidget("query", int(0.2*float32(maxX)), int(0.09*float32(maxY)), maxX-1, int(0.27*float32(maxY)), "SQL Query")
	indexes := NewOutputWidget("indexes", int(0.2*float32(maxX)), int(0.29*float32(maxY)), maxX-1, int(0.94*float32(maxY)), "Indexes", "Please select a table!")
	constraints := NewOutputWidget("constraints", int(0.2*float32(maxX)), int(0.29*float32(maxY)), maxX-1, int(0.94*float32(maxY)), "Constraints", "Please select a table!")
	structure := NewOutputWidget("structure", int(0.2*float32(maxX)), int(0.29*float32(maxY)), maxX-1, int(0.94*float32(maxY)), "Structure", "Please select a table!")
	rows := NewOutputWidget("rows", int(0.2*float32(maxX)), int(0.29*float32(maxY)), maxX-1, int(0.94*float32(maxY)), "Rows", "Type the sql query above. Press Ctrl-c to quit.")

	rowsPerPage := NewLabelWidget("rows-per-page", int(0.81*float32(maxX)), int(0.96*float32(maxY)), "00 rows", gocui.ColorWhite)
	currentPage := NewLabelWidget("current-page", int(0.90*float32(maxX)), int(0.96*float32(maxY)), "00", gocui.ColorWhite)
	slash := NewLabelWidget("slash", int(0.92*float32(maxX)), int(0.96*float32(maxY)), "/", gocui.ColorWhite)
	totalPages := NewLabelWidget("total-pages", int(0.93*float32(maxX)), int(0.96*float32(maxY)), "00", gocui.ColorWhite)

	back := NewButtonWidget("back", int(0.86*float32(maxX)), int(0.96*float32(maxY)), "< BACK", gocui.ColorGreen)
	next := NewButtonWidget("next", int(0.95*float32(maxX)), int(0.96*float32(maxY)), "NEXT >", gocui.ColorGreen)

	gui.g.SetManager(banner, tables, navigation, editor, indexes, constraints, structure, rows, rowsPerPage, back, currentPage, slash, totalPages, next)
}
