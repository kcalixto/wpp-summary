package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run . <action(calculate,fix)> <chat.txt>")
		os.Exit(1)
	}

	action := os.Args[1]
	filePath := os.Args[2]

	// Abrir o arquivo
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Erro ao abrir arquivo: %v\n", err)
		os.Exit(1)
	}

	defer file.Close()
	switch action {
	case "calculate":
		CalculateTokens(file)
	case "fix":
		FixTokens(file)
	default:
		fmt.Println("Ação inválida. Use 'calculate' ou 'fix'.")
		os.Exit(1)
	}
}
