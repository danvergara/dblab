package bubbletui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
)

const (
	// dblabJSONStyle names the chroma style, used for the query
	// resulsets's JSON syntax highlighting.
	dblabJSONStyle = "dblab-json"
	// dblabSQLStyle names the chroma style, used for the query
	// editor's SQL syntax highlighting.
	dblabSQLStyle = "dblab-sql"
)

const (
	red   = "#FF0000"
	white = "#ffffff"
	black = "#000000"

	green          = "#1fb009" // Grass green
	platinum       = "#E0E0E0"
	purple         = "#800080"
	cyberGreen     = "#39FF14" // High-visibility neon green
	hiMagenta      = "#FF00FF" // High-visibility Magenta
	mutedGreen     = "#2ECC71" // Softer green for standard text
	neonPurple     = "#BF40BF" // Bright purple for highlights
	darkPurple     = "#4B0082" // Off-white for readability
	deepGreen      = "#116800" // Japanese Laurel
	greenishYellow = "#a2cb27" // Lime green
	neonOrange     = "#FF6600" // neonOrange
	neonViolet     = "#9D00FF" // neonViolet
	mediumGray     = "#999999" // mediumGray
	charcoal       = "#393939" // charchoal, dark gray, onyx
	soapstone      = "#fcfcfc" // off-white, snow
	ashGray        = "#bbbbbb"
	tundora        = "#444444" // dark eel, emperor
	darkGray       = "#323232" // dark gray, dune
)

var (
	// General purpose lipgloss Colors.
	greenColor      = lipgloss.Color(green)
	purpleColor     = lipgloss.Color(purple)
	cyberGreenColor = lipgloss.Color(cyberGreen)
	hiMagentaColor  = lipgloss.Color(hiMagenta)
	mutedGreenColor = lipgloss.Color(mutedGreen)
	neonPurpleColor = lipgloss.Color(neonPurple)
	darkPurpleColor = lipgloss.Color(darkPurple)
	whiteColor      = lipgloss.Color(white)
	blackColor      = lipgloss.Color(black)

	// StatusBar Colors.
	FocusText  = lipgloss.Color(deepGreen)
	FocusBg    = lipgloss.Color(greenishYellow)
	SchemaText = lipgloss.Color(charcoal)
	SchemaBg   = lipgloss.Color(soapstone)
	KbOddText  = lipgloss.Color(ashGray)
	KbOddBg    = lipgloss.Color(tundora)
	KbEvenText = lipgloss.Color(mediumGray)
	KbEvenBg   = lipgloss.Color(darkGray)
)

// Register the dblab-cyberpunk-sql style the query editor highlights SQL with,
// built from the same neon palette used across the rest of the app (see
// cyberGreen, neonViolet, neonPurple and friends in bubbletui.go).
var _ = styles.Register(chroma.MustNewStyle(dblabSQLStyle, chroma.StyleEntries{
	chroma.Background:    platinum,             // default text
	chroma.Keyword:       neonViolet + " bold", // SELECT/FROM/...
	chroma.NameBuiltin:   neonPurple + " bold", // VARCHAR/INT/...
	chroma.Name:          mutedGreen,           // identifiers
	chroma.LiteralString: cyberGreen,           // 'strings'
	chroma.LiteralNumber: neonOrange,           // numbers
	chroma.Operator:      platinum,             // + - * / = ...
	chroma.Punctuation:   platinum,             // ; , ( ) ...
	chroma.Comment:       mediumGray + " italic",
	chroma.Error:         red + " bold",
}))

// Register the dblab-cyberpunk style Highlight json inspects.
var _ = styles.Register(chroma.MustNewStyle(dblabJSONStyle, chroma.StyleEntries{
	chroma.NameTag:             hiMagenta + " bold", // keys
	chroma.LiteralStringDouble: cyberGreen,          // strings
	chroma.LiteralString:       cyberGreen,
	chroma.LiteralNumber:       neonPurple,           // numbers
	chroma.KeywordConstant:     neonPurple + " bold", // true/false/null
	chroma.Punctuation:         platinum,             // {} [] : ,
	chroma.Error:               red + " bold",
	chroma.Background:          platinum, // default fg, no bg set
}))

type Theme struct {
	Banner    BannerTheme
	ResultSet ResultSetTheme
	Modal     ModalTheme
	Sidebar   SidebarTheme
	Editor    EditorTheme
}

type ResultSetTheme struct {
	TableBorder       color.Color
	ActiveTableBorder color.Color
	Header            HeaderTheme
	SelectedContent   SelectedContentTheme
}

type HeaderTheme struct {
	Foreground       color.Color
	BorderForeground color.Color
}

type SelectedContentTheme struct {
	Foreground color.Color
	Background color.Color
}

// ModalTheme is them used for the ConnectModel, HistoryModel and the SchemaModel models.
// Those model share the same styles and colors since they are modal used to show options.
type ModalTheme struct {
	BorderForeground                  color.Color
	ListTitleBackground               color.Color
	ListTitleForeground               color.Color
	NormalTitle                       color.Color
	SelectedTitle                     color.Color
	SelectedTitleBorderLeftForeground color.Color
}

type SidebarTheme struct {
	BorderForeground        color.Color
	FocusedBorderForeground color.Color

	TextForeground color.Color

	FocusedTextForeground       color.Color
	FocusedTextBackground       color.Color
	FocusedTextBorderForeground color.Color
}

type EditorTheme struct {
	Border        color.Color
	FocusedBorder color.Color
}

type BannerTheme struct {
	Text   color.Color
	Border color.Color
}

func GetTheme(theme string) Theme {
	switch theme {
	default:
		// dblab-cyberpunk is the default theme.
		return cyberPunkTheme()
	}
}

func cyberPunkTheme() Theme {
	return Theme{
		Banner: BannerTheme{
			Text:   purpleColor,
			Border: greenColor,
		},
		Modal: ModalTheme{
			NormalTitle:                       whiteColor,
			SelectedTitle:                     cyberGreenColor,
			SelectedTitleBorderLeftForeground: hiMagentaColor,
			ListTitleBackground:               darkPurpleColor,
			ListTitleForeground:               hiMagentaColor,
		},
		ResultSet: ResultSetTheme{
			TableBorder:       darkPurpleColor,
			ActiveTableBorder: neonPurpleColor,
			Header: HeaderTheme{
				Foreground:       cyberGreenColor,
				BorderForeground: hiMagentaColor,
			},
			SelectedContent: SelectedContentTheme{
				Foreground: blackColor,
				Background: cyberGreenColor,
			},
		},
		Sidebar: SidebarTheme{
			BorderForeground:            darkPurpleColor,
			FocusedBorderForeground:     neonPurpleColor,
			TextForeground:              whiteColor,
			FocusedTextForeground:       cyberGreenColor,
			FocusedTextBackground:       darkPurpleColor,
			FocusedTextBorderForeground: hiMagentaColor,
		},
		Editor: EditorTheme{
			FocusedBorder: neonPurpleColor,
			Border:        darkPurpleColor,
		},
	}
}
