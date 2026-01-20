package main

import (
	"fmt"
	"os"

	"ufmt.br/luiz-crema/compilador/lexer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run . <arquivo_entrada>")
		fmt.Println("Exemplo: go run . exemplo.txt")
		os.Exit(1)
	}

	nomeArquivo := os.Args[1]

	fmt.Printf("Processando arquivo: %s\n", nomeArquivo)

	err := lexer.GerarTabelaLexemas(nomeArquivo)
	if err != nil {
		fmt.Printf("Erro ao processar arquivo: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Tabela de lexemas gerada com sucesso em lexemas.txt!")
}
