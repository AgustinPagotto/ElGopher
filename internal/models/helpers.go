package models

import (
	"bytes"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func slugifyTitle(title string) string {
	newTitle := norm.NFD.String(title)
	var runes []rune
	for _, r := range newTitle {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			runes = append(runes, r)
		}
	}
	newTitleSt := string(runes)
	newTitleB := []byte(newTitleSt)
	newTitleB = bytes.ToLower(newTitleB)
	newTitleB = bytes.ReplaceAll(newTitleB, []byte("  "), []byte(" "))
	newTitleB = bytes.ReplaceAll(newTitleB, []byte(" "), []byte("-"))
	return string(newTitleB)
}

func generateExcerpt(body string) string {
	// Strip markdown syntax
	cleaned := body
	cleaned = strings.ReplaceAll(cleaned, "#", "")
	cleaned = strings.ReplaceAll(cleaned, "*", "")
	cleaned = strings.ReplaceAll(cleaned, "_", "")
	cleaned = strings.ReplaceAll(cleaned, "`", "")
	cleaned = strings.ReplaceAll(cleaned, "[", "")
	cleaned = strings.ReplaceAll(cleaned, "]", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	cleaned = strings.TrimSpace(cleaned)

	const maxLength = 155
	if len(cleaned) <= maxLength {
		return cleaned
	}

	truncated := cleaned[:maxLength]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}
	return truncated + "..."
}
