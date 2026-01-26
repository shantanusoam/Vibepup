package ui

import (
	"unicode/utf8"
)

// ClampWidth trims a string to fit within width runes.
func ClampWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= width {
		return s
	}
	runes := []rune(s)
	return string(runes[:width])
}
