package sintatico

import (
	"fmt"

	"ufmt.br/luiz-crema/compilador/lexer"
)

// Corpo: <corpo> -> <dc> <comandos>
func (parser *Parser) Corpo() error {
	if err := parser.DC(); err != nil {
		return err
	}
	linhaInicioMain := parser.gerador.ProximaLinha()
	for _, idx := range parser.indicesDSVI {
		parser.gerador.Preencher(idx, linhaInicioMain, "DSVI")
	}
	parser.indicesDSVI = nil
	return parser.Comandos()
}

// DC: <dc> -> <dc_v> <mais_dc> | <dc_f> | λ
func (parser *Parser) DC() error {
	for parser.proximoToken.Tag == lexer.TOKEN_FIM_LINHA {
		parser.Avancar()
	}
	if parser.proximoToken.Tag == lexer.TOKEN_INICIO_LINHA_SEM_INDENT {
		parser.Avancar()
	}
	switch parser.proximoToken.Tag {
	case lexer.TOKEN_VARIAVEL:
		_, encontrou := parser.tabelaDeSimbolos.Buscar(parser.proximoToken.Lexeme)
		if encontrou {
			return nil
		}
		return parser.dcVar()
	case lexer.TOKEN_DEF:
		return parser.dcFunc()
	default:
		return nil
	}
}

// dcVar: <dc_v> -> ident = expressao
func (parser *Parser) dcVar() error {
	nome := parser.proximoToken.Lexeme
	parser.Avancar()
	if parser.proximoToken.Tag != lexer.TOKEN_ATRIBUICAO {
		return fmt.Errorf("erro sintático: esperado '=' em declaração de variável, veio %v (%q)", parser.proximoToken.Tag, parser.proximoToken.Lexeme)
	}
	parser.tabelaDeSimbolos.DeclararVariavel(nome)
	parser.gerador.GerarAlocacao(1)
	parser.Confirmar(lexer.TOKEN_ATRIBUICAO)
	if parser.proximoToken.Tag == lexer.TOKEN_NUMERO && (parser.proximoToken.Lexeme == "0" || parser.proximoToken.Lexeme == "0.0") {
		parser.Avancar()
		return parser.MaisDC()
	}
	if err := parser.Expressao(); err != nil {
		return err
	}
	return parser.MaisDC()
}

// dcFunc: <dc_f> -> def ident parametros : corpo_f
func (parser *Parser) dcFunc() error {
	if parser.enderecoBaseFuncoes < 0 {
		parser.enderecoBaseFuncoes = parser.tabelaDeSimbolos.ProximoEndereco()
	}
	if err := parser.Funcao(); err != nil {
		return err
	}
	return parser.MaisDC()
}

// MaisDC: <mais_dc> -> <dc> | λ
func (parser *Parser) MaisDC() error {
	for parser.proximoToken.Tag == lexer.TOKEN_FIM_LINHA {
		parser.Avancar()
	}
	return parser.DC()
}
