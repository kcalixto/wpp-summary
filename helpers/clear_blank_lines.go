package helpers

import "strings"

func ClearIfJustTimeAndName(line *Line) *Line {
	msg := strings.ToLower(line.Message)
	replace := []string{
		"\n",
		".",
		" ",
		"obrigado",
		"bom dia",
		"boa tarde",
		"boa noite",
	}

	for _, r := range replace {
		msg = strings.ReplaceAll(msg, r, "")
	}

	if msg == "" {
		return nil
	}
	return line
}
