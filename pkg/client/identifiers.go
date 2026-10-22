package client

import "strings"

func quoteDoubleQuotedIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func quoteMySQLIdentifier(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

func quoteMSSQLIdentifier(name string) string {
	return "[" + strings.ReplaceAll(name, "]", "]]") + "]"
}
