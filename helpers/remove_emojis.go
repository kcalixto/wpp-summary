package helpers

import (
	"regexp"
	"strings"
)

func RemoveEmojis(line *Line) *Line {
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=~`[]{};':\",.<>?/\\|éêáãíóôúçÁÉÊÍÓÔÚÇÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÐÑÒÓÔÕÖØÙÚÛÜÝÞßàáâãäåæçèéêëìíîïðñòóôõöøùúûüýþÿ"

	// Remover caracteres não ASCII (que inclui emojis)
	const name = 0
	const message = 1
	for i, part := range []string{
		line.Name,
		line.Message,
	} {
		var result strings.Builder
		for _, r := range part {
			if r < 128 || strings.ContainsRune(allowedChars, r) {
				result.WriteRune(r)
			}
		}

		// Remover espaços extras e fazer trim
		cleaned := strings.TrimSpace(result.String())

		// Remover espaços duplicados
		spaceRegex := regexp.MustCompile(`\s+`)
		cleaned = spaceRegex.ReplaceAllString(cleaned, " ")

		if i == name {
			if cleaned == "" {
				line.Name = "Unknown"
			} else {
				line.Name = cleaned
			}
		} else {
			line.Message = cleaned
		}
	}

	return line
}
