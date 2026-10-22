package bubbletui

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
)

// Foreground = Fg
// Background = Bg
// LeftForeground = LFg

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
	Help      HelpTheme
}

type ResultSetTheme struct {
	TableBorder       color.Color
	ActiveTableBorder color.Color
	Header            HeaderTheme
	SelectedContent   SelectedContentTheme
}

type HeaderTheme struct {
	Fg       color.Color
	BorderFg color.Color
}

type SelectedContentTheme struct {
	Fg color.Color
	Bg color.Color
}

// ModalTheme is them used for the ConnectModel, HistoryModel and the SchemaModel models.
// Those model share the same styles and colors since they are modal used to show options.
type ModalTheme struct {
	BorderFg               color.Color
	ListTitleBg            color.Color
	ListTitleFg            color.Color
	NormalTitle            color.Color
	SelectedTitle          color.Color
	SelectedTitleBorderLFg color.Color
}

type SidebarTheme struct {
	BorderFg            color.Color
	FocusedBorderFg     color.Color
	TextFg              color.Color
	FocusedTextFg       color.Color
	FocusedTextBg       color.Color
	FocusedTextBorderFg color.Color
}

type EditorTheme struct {
	Border        color.Color
	FocusedBorder color.Color
}

type BannerTheme struct {
	Text   color.Color
	Border color.Color
}

type HelpTheme struct {
	FullKeyFg       color.Color
	FullDescFg      color.Color
	FullSeparatorFg color.Color
}

func GetTheme(theme string) Theme {
	switch strings.ToLower(strings.TrimSpace(theme)) {
	case "catppuccin":
		return catppuccinMochaTheme()
	case "tokyonight":
		return tokyoNightTheme()
	case "gruvbox":
		return gruvboxDarkTheme()
	case "nord":
		return nordTheme()
	case "dracula":
		return draculaTheme()
	case "monokai":
		return monokaiTheme()
	case "onedark":
		return oneDarkTheme()
	case "cyberpunk":
		fallthrough
	default:
		// dblab-cyberpunk is the default theme.
		return cyberPunkTheme()
	}
}

func catppuccinMochaTheme() Theme {
	mauve := lipgloss.Color("#cba6f7")
	lavender := lipgloss.Color("#b4befe")
	peach := lipgloss.Color("#fab387")
	yellow := lipgloss.Color("#f9e2af")
	text := lipgloss.Color("#cdd6f4")
	subtext := lipgloss.Color("#a6adc8")
	surface1 := lipgloss.Color("#45475a")
	surface0 := lipgloss.Color("#313244")
	mantle := lipgloss.Color("#181825")
	crust := lipgloss.Color("#11111b")

	return Theme{
		Banner: BannerTheme{
			Text:   mauve,
			Border: lavender,
		},
		Modal: ModalTheme{
			BorderFg:               surface1,
			NormalTitle:            text,
			SelectedTitle:          peach,
			SelectedTitleBorderLFg: mauve,
			ListTitleBg:            mantle,
			ListTitleFg:            peach,
		},
		ResultSet: ResultSetTheme{
			TableBorder:       surface0,
			ActiveTableBorder: mauve,
			Header: HeaderTheme{
				Fg:       peach,
				BorderFg: mauve,
			},
			SelectedContent: SelectedContentTheme{
				Fg: crust,
				Bg: peach,
			},
		},
		Sidebar: SidebarTheme{
			BorderFg:            surface0,
			FocusedBorderFg:     mauve,
			TextFg:              text,
			FocusedTextFg:       peach,
			FocusedTextBg:       mantle,
			FocusedTextBorderFg: yellow,
		},
		Editor: EditorTheme{
			Border:        surface0,
			FocusedBorder: mauve,
		},
		Help: HelpTheme{
			FullKeyFg:       peach,
			FullDescFg:      text,
			FullSeparatorFg: subtext,
		},
	}
}

