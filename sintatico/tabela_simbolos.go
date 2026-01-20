package parser

import "fmt"

type TipoSimbolo int

const (
	SIMBOLO_VARIAVEL TipoSimbolo = iota
	SIMBOLO_FUNCAO
)

type Simbolo struct {
	Nome           string
	Tipo           TipoSimbolo
	QtdeParametros int // somente para funções
}

type TabelaDeSimbolos struct {
	simbolos map[string]Simbolo
}

func NovaTabelaDeSimbolos() *TabelaDeSimbolos {
	return &TabelaDeSimbolos{
		simbolos: make(map[string]Simbolo),
	}
}

func (tabelaDeSimbolos *TabelaDeSimbolos) DeclararVariavel(nome string) {
	tabelaDeSimbolos.simbolos[nome] = Simbolo{
		Nome: nome,
		Tipo: SIMBOLO_VARIAVEL,
	}
}

func (tabelaDeSimbolos *TabelaDeSimbolos) DeclararFuncao(nome string, params int) error {
	if _, exist := tabelaDeSimbolos.simbolos[nome]; exist {
		return fmt.Errorf("erro semântico: função '%s' já declarada", nome)
	}
	tabelaDeSimbolos.simbolos[nome] = Simbolo{
		Nome:           nome,
		Tipo:           SIMBOLO_FUNCAO,
		QtdeParametros: params,
	}
	return nil
}

func (tabelaDeSimbolos *TabelaDeSimbolos) Existe(nome string) bool {
	_, ok := tabelaDeSimbolos.simbolos[nome]
	return ok
}

func (tabelaDeSimbolos *TabelaDeSimbolos) Adicionar(nome string, tipo TipoSimbolo) error {
	if _, existe := tabelaDeSimbolos.simbolos[nome]; existe {
		return fmt.Errorf("erro semântico: '%s' já declarado", nome)
	}

	tabelaDeSimbolos.simbolos[nome] = Simbolo{
		Nome: nome,
		Tipo: tipo,
	}

	return nil
}

func (tabelaDeSimbolos *TabelaDeSimbolos) Listar() {
	for nome, simbolo := range tabelaDeSimbolos.simbolos {
		fmt.Printf("Nome: %s, Tipo: %v\n", nome, simbolo.Tipo)
	}
}

func (tabelaDeSimbolos *TabelaDeSimbolos) Buscar(nome string) (Simbolo, bool) {
	simbolo, existe := tabelaDeSimbolos.simbolos[nome]
	return simbolo, existe
}
