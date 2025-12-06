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
	words := strings.Fields(body)
	if len(words) <= 30 {
		return body
	}
	excerpt := strings.Join(words[:30], " ")
	return excerpt
}