func tokyoNightTheme() Theme {
	blue := lipgloss.Color("#7aa2f7")
	cyan := lipgloss.Color("#7dcfff")
	magenta := lipgloss.Color("#bb9af7")
	orange := lipgloss.Color("#ff9e64")
	fg := lipgloss.Color("#c0caf5")
	comment := lipgloss.Color("#565f89")
	borderDark := lipgloss.Color("#292e42")
	borderLit := lipgloss.Color("#3b4261")
	bgDark := lipgloss.Color("#16161e")
	panelBg := lipgloss.Color("#1f2335")

	return Theme{
		Banner: BannerTheme{
			Text:   blue,
			Border: cyan,
		},
		Modal: ModalTheme{
			BorderFg:               borderLit,
			NormalTitle:            fg,
			SelectedTitle:          orange,
			SelectedTitleBorderLFg: magenta,
			ListTitleBg:            panelBg,
			ListTitleFg:            cyan,
		},
		ResultSet: ResultSetTheme{
			TableBorder:       borderDark,
			ActiveTableBorder: cyan,
			Header: HeaderTheme{
				Fg:       cyan,
				BorderFg: magenta,
			},
			SelectedContent: SelectedContentTheme{
				Fg: bgDark,
				Bg: orange,
			},
		},
		Sidebar: SidebarTheme{
			BorderFg:            borderDark,
			FocusedBorderFg:     cyan,
			TextFg:              fg,
			FocusedTextFg:       cyan,
			FocusedTextBg:       panelBg,
			FocusedTextBorderFg: magenta,
		},
		Editor: EditorTheme{
			Border:        borderDark,
			FocusedBorder: cyan,
		},
		Help: HelpTheme{
			FullKeyFg:       cyan,
			FullDescFg:      fg,
			FullSeparatorFg: comment,
		},
	}
}

func gruvboxDarkTheme() Theme {
	yellow := lipgloss.Color("#fabd2f")
	orange := lipgloss.Color("#fe8019")
	aqua := lipgloss.Color("#8ec07c")
	green := lipgloss.Color("#b8bb26")
	fg := lipgloss.Color("#ebdbb2")
	gray := lipgloss.Color("#928374")
	bg0 := lipgloss.Color("#282828")
	bg1 := lipgloss.Color("#3c3836")
	bg2 := lipgloss.Color("#504945")

	return Theme{
		Banner: BannerTheme{
			Text:   aqua,
			Border: green,
		},
		Modal: ModalTheme{
			BorderFg:               bg2,
			NormalTitle:            fg,
			SelectedTitle:          yellow,
			SelectedTitleBorderLFg: orange,
			ListTitleBg:            bg1,
			ListTitleFg:            yellow,
		},
		ResultSet: ResultSetTheme{
			TableBorder:       bg1,
			ActiveTableBorder: orange,
			Header: HeaderTheme{
				Fg:       yellow,
				BorderFg: orange,
			},
			SelectedContent: SelectedContentTheme{
				Fg: bg0,
				Bg: yellow,
			},
		},
		Sidebar: SidebarTheme{
			BorderFg:            bg1,
			FocusedBorderFg:     orange,
			TextFg:              fg,
			FocusedTextFg:       yellow,
			FocusedTextBg:       bg1,
			FocusedTextBorderFg: orange,
		},
		Editor: EditorTheme{
			Border:        bg1,
			FocusedBorder: orange,
		},
		Help: HelpTheme{
			FullKeyFg:       yellow,
			FullDescFg:      fg,
			FullSeparatorFg: gray,
		},
	}
}

