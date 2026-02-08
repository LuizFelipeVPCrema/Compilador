package main

import (
	"fmt"
	"os"

	"ufmt.br/luiz-crema/compilador/executor"
	"ufmt.br/luiz-crema/compilador/gerador_codigo"
	"ufmt.br/luiz-crema/compilador/lexer"
	"ufmt.br/luiz-crema/compilador/semantico"
	"ufmt.br/luiz-crema/compilador/sintatico"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run . <arquivo_entrada> [lexemas|run]")
		fmt.Println("     go run . <codigo.objeto.txt> --exec  (apenas executar código objeto)")
		fmt.Println("Exemplo: go run . descricao/correto.python.txt")
		fmt.Println("         go run . descricao/correto.python.txt lexemas  (gera lexemas.txt)")
		fmt.Println("         go run . descricao/correto.python.txt run      (compila e executa)")
		fmt.Println("         go run . codigo.objeto.txt --exec              (só executa)")
		os.Exit(1)
	}

	nomeArquivo := os.Args[1]

	if len(os.Args) > 2 && os.Args[2] == "--exec" {
		vm := executor.NovaVM(os.Stdin, os.Stdout)
		if err := vm.Executar(nomeArquivo); err != nil {
			fmt.Fprintf(os.Stderr, "Erro na execução: %v\n", err)
			os.Exit(1)
		}
		return
	}

	arquivoSaida := "codigo.objeto.txt"
	gerarLexemas := len(os.Args) > 2 && os.Args[2] == "lexemas"
	executarDepois := len(os.Args) > 2 && os.Args[2] == "run"

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

	if executarDepois {
		fmt.Println("--- Executando Compilador ---")
		vm := executor.NovaVM(os.Stdin, os.Stdout)
		if err := vm.Executar(arquivoSaida); err != nil {
			fmt.Fprintf(os.Stderr, "Erro na execução: %v\n", err)
			os.Exit(1)
		}
	}
}
