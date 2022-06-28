package gui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/danvergara/gocui"
)

// ButtonWidget struct used to build buttons.
type ButtonWidget struct {
	name  string
	x, y  int
	w     int
	color gocui.Attribute
	label string
}

// NewButtonWidget returns a pointer to a ButtonWidget instance.
func NewButtonWidget(name string, x, y int, label string, color gocui.Attribute) *ButtonWidget {
	return &ButtonWidget{name: name, x: x, y: y, w: len(label) + 1, label: label, color: color}
}

// Layout implements the gocui.Manager interface.
func (w *ButtonWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.FrameColor = w.color

		fmt.Fprint(v, w.label)
	}

	return nil
}

// LabelWidget struct used to display data to dynamic data to the user.
type LabelWidget struct {
	name  string
	x, y  int
	w     int
	color gocui.Attribute
	label string
}

// NewLabelWidget returns a pointer to a LabelWidget instance.
func NewLabelWidget(name string, x, y int, label string, color gocui.Attribute) *ButtonWidget {
	return &ButtonWidget{name: name, x: x, y: y, w: len(label) + 1, label: label, color: color}
}

// Layout implements the gocui.Manager interface.
func (w *LabelWidget) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.FrameColor = w.color

		fmt.Fprint(v, w.label)
	}

	return nil
}

// BannerWidget struct used to build the banner where
// we show the name of the app.
type BannerWidget struct {
	name           string
	x0, y0, x1, y1 int
	color          gocui.Attribute
	label          string
}

// NewBannerWidget returns a pointer to a BannerWidget instance.
func NewBannerWidget(name string, x0, y0, x1, y1 int, label string, color gocui.Attribute) *BannerWidget {
	return &BannerWidget{name: name, x0: x0, y0: y0, x1: x1, y1: y1, label: label, color: color}
}

// Layout implements the gocui.Manager interface.
func (w *BannerWidget) Layout(g *gocui.Gui) error {
	if v, err := g.SetView(w.name, w.x0, w.y0, w.x1, w.y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.FrameColor = w.color
		myFigure := figure.NewFigure(w.label, "", true)
		figure.Write(v, myFigure)
	}

	return nil
}

// TableWidget struct used to build the section for the
// sql tables content.
type TableWidget struct {
	name             string
	x0, y0, x1, y1   int
	gui              *Gui
	bgcolor, fgcolor gocui.Attribute
	label            string
}

// NewTableWidget returns a pointer to a TableWidget instance.
func NewTableWidget(name string, x0, y0, x1, y1 int, label string, bgcolor, fgcolor gocui.Attribute, gui *Gui) *TableWidget {
	return &TableWidget{
		name:    name,
		x0:      x0,
		y0:      y0,
		x1:      x1,
		y1:      y1,
		gui:     gui,
		label:   label,
		bgcolor: bgcolor,
		fgcolor: fgcolor,
	}
}

// Layout implements the gocui.Manager interface.
func (w *TableWidget) Layout(g *gocui.Gui) error {
	if v, err := g.SetView(w.name, w.x0, w.y0, w.x1, w.y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		if err := w.gui.showTables(); err != nil {
			return err
		}

		v.Title = w.label
		v.Highlight = true
		v.SelBgColor = w.bgcolor
		v.SelFgColor = w.fgcolor
	}

	return nil
}

// NavigationWidget struct used to show the navigation panel.
type NavigationWidget struct {
	name           string
	x0, y0, x1, y1 int
	options        []string
	label          string
}

// NewNavigationWidget returns a pointer to a NavigationWidget instance.
func NewNavigationWidget(name string, x0, y0, x1, y1 int, label string, options []string) *NavigationWidget {
	return &NavigationWidget{name: name, x0: x0, y0: y0, x1: x1, y1: y1, label: label, options: options}
}

// Layout implements the gocui.Manager interface.
func (w *NavigationWidget) Layout(g *gocui.Gui) error {
	if v, err := g.SetView(w.name, w.x0, w.y0, w.x1, w.y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = w.label

		tmpOptions := make([]string, len(w.options))
		copy(tmpOptions, w.options)
		tmpOptions[0] = green.Sprint("Rows")

		fmt.Fprint(v, strings.Join(tmpOptions, "   "))
	}

	return nil
}

// EditorWidget struct used as an editor to perform queries to the databases.
type EditorWidget struct {
	name           string
	x0, y0, x1, y1 int
	label          string
}

// NewEditorWidget returns a pointer to a EditorWidget instance.
func NewEditorWidget(name string, x0, y0, x1, y1 int, label string) *EditorWidget {
	return &EditorWidget{name: name, x0: x0, y0: y0, x1: x1, y1: y1, label: label}
}

// Layout implements the gocui.Manager interface.
func (w *EditorWidget) Layout(g *gocui.Gui) error {
	if v, err := g.SetView(w.name, w.x0, w.y0, w.x1, w.y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = w.label
		v.Editable = true
		v.Wrap = true
		v.Highlight = true

		if _, err := g.SetCurrentView(w.name); err != nil {
			return err
		}
	}

	return nil
}

// OutputWidget struct used to show important data to the user
// based off the context.
type OutputWidget struct {
	name           string
	x0, y0, x1, y1 int
	label          string
	initMsg        string
}

// NewOutputWidget returns a pointer to a OutputWidget instance.
func NewOutputWidget(name string, x0, y0, x1, y1 int, label string, initMsg string) *OutputWidget {
	return &OutputWidget{name: name, x0: x0, y0: y0, x1: x1, y1: y1, label: label, initMsg: initMsg}
}

// Layout implements the gocui.Manager interface.
func (w *OutputWidget) Layout(g *gocui.Gui) error {
	if v, err := g.SetView(w.name, w.x0, w.y0, w.x1, w.y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = w.label

		fmt.Fprintln(v, w.initMsg)
	}

	return nil
}
