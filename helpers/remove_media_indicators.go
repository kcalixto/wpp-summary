package helpers

import "strings"

func RemoveMediaIndicators(line string) string {
	// Remover indicadores de mídia
	// Exemplo: "[image omitted]", "[sticker omitted]", "[video omitted]", "[Contact card omitted]"
	indicators := []string{
		"image omitted",
		"sticker omitted",
		"video omitted",
		"document omitted",
		"audio omitted",
		"Contact card omitted",
		"\u200E",
		"<This message was edited>",
		"This message was deleted.",
	}
	for _, indicator := range indicators {
		line = strings.ReplaceAll(line, indicator, "")
	}

	return line
}
