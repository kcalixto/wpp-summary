package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func CalculateTokens(file *os.File) error {
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

	// Estatísticas básicas
	lineCount := len(strings.Split(fullText, "\n"))
	charCount := utf8.RuneCountInString(fullText)
	byteCount := len(fullText)
	wordCount := len(strings.Fields(fullText))

	// Estimar tokens
	tokenEstimate := estimateTokens(fullText)

	// Estimar tokens com prompt
	promptText := readPromptFile()

	promptTokens := estimateTokens(promptText)
	totalTokens := promptTokens + tokenEstimate

	// Resultados
	fmt.Println("=== Estatísticas do Arquivo ===")
	fmt.Printf("Tamanho em bytes: %d\n", byteCount)
	fmt.Printf("Número de caracteres: %d\n", charCount)
	fmt.Printf("Número de palavras: %d\n", wordCount)
	fmt.Printf("Número de linhas: %d\n", lineCount)

	fmt.Println("\n=== Estimativa de Tokens ===")
	fmt.Printf("Tokens do prompt: ~%d\n", promptTokens)
	fmt.Printf("Tokens do arquivo: ~%d\n", tokenEstimate)
	fmt.Printf("Total estimado: ~%d tokens\n", totalTokens)

	// Estimar se está dentro dos limites do modelo
	fmt.Println("\n=== Compatibilidade com Claude 3.5 Sonnet ===")

	const claudeSonnetLimit = 50000

	if totalTokens <= claudeSonnetLimit {
		fmt.Printf("✅ O arquivo está dentro do limite (~%d tokens). O Claude 3.5 Sonnet pode processar.\n", totalTokens)
		// Estimativa de custo muito aproximada
		fmt.Printf("   Estimativa de custo: $%.4f para input + $%.4f para output (considerando resposta curta)\n",
			float64(totalTokens)*0.000003, float64(2000)*0.000015)
	} else {
		fmt.Printf("❌ O arquivo excede o limite do Claude 3.5 Sonnet em ~%d tokens.\n", totalTokens-claudeSonnetLimit)
		fmt.Println("   Recomendações:")
		fmt.Println("   - Divida o arquivo em partes menores")
		fmt.Println("   - Remova mensagens menos importantes")
		fmt.Println("   - Pré-processe o arquivo para reduzir conteúdo irrelevante")
		return fmt.Errorf("o arquivo excede o limite do Claude 3.5 Sonnet em ~%d tokens", totalTokens-claudeSonnetLimit)
	}

	return nil
}

func readPromptFile() string {
	promptFilePath := ".goosehints"
	promptFile, err := os.Open(promptFilePath)
	if err != nil {
		fmt.Printf("Erro ao abrir arquivo %s: %v\n", promptFilePath, err)
		os.Exit(1)
	}

	promptScanner := bufio.NewScanner(promptFile)
	var promptText string
	for promptScanner.Scan() {
		promptText += promptScanner.Text() + "\n"
	}

	if err := promptScanner.Err(); err != nil {
		fmt.Printf("Erro ao ler prompt: %v\n", err)
		os.Exit(1)
	}

	return promptText
}

// Função para estimar tokens na string conforme regras comuns de tokenização
func estimateTokens(text string) int {
	// Método simplificado de estimativa de tokens
	// Regras aproximadas para tokenização:
	// 1. Palavras separadas por espaço são tokens
	// 2. Pontuação geralmente é um token separado
	// 3. Números podem ser múltiplos tokens dependendo do comprimento
	// 4. Emojis e caracteres especiais são tokens separados

	// Remover quebras de linha extras e espaços extras
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// Contar palavras básicas (separadas por espaço)
	words := strings.Fields(text)
	tokenCount := len(words)

	// Ajustar para pontuação e símbolos (estimativa aproximada)
	punctCount := 0
	for _, r := range text {
		if unicode.IsPunct(r) || (r > 127 && !unicode.IsLetter(r) && !unicode.IsNumber(r)) {
			punctCount++
		}
	}

	// Ajustar para números longos e palavras longas
	// (números e palavras longas geralmente são divididos em múltiplos tokens)
	longItemsAdjustment := 0
	for _, word := range words {
		// Se for um número ou palavra longa (>6 caracteres), pode gerar tokens adicionais
		if len(word) > 6 {
			// Estimar tokens adicionais com base no comprimento
			longItemsAdjustment += (len(word) / 6)
		}
	}

	// Considerar caracteres não-ASCII (como emojis, acentos, etc.)
	nonAsciiCount := 0
	for _, r := range text {
		if r > 127 {
			nonAsciiCount++
		}
	}

	// Fórmula de estimativa final
	// Esta é uma aproximação e não uma contagem exata como faria um tokenizador real
	estimatedTokens := tokenCount + (punctCount / 2) + (longItemsAdjustment / 2) + (nonAsciiCount / 4)

	return estimatedTokens
}
