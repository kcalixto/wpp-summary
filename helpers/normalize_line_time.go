package helpers

import (
	"fmt"
	"regexp"
	"strconv"
)

func NormalizeLineTime(line *Line) *Line {
	// Normalizar a data/hora
	// Exemplo: [25/04/25, 08:34:22] -> 25/04/25 08h
	// Exemplo: [25/04/25, 19:34:22] -> 25/04/25 19h
	regex := regexp.MustCompile(`\[(\d{2}/\d{2}/\d{2}), (\d{2}):(\d{2}):\d{2}\]`)

	return &Line{
		Time: regex.ReplaceAllStringFunc(line.Time, func(match string) string {
			// Extrair os grupos capturados
			submatch := regex.FindStringSubmatch(match)

			date := submatch[1] // Data (25/04/25)
			hour := submatch[2] // Hora (08)

			// Converter para exibir apenas a hora
			hourInt, _ := strconv.Atoi(hour)

			// Formatar como "25/04/25 08h" ou "25/04/25 19h"
			return fmt.Sprintf("[%s %02dh]", date, hourInt)
		}),
		Name:    line.Name,
		Message: line.Message,
	}
}
