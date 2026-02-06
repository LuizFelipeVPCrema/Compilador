package sintatico

import (
	"fmt"

	"ufmt.br/luiz-crema/compilador/lexer"
)

// <relacao> -> == | != | >= | <= | > | <
func (parser *Parser) Relacao() error {
	switch parser.proximoToken.Tag {
	case lexer.TOKEN_IGUAL:
		return parser.Confirmar(lexer.TOKEN_IGUAL)
	case lexer.TOKEN_DIFERENTE:
		return parser.Confirmar(lexer.TOKEN_DIFERENTE)
	case lexer.TOKEN_MAIOR_IGUAL_QUE:
		return parser.Confirmar(lexer.TOKEN_MAIOR_IGUAL_QUE)
	case lexer.TOKEN_MENOR_IGUAL_QUE:
		return parser.Confirmar(lexer.TOKEN_MENOR_IGUAL_QUE)
	case lexer.TOKEN_MAIOR_QUE:
		return parser.Confirmar(lexer.TOKEN_MAIOR_QUE)
	case lexer.TOKEN_MENOR_QUE:
		return parser.Confirmar(lexer.TOKEN_MENOR_QUE)
	default:
		return fmt.Errorf("erro sintático: operador de comparação inválido %v", parser.proximoToken.Tag)
	}
}

// Condicao: <condicao> -> <expressao> <relacao> <expressao>.
func (parser *Parser) Condicao() error {
	if err := parser.Expressao(); err != nil {
		return err
	}
	if err := parser.Relacao(); err != nil {
		return err
	}
	return parser.Expressao()
}

// emitirCondicao parse a condição e emite o código (expr1, expr2, comparação).
func (parser *Parser) emitirCondicao() error {
	if err := parser.Expressao(); err != nil {
		return err
	}
	relTag := parser.proximoToken.Tag
	if err := parser.Relacao(); err != nil {
		return err
	}
	if err := parser.Expressao(); err != nil {
		return err
	}
	parser.emitirComparacao(relTag)
	return nil
}

func (parser *Parser) emitirComparacao(relTag lexer.Tag) {
	switch relTag {
	case lexer.TOKEN_IGUAL:
		parser.gerador.GerarIgual()
	case lexer.TOKEN_DIFERENTE:
		parser.gerador.GerarDiferente()
	case lexer.TOKEN_MAIOR_IGUAL_QUE:
		parser.gerador.GerarMaiorIgual()
	case lexer.TOKEN_MENOR_IGUAL_QUE:
		parser.gerador.GerarMenorIgual()
	case lexer.TOKEN_MAIOR_QUE:
		parser.gerador.GerarMaior()
	case lexer.TOKEN_MENOR_QUE:
		parser.gerador.GerarMenor()
	}
}

// Condicional: condição (emitida), DSVF, bloco então, DSVI, [bloco senão].
func (parser *Parser) Condicional() error {
	if err := parser.Confirmar(lexer.TOKEN_IF); err != nil {
		return err
	}
	return parser.gerador.GerarIf(
		parser.emitirCondicao,
		func() error {
			if parser.proximoToken.Tag == lexer.TOKEN_FECHA_PARENTESES {
				parser.Avancar()
			}
			if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
				return err
			}
			return parser.Bloco()
		},
		func() error {
			if parser.proximoToken.Tag != lexer.TOKEN_ELSE {
				return nil
			}
			parser.Avancar()
			if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
				return err
			}
			return parser.Bloco()
		},
	)
}

// Bloco analisa: <bloco> -> tabulacao <comandos>
func (parser *Parser) Bloco() error {
	antigaIndentacaoBloco := parser.indentacaoBloco
	defer func() {
		parser.indentacaoBloco = antigaIndentacaoBloco
	}()

	parser.exigeIndentacao = true

	parser.indentacaoPendente = ""

	for parser.proximoToken.Tag == lexer.TOKEN_FIM_LINHA {
		parser.Avancar()
	}

	if parser.proximoToken.Tag != lexer.TOKEN_IDENTACAO {
		return nil
	}
	parser.indentacaoBloco = ""
	for parser.proximoToken.Tag == lexer.TOKEN_IDENTACAO {
		parser.indentacaoBloco += parser.proximoToken.Lexeme
		parser.Avancar()
	}
	return parser.Comandos()
}

// While: condição (emitida), DSVF, bloco, DSVI (início).
func (parser *Parser) While() error {
	if err := parser.Confirmar(lexer.TOKEN_WHILE); err != nil {
		return err
	}
	return parser.gerador.GerarWhile(
		parser.emitirCondicao,
		func() error {
			if parser.proximoToken.Tag == lexer.TOKEN_FECHA_PARENTESES {
				parser.Avancar()
			}
			if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
				return err
			}
			return parser.Bloco()
		},
	)
}
