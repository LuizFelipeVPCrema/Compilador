package main

import (
	"fmt"
	"os"

	"ufmt.br/luiz-crema/compilador/gerador_codigo"
	"ufmt.br/luiz-crema/compilador/lexer"
	"ufmt.br/luiz-crema/compilador/semantico"
	"ufmt.br/luiz-crema/compilador/sintatico"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run . <arquivo_entrada> [lexemas]")
		fmt.Println("Exemplo: go run . descricao/correto.python.txt")
		fmt.Println("         go run . descricao/correto.python.txt lexemas  (gera lexemas.txt)")
		os.Exit(1)
	}

	nomeArquivo := os.Args[1]
	arquivoSaida := "codigo.objeto.txt"
	gerarLexemas := len(os.Args) > 2 && os.Args[2] == "lexemas"

	fmt.Printf("Processando arquivo: %s\n", nomeArquivo)
	if gerarLexemas {
		if err := lexer.GerarTabelaLexemas(nomeArquivo); err != nil {
			fmt.Printf("Erro ao gerar lexemas: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Lexemas salvos em lexemas.txt")
	}

	lex, err := lexer.NovoLexer(nomeArquivo)
	if err != nil {
		fmt.Printf("Erro ao abrir arquivo: %v\n", err)
		os.Exit(1)
	}

	// Tabela de símbolos
	tabela := semantico.NovaTabelaDeSimbolos()

	// Gerador de código objeto
	gerador := gerador_codigo.NovoGerador(arquivoSaida)

	// Parser com lexer, tabela e gerador
	parser := sintatico.NovoParser(lex, tabela, gerador)

	// Analisa o programa e gera o código objeto
	err = parser.Programa()
	if err != nil {
		fmt.Printf("Erro na compilação: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Compilação concluída. Código objeto salvo em %s\n", arquivoSaida)
}
