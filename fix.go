package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kcalixto/wpp-summary/helpers"
)

/*
	Estratégias para redução do arquivo

[X] Remover todas as entradas "[image omitted]", "[sticker omitted]", "[video omitted]", "[Contact card omitted]"
[ ] Simplificar cabeçalhos de data/hora para apenas data se as mensagens são do mesmo dia
[ ] Mensagens curtas ou reações ("Ameeeeei", "Thanks", "Obrigado", "Bom dia, pessoal!")
[ ] Consolidar mensagens sequenciais do mesmo remetente
*/

func FixTokens(file *os.File) {
	// Leitura do arquivo
	scanner := bufio.NewScanner(file)
	var fullText string
	for scanner.Scan() {
		fullText += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Erro ao ler arquivo: %v\n", err)
		os.Exit(1)
	}

	fullText = StartFix(fullText)

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

func StartFix(str string) string {
	// PRE FIXER
	preFixedLines := make([]*helpers.Line, 0)
	for i, lineStr := range strings.Split(str, "\n") {
		for _, fix := range []helpers.LinePreFixer{
			helpers.RemoveMediaIndicators,
		} {
			lineStr = fix(lineStr)
		}

		previousLine := &helpers.Line{}
		if i > 0 {
			previousLine = preFixedLines[len(preFixedLines)-1]
		}

		next, line := ReadLine(previousLine, lineStr)
		if line == nil {
			// Se a linha for nil, significa que não há mais nada a fazer
			// Então, podemos parar de processar essa linha
			continue
		}
		if next { // juntou com a linha anterior
			idx := len(preFixedLines) - 1
			if idx < 0 {
				// Se não houver linha anterior, não há nada a fazer
				continue
			}

			preFixedLines[idx] = line
			continue
		}

		preFixedLines = append(preFixedLines, line)
	}

	// FIXER
	fixedLines := make([]*helpers.Line, 0)
	for _, line := range preFixedLines {
		for _, fix := range []helpers.LineFixer{
			helpers.RemoveEmojis,
			helpers.NormalizeLineTime,
			helpers.ClearIfJustTimeAndName,
		} {
			line = fix(line)
			if line == nil {
				// Se a linha for nil, significa que não há mais nada a fazer
				// Então, podemos parar de processar essa linha
				break
			}
		}
		if line != nil {
			fixedLines = append(fixedLines, line)
		}
	}

	// POST FIXER
	postFixedLines := make([]*helpers.Line, 0)
	for _, fix := range []helpers.LinesPostFixer{
		helpers.GroupByDateTime,
	} {
		postFixedLines = fix(fixedLines)
	}

	// * Recria a linha com os dados fixados
	fullText := ""
	for _, line := range postFixedLines {
		if line.Name == "" {
			if line.Time == "" { // message
				fullText += line.Message + "\n"
				continue
			}
			if line.Message == "" { // time header
				fullText += "\n" + line.Time + "\n"
				continue
			}
		}

		fullText += fmt.Sprintf("%s: %s", line.Name, line.Message) + "\n" // name header

		// fullText += fmt.Sprintf("%s %s: %s", line.Time, line.Name, line.Message) + "\n" // name header
	}

	return fullText
}

func ReadLine(previousLine *helpers.Line, lineStr string) (next bool, line *helpers.Line) {
	re := regexp.MustCompile(`^\[([^\]]+)\]\s*~([^:]+):\s*(.*)$`)
	matches := re.FindStringSubmatch(lineStr)
	if matches == nil {
		// empty line or just message
		if lineStr == "" {
			// empty line
			// fmt.Println("No match found")
			return false, nil
		}

		// just message
		// fmt.Println("Just message found")
		msg := lineStr
		if previousLine.Message != "" {
			msg = fmt.Sprintf("%s. %s", previousLine.Message, lineStr)
		}
		return true, &helpers.Line{
			Time:    previousLine.Time,
			Name:    previousLine.Name,
			Message: msg,
		}
	}

	line = &helpers.Line{
		Time:    fmt.Sprintf("[%s]", matches[1]),
		Name:    matches[2],
		Message: matches[3],
	}

	if previousLine.Name == line.Name {
		// Se o nome for o mesmo remetente da linha anterior, não é necessário criar uma nova linha
		// Então, podemos juntar as mensagens
		return true, &helpers.Line{
			Time:    line.Time,
			Name:    line.Name,
			Message: fmt.Sprintf("%s. %s", previousLine.Message, line.Message),
		}
	}

	return false, line
}
