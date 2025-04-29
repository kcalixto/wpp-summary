package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func FixTokens(file *os.File) {
	// Leitura do arquivo
	scanner := bufio.NewScanner(file)
	var fullText string
	for scanner.Scan() {
		line := scanner.Text()
		line = fixLine(line)
		fullText += line + "\n"
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Erro ao ler arquivo: %v\n", err)
		os.Exit(1)
	}

	// Estratégias para redução do arquivo
	//
	// Remover indicadores de mídia
	//
	// Remover todas as entradas "[image omitted]", "[sticker omitted]", "[video omitted]", "[Contact card omitted]"
	// Estas não contêm informação útil para análise
	//
	//
	// Remover metadados redundantes
	//
	// Simplificar cabeçalhos de data/hora para apenas data se as mensagens são do mesmo dia
	// Ou remover completamente se a data for evidente pelo contexto
	//
	//
	// Remover mensagens de baixa prioridade
	//
	// Mensagens curtas ou reações ("Ameeeeei", "Thanks", "Obrigado")
	// Saudações simples ("Bom dia, pessoal!")
	// Mensagens deletadas
	//
	//
	// Consolidar mensagens sequenciais do mesmo remetente
	//
	// Quando uma pessoa envia várias mensagens em sequência, agrupá-las
	//
	//
	// Focar em mensagens de alta prioridade
	//
	// Manter apenas comunicados da sub-síndica (Maria Eduarda)
	// Manter apenas comunicados do assistente de gestão (Jeff Lima)
	// Manter apenas informações sobre serviços essenciais

	fullText = groupByDateTime(fullText)

	// Escrever o texto reduzido em um novo arquivo
	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Printf("Erro ao criar arquivo de saída: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(fullText)
	if err != nil {
		fmt.Printf("Erro ao escrever no arquivo de saída: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Arquivo reduzido salvo como output.txt")
}

func groupByDateTime(fullText string) string {
	pattern := `\[(\d{2}/\d{2}/\d{2}\s\d{2})h\]`
	datetimes := make(map[string]string)

	for _, line := range strings.Split(fullText, "\n") {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(line)

		if len(matches) < 2 {
			continue
		}

		datetime := matches[1]
		if _, exists := datetimes[datetime]; !exists {
			datetimes[datetime] = line
		} else {
			datetimes[datetime] += "\n" + line
		}
	}

	// Recria o texto completo com as mensagens agrupadas
	fullText = ""
	for _, line := range datetimes {
		fullText += line + "\n"
	}

	return fullText
}

func fixLine(line string) string {
	line = removeMediaIndicators(line)
	line = normalizeLineTime(line)
	line = removeBlankLines(line)
	line = removeEmojis(line)
	line = fixName(line)
	line = clearIfJustTimeAndName(line)
	return line
}

func removeMediaIndicators(line string) string {
	// Remover indicadores de mídia
	// Exemplo: "[image omitted]", "[sticker omitted]", "[video omitted]", "[Contact card omitted]"
	indicators := []string{
		"image omitted",
		"sticker omitted",
		"video omitted",
		"Contact card omitted",
		"\u200E",
		"<This message was edited>",
	}
	for _, indicator := range indicators {
		line = strings.ReplaceAll(line, indicator, "")
	}

	return line
}

func normalizeLineTime(line string) string {
	// Normalizar a data/hora
	// Exemplo: [25/04/25, 08:34:22] -> 25/04/25 08h
	// Exemplo: [25/04/25, 19:34:22] -> 25/04/25 19h
	regex := regexp.MustCompile(`\[(\d{2}/\d{2}/\d{2}), (\d{2}):(\d{2}):\d{2}\]`)

	// Substituir pela formatação desejada
	result := regex.ReplaceAllStringFunc(line, func(match string) string {
		// Extrair os grupos capturados
		submatch := regex.FindStringSubmatch(match)
		if len(submatch) >= 4 {
			date := submatch[1] // Data (25/04/25)
			hour := submatch[2] // Hora (08)

			// Converter para exibir apenas a hora
			hourInt, _ := strconv.Atoi(hour)

			// Formatar como "25/04/25 08h" ou "25/04/25 19h"
			return fmt.Sprintf("[%s %02dh]", date, hourInt)
		}
		return match
	})

	return result
}

func removeBlankLines(line string) string { // unused
	// Remover linhas em branco
	if strings.TrimSpace(line) == " " {
		return ""
	}
	return line
}

func removeEmojis(line string) string {
	// Remover caracteres não ASCII (que inclui emojis)
	var result strings.Builder
	for _, r := range line {
		if r < 128 {
			result.WriteRune(r)
		}
	}

	// Remover espaços extras e fazer trim
	cleaned := strings.TrimSpace(result.String())

	// Remover espaços duplicados
	spaceRegex := regexp.MustCompile(`\s+`)
	cleaned = spaceRegex.ReplaceAllString(cleaned, " ")

	return cleaned
}

func fixName(line string) string {
	// Se não encontrar o padrão, retorna a linha original
	return line
}

func clearIfJustTimeAndName(line string) string {
	// Regex para capturar o formato: [DATA HORA] ~ NOME: +
	reWithSpace := regexp.MustCompile(`^\[([^\]]+)\]\s*~([^:]+ ):\s*(.*)$`)
	matches := reWithSpace.FindStringSubmatch(line)
	if matches == nil {
		return line
	}

	message := matches[3]
	if message == "" {
		return ""
	}

	// no space
	re := regexp.MustCompile(`^\[([^\]]+)\]\s*~([^:]+):\s*(.*)$`)
	matches = re.FindStringSubmatch(line)
	if matches == nil {
		return line
	}

	message = matches[3]
	if message == "" {
		return ""
	}

	return line
}