func nordTheme() Theme {
	nord4 := lipgloss.Color("#d8dee9")
	nord3 := lipgloss.Color("#4c566a")
	nord2 := lipgloss.Color("#434c5e")
	nord1 := lipgloss.Color("#3b4252")
	nord0 := lipgloss.Color("#2e3440")
	nord7 := lipgloss.Color("#8fbcbb")
	nord8 := lipgloss.Color("#88c0d0")
	nord9 := lipgloss.Color("#81a1c1")
	nord12 := lipgloss.Color("#d08770")
	nord13 := lipgloss.Color("#ebcb8b")

	return Theme{
		Banner: BannerTheme{
			Text:   nord9,
			Border: nord7,
		},
		Modal: ModalTheme{
			BorderFg:               nord2,
			NormalTitle:            nord4,
			SelectedTitle:          nord13,
			SelectedTitleBorderLFg: nord8,
			ListTitleBg:            nord1,
			ListTitleFg:            nord13,
		},
		ResultSet: ResultSetTheme{
			TableBorder:       nord2,
			ActiveTableBorder: nord8,
			Header: HeaderTheme{
				Fg:       nord13,
				BorderFg: nord8,
			},
			SelectedContent: SelectedContentTheme{
				Fg: nord0,
				Bg: nord13,
			},
		},
		Sidebar: SidebarTheme{
			BorderFg:            nord2,
			FocusedBorderFg:     nord8,
			TextFg:              nord4,
			FocusedTextFg:       nord13,
			FocusedTextBg:       nord1,
			FocusedTextBorderFg: nord12,
		},
		Editor: EditorTheme{
			Border:        nord2,
			FocusedBorder: nord8,
		},
		Help: HelpTheme{
			FullKeyFg:       nord8,
			FullDescFg:      nord4,
			FullSeparatorFg: nord3,
		},
	}
}

func draculaTheme() Theme {
	bg := lipgloss.Color("#282a36")
	currLine := lipgloss.Color("#44475a")
	fg := lipgloss.Color("#f8f8f2")
	comment := lipgloss.Color("#6272a4")
	cyan := lipgloss.Color("#8be9fd")
	green := lipgloss.Color("#50fa7b")
	pink := lipgloss.Color("#ff79c6")
	purple := lipgloss.Color("#bd93f9")
	yellow := lipgloss.Color("#f1fa8c")

	return Theme{
		Banner: BannerTheme{
			Text:   purple,
			Border: cyan,
		},
		Modal: ModalTheme{
			BorderFg:               currLine,
			NormalTitle:            fg,
			SelectedTitle:          yellow,
			SelectedTitleBorderLFg: pink,
			ListTitleBg:            currLine,
			ListTitleFg:            pink,
		},
		ResultSet: ResultSetTheme{
			TableBorder:       currLine,
			ActiveTableBorder: pink,
			Header: HeaderTheme{
				Fg:       yellow,
				BorderFg: pink,
			},
			SelectedContent: SelectedContentTheme{
				Fg: bg,
				Bg: green,
			},
		},
		Sidebar: SidebarTheme{
			BorderFg:            currLine,
			FocusedBorderFg:     pink,
			TextFg:              fg,
			FocusedTextFg:       yellow,
			FocusedTextBg:       currLine,
			FocusedTextBorderFg: pink,
		},
		Editor: EditorTheme{
			Border:        currLine,
			FocusedBorder: pink,
		},
		Help: HelpTheme{
			FullKeyFg:       yellow,
			FullDescFg:      fg,
			FullSeparatorFg: comment,
		},
	}
}

func monokaiTheme() Theme {
	bg := lipgloss.Color("#272822")
	panelBg := lipgloss.Color("#1e1f1c")
	border := lipgloss.Color("#49483e")
	fg := lipgloss.Color("#f8f8f2")
	comment := lipgloss.Color("#75715e")
	yellow := lipgloss.Color("#e6db74")
	green := lipgloss.Color("#a6e22e")
	orange := lipgloss.Color("#fd971f")
	pink := lipgloss.Color("#f92672")
	cyan := lipgloss.Color("#66d9ef")

	return Theme{
		Banner: BannerTheme{
			Text:   orange,
			Border: cyan,
		},
		Modal: ModalTheme{
			BorderFg:               border,
			NormalTitle:            fg,
			SelectedTitle:          yellow,
			SelectedTitleBorderLFg: pink,
			ListTitleBg:            panelBg,
			ListTitleFg:            yellow,
		},
		ResultSet: ResultSetTheme{
			TableBorder:       border,
			ActiveTableBorder: pink,
			Header: HeaderTheme{
				Fg:       yellow,
				BorderFg: pink,
			},
			SelectedContent: SelectedContentTheme{
				Fg: bg,
				Bg: green,
			},
		},
		Sidebar: SidebarTheme{
			BorderFg:            border,
			FocusedBorderFg:     pink,
			TextFg:              fg,
			FocusedTextFg:       yellow,
			FocusedTextBg:       panelBg,
			FocusedTextBorderFg: orange,
		},
		Editor: EditorTheme{
			Border:        border,
			FocusedBorder: pink,
		},
		Help: HelpTheme{
			FullKeyFg:       yellow,
			FullDescFg:      fg,
			FullSeparatorFg: comment,
		},
	}
}

