package bubbletui

import (
	"github.com/danvergara/dblab/pkg/config"
)

var icons map[string]string

func init() {
	icons = config.GetIcons()
}

func getIcon(icon string) string {
	if v, ok := icons[icon]; ok {
		return v
	} else {
		switch icon {
		case "database":
			return "⛃"
		case "schema":
			return "📁"
		case "table":
			return "📋"
		case "view":
			return "📑"
		}
		return ""
	}
}
