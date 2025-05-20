package main

import (
	"fmt"
	"os"
	"strings"
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
		if len(os.Args) < 4 {
			fmt.Println("Uso: go run . fix <chat.txt> <destino.txt>")
			os.Exit(1)
		}
		dest := os.Args[3]
		if dest == "" {
			fmt.Println("Uso: go run . fix <chat.txt> <destino.txt>")
			os.Exit(1)
		}
		if !strings.HasSuffix(dest, ".txt") {
			fmt.Println("O arquivo de destino deve ter a extensão .txt")
			os.Exit(1)
		}
		FixTokens(file, dest)

		destFile, err := os.Open(dest)
		if err != nil {
			fmt.Printf("Erro ao abrir arquivo de destino: %v\n", err)
			os.Exit(1)
		}
		defer destFile.Close()

		err = CalculateTokens(destFile)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	default:
		fmt.Println("Ação inválida. Use 'calculate' ou 'fix'.")
		os.Exit(1)
	}
}