func oneDarkTheme() Theme {
	bg := lipgloss.Color("#282c34")
	panelBg := lipgloss.Color("#21252b")
	border := lipgloss.Color("#3e4451")
	borderLit := lipgloss.Color("#4b5263")
	fg := lipgloss.Color("#abb2bf")
	comment := lipgloss.Color("#5c6370")
	blue := lipgloss.Color("#61afef")
	cyan := lipgloss.Color("#56b6c2")
	purple := lipgloss.Color("#c678dd")
	yellow := lipgloss.Color("#e5c07b")

	return Theme{
		Banner: BannerTheme{
			Text:   blue,
			Border: cyan,
		},
		Modal: ModalTheme{
			BorderFg:               borderLit,
			NormalTitle:            fg,
			SelectedTitle:          yellow,
			SelectedTitleBorderLFg: purple,
			ListTitleBg:            panelBg,
			ListTitleFg:            blue,
		},
		ResultSet: ResultSetTheme{
			TableBorder:       border,
			ActiveTableBorder: blue,
			Header: HeaderTheme{
				Fg:       yellow,
				BorderFg: purple,
			},
			SelectedContent: SelectedContentTheme{
				Fg: bg,
				Bg: yellow,
			},
		},
		Sidebar: SidebarTheme{
			BorderFg:            border,
			FocusedBorderFg:     blue,
			TextFg:              fg,
			FocusedTextFg:       yellow,
			FocusedTextBg:       panelBg,
			FocusedTextBorderFg: purple,
		},
		Editor: EditorTheme{
			Border:        border,
			FocusedBorder: blue,
		},
		Help: HelpTheme{
			FullKeyFg:       yellow,
			FullDescFg:      fg,
			FullSeparatorFg: comment,
		},
	}
}

func cyberPunkTheme() Theme {
	return Theme{
		Banner: BannerTheme{
			Text:   purpleColor,
			Border: greenColor,
		},
		Modal: ModalTheme{
			BorderFg:               neonPurpleColor,
			NormalTitle:            whiteColor,
			SelectedTitle:          cyberGreenColor,
			SelectedTitleBorderLFg: hiMagentaColor,
			ListTitleBg:            darkPurpleColor,
			ListTitleFg:            hiMagentaColor,
		},
		ResultSet: ResultSetTheme{
			TableBorder:       darkPurpleColor,
			ActiveTableBorder: neonPurpleColor,
			Header: HeaderTheme{
				Fg:       cyberGreenColor,
				BorderFg: hiMagentaColor,
			},
			SelectedContent: SelectedContentTheme{
				Fg: blackColor,
				Bg: cyberGreenColor,
			},
		},
		Sidebar: SidebarTheme{
			BorderFg:            darkPurpleColor,
			FocusedBorderFg:     neonPurpleColor,
			TextFg:              whiteColor,
			FocusedTextFg:       cyberGreenColor,
			FocusedTextBg:       darkPurpleColor,
			FocusedTextBorderFg: hiMagentaColor,
		},
		Editor: EditorTheme{
			FocusedBorder: neonPurpleColor,
			Border:        darkPurpleColor,
		},
		Help: HelpTheme{
			FullKeyFg:       cyberGreenColor,
			FullDescFg:      whiteColor,
			FullSeparatorFg: mutedGreenColor,
		},
	}
}
