package sintatico

import (
	"fmt"

	"ufmt.br/luiz-crema/compilador/gerador_codigo"
	"ufmt.br/luiz-crema/compilador/lexer"
	"ufmt.br/luiz-crema/compilador/semantico"
)

type Parser struct {
	lexer              *lexer.Lexer
	proximoToken       lexer.Token
	tabelaDeSimbolos   *semantico.TabelaDeSimbolos
	gerador            *gerador_codigo.Gerador
	dentroDeFuncao     bool
	exigeIndentacao    bool
	indentacaoBloco    string
	indentacaoPendente string

	numParamsAtual  int
	numLocaisAtual  int
	indiceDSVI      int
	nomeFuncaoAtual string
	indicesDSVI     []int

	enderecoBaseFuncoes int
}

func NovoParser(lexer *lexer.Lexer, tabela *semantico.TabelaDeSimbolos, gerador *gerador_codigo.Gerador) *Parser {
	parser := &Parser{
		lexer:               lexer,
		tabelaDeSimbolos:    tabela,
		gerador:             gerador,
		enderecoBaseFuncoes: -1, // -1 = ainda não definido (será definido na primeira função)
	}
	parser.proximoToken = lexer.ProximoToken()
	return parser
}

func (parser *Parser) Avancar() {
	parser.proximoToken = parser.lexer.ProximoToken()
}

func (parser *Parser) Confirmar(tag lexer.Tag) error {
	if parser.proximoToken.Tag != tag {
		return fmt.Errorf("erro sintático: esperado %v, veio %v (%q)", tag, parser.proximoToken.Tag, parser.proximoToken.Lexeme)
	}
	parser.Avancar()
	return nil
}

// Programa: início, corpo (declarações + comandos) e fim.
func (parser *Parser) Programa() error {
	parser.gerador.InicioPrograma()
	if err := parser.Corpo(); err != nil {
		return err
	}
	parser.gerador.FimPrograma()
	return nil
}
