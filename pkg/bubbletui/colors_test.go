package bubbletui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTheme(t *testing.T) {
	tests := []struct {
		name          string
		themeInput    string
		expectedTheme Theme
	}{
		{
			name:          "default fallback when empty",
			themeInput:    "",
			expectedTheme: cyberPunkTheme(),
		},
		{
			name:          "default fallback when unknown",
			themeInput:    "nonexistent-theme",
			expectedTheme: cyberPunkTheme(),
		},
		{
			name:          "cyberpunk explicit",
			themeInput:    "cyberpunk",
			expectedTheme: cyberPunkTheme(),
		},
		{
			name:          "cyberpunk uppercase with whitespace",
			themeInput:    "  CYBERPUNK  ",
			expectedTheme: cyberPunkTheme(),
		},
		{
			name:          "catppuccin",
			themeInput:    "catppuccin",
			expectedTheme: catppuccinMochaTheme(),
		},
		{
			name:          "tokyonight",
			themeInput:    "tokyonight",
			expectedTheme: tokyoNightTheme(),
		},
		{
			name:          "gruvbox",
			themeInput:    "gruvbox",
			expectedTheme: gruvboxDarkTheme(),
		},
		{
			name:          "nord",
			themeInput:    "nord",
			expectedTheme: nordTheme(),
		},
		{
			name:          "nord uppercase",
			themeInput:    "NORD",
			expectedTheme: nordTheme(),
		},
		{
			name:          "dracula",
			themeInput:    "dracula",
			expectedTheme: draculaTheme(),
		},
		{
			name:          "dracula mixed case",
			themeInput:    "Dracula",
			expectedTheme: draculaTheme(),
		},
		{
			name:          "monokai",
			themeInput:    "monokai",
			expectedTheme: monokaiTheme(),
		},
		{
			name:          "monokai uppercase",
			themeInput:    "MONOKAI",
			expectedTheme: monokaiTheme(),
		},
		{
			name:          "onedark",
			themeInput:    "onedark",
			expectedTheme: oneDarkTheme(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTheme(tt.themeInput)
			assert.Equal(t, tt.expectedTheme, got)
		})
	}
}

func TestThemesHaveValidColors(t *testing.T) {
	themes := []struct {
		name  string
		theme Theme
	}{
		{name: "cyberpunk", theme: cyberPunkTheme()},
		{name: "catppuccin", theme: catppuccinMochaTheme()},
		{name: "tokyonight", theme: tokyoNightTheme()},
		{name: "gruvbox", theme: gruvboxDarkTheme()},
		{name: "nord", theme: nordTheme()},
		{name: "dracula", theme: draculaTheme()},
		{name: "monokai", theme: monokaiTheme()},
		{name: "onedark", theme: oneDarkTheme()},
	}

	for _, tt := range themes {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.theme.Banner.Text)
			assert.NotNil(t, tt.theme.Banner.Border)

			assert.NotNil(t, tt.theme.Modal.BorderFg)
			assert.NotNil(t, tt.theme.Modal.NormalTitle)
			assert.NotNil(t, tt.theme.Modal.SelectedTitle)
			assert.NotNil(t, tt.theme.Modal.SelectedTitleBorderLFg)
			assert.NotNil(t, tt.theme.Modal.ListTitleBg)
			assert.NotNil(t, tt.theme.Modal.ListTitleFg)

			assert.NotNil(t, tt.theme.ResultSet.TableBorder)
			assert.NotNil(t, tt.theme.ResultSet.ActiveTableBorder)
			assert.NotNil(t, tt.theme.ResultSet.Header.Fg)
			assert.NotNil(t, tt.theme.ResultSet.Header.BorderFg)
			assert.NotNil(t, tt.theme.ResultSet.SelectedContent.Fg)
			assert.NotNil(t, tt.theme.ResultSet.SelectedContent.Bg)

			assert.NotNil(t, tt.theme.Sidebar.BorderFg)
			assert.NotNil(t, tt.theme.Sidebar.FocusedBorderFg)
			assert.NotNil(t, tt.theme.Sidebar.TextFg)
			assert.NotNil(t, tt.theme.Sidebar.FocusedTextFg)
			assert.NotNil(t, tt.theme.Sidebar.FocusedTextBg)
			assert.NotNil(t, tt.theme.Sidebar.FocusedTextBorderFg)

			assert.NotNil(t, tt.theme.Editor.Border)
			assert.NotNil(t, tt.theme.Editor.FocusedBorder)

			assert.NotNil(t, tt.theme.Help.FullKeyFg)
			assert.NotNil(t, tt.theme.Help.FullDescFg)
			assert.NotNil(t, tt.theme.Help.FullSeparatorFg)

			assert.NotEmpty(t, tt.theme.SQL.Background)
			assert.NotEmpty(t, tt.theme.SQL.Comment)
			assert.NotEmpty(t, tt.theme.SQL.Error)
			assert.NotEmpty(t, tt.theme.SQL.Keyword)
			assert.NotEmpty(t, tt.theme.SQL.LiteralNumber)
			assert.NotEmpty(t, tt.theme.SQL.LiteralString)
			assert.NotEmpty(t, tt.theme.SQL.Name)
			assert.NotEmpty(t, tt.theme.SQL.NameBuiltin)
			assert.NotEmpty(t, tt.theme.SQL.Operator)
			assert.NotEmpty(t, tt.theme.SQL.Punctuation)

			assert.NotEmpty(t, tt.theme.JSON.Background)
			assert.NotEmpty(t, tt.theme.JSON.Punctuation)
			assert.NotEmpty(t, tt.theme.JSON.Error)
			assert.NotEmpty(t, tt.theme.JSON.KeywordConstant)
			assert.NotEmpty(t, tt.theme.JSON.LiteralNumber)
			assert.NotEmpty(t, tt.theme.JSON.LiteralString)
			assert.NotEmpty(t, tt.theme.JSON.LiteralStringDouble)
			assert.NotEmpty(t, tt.theme.JSON.NameTag)

			assert.NotPanics(t, func() {
				registerSQLSyntax(tt.theme.SQL)
				registerJSONSyntax(tt.theme.JSON)
			})
		})
	}
}
