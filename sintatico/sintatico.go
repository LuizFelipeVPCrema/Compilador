package parser

import (
	"fmt"

	"ufmt.br/luiz-crema/compilador/lexer"
	"ufmt.br/luiz-crema/compilador/semantico"
)

type Parser struct {
	lexer            *lexer.Lexer
	lookahead        lexer.Token
	tabelaDeSimbolos *semantico.TabelaDeSimbolos
}

func NovoLexer(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer}
	parser.lookahead = lexer.ProximoToken()
	return parser
}

func (parser *Parser) Avancar() {
	parser.lookahead = parser.lexer.ProximoToken()
}

func (parser *Parser) Confirmar(tag lexer.Tag) error {
	if parser.lookahead.Tag != tag {
		return fmt.Errorf("error sintático: esperado %v, veio %v (%q)", tag, parser.lookahead.Tag, parser.lookahead.Lexeme)
	}
	parser.Avancar()
	return nil
}

func (parser *Parser) Comando() error {
	switch parser.lookahead.Tag {
	case lexer.TOKEN_VARIAVEL:
		return parser.AtribuicaoOuChamada()
	case lexer.TOKEN_PRINT:
		return parser.Imprimir()
	case lexer.TOKEN_IF:
		return parser.Condicional()
	case lexer.TOKEN_WHILE:
		return parser.While()
	case lexer.TOKEN_FOR:
		return parser.For()
	case lexer.TOKEN_READ:
		return parser.Read()
	case lexer.TOKEN_INPUT:
		return parser.Input()
	default:
		return fmt.Errorf("erro sintático: comando inesperado %v (%q)", parser.lookahead.Tag, parser.lookahead.Lexeme)
	}
}

func (parser *Parser) AtribuicaoOuChamada() error {

	nome := parser.lookahead.Lexeme

	if err := parser.tabelaDeSimbolos.VerificarVariavelDeclarada(nome); err != nil {
		return err
	}

	if parser.lookahead.Tag == lexer.TOKEN_ATRIBUICAO {
		parser.Avancar()
		return parser.Expressao()
	}

	if parser.lookahead.Tag == lexer.TOKEN_ABRE_PARENTESES {
		parser.Avancar()
		if parser.lookahead.Tag != lexer.TOKEN_FECHA_PARENTESES {
			if err := parser.Expressao(); err != nil {
				return err
			}
			for parser.lookahead.Tag == lexer.TOKEN_VIRGULA {
				parser.Avancar()
				if err := parser.Expressao(); err != nil {
					return err
				}
			}
		}
		return parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES)
	}

	return fmt.Errorf("erro sintático: esperado '=' ou '(' após uma variável")
}

func (parser *Parser) Expressao() error {
	if err := parser.Termo(); err != nil {
		return err
	}
	for parser.lookahead.Tag == lexer.TOKEN_MAIS || parser.lookahead.Tag == lexer.TOKEN_MENOS {
		parser.Avancar()
		if err := parser.Termo(); err != nil {
			return err
		}
	}
	return nil
}

func (parser *Parser) Termo() error {
	if err := parser.Fator(); err != nil {
		return err
	}
	for parser.lookahead.Tag == lexer.TOKEN_VEZES || parser.lookahead.Tag == lexer.TOKEN_DIVIDIDO {
		parser.Avancar()
		if err := parser.Fator(); err != nil {
			return err
		}
	}
	return nil
}

func (parser *Parser) Fator() error {
	switch parser.lookahead.Tag {
	case lexer.TOKEN_NUMERO:
		return parser.Confirmar(lexer.TOKEN_NUMERO)
	case lexer.TOKEN_VARIAVEL:
		return parser.Confirmar(lexer.TOKEN_VARIAVEL)
	case lexer.TOKEN_ABRE_PARENTESES:
		parser.Avancar()
		if err := parser.Expressao(); err != nil {
			return err
		}
		return parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES)
	default:
		return fmt.Errorf("erro sintático: fator inválido %v (%q)", parser.lookahead.Tag, parser.lookahead.Lexeme)
	}
}

// <relacao> -> == | != | >= | <= | > | <
func (parser *Parser) Relacao() error {
	switch parser.lookahead.Tag {
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
		return fmt.Errorf("erro sintático: operador de comparação inválido %v", parser.lookahead.Tag)
	}
}

// print(imprimido)
func (parser *Parser) Imprimir() error {
	if err := parser.Confirmar(lexer.TOKEN_PRINT); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES); err != nil {
		return err
	}
	if err := parser.Expressao(); err != nil {
		return err
	}
	return parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES)
}

//		if (condicao) :
//			bloco_if
//		else :
//	   	bloco_else
func (parser *Parser) Condicional() error {
	if err := parser.Confirmar(lexer.TOKEN_IF); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES); err != nil {
		return err
	}
	if err := parser.Expressao(); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
		return err
	}
	if err := parser.Bloco(); err != nil {
		return err
	}
	if parser.lookahead.Tag == lexer.TOKEN_ELSE {
		parser.Avancar()
		if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
			return err
		}
		return parser.Bloco()
	}
	return nil
}

func (parser *Parser) Bloco() error {
	if parser.lookahead.Tag == lexer.TOKEN_IDENTACAO {
		parser.Avancar()
		return parser.Comando()
	}
	return nil
}

// while (condicao):
//
//	bloco_while
func (parser *Parser) While() error {
	if err := parser.Confirmar(lexer.TOKEN_WHILE); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES); err != nil {
		return err
	}
	if err := parser.Expressao(); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
		return err
	}
	return parser.Bloco()
}

// def funcao(a, b):
//
//	bloco_funcao
func (parser *Parser) Funcao() error {
	if err := parser.Confirmar(lexer.TOKEN_DEF); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_VARIAVEL); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES); err != nil {
		return err
	}
	if err := parser.Parametros(); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
		return err
	}
	return parser.Bloco()
}

func (parser *Parser) Parametros() error {
	if parser.lookahead.Tag == lexer.TOKEN_VARIAVEL {
		if err := parser.Confirmar(lexer.TOKEN_VARIAVEL); err != nil {
			return err
		}
		for parser.lookahead.Tag == lexer.TOKEN_VIRGULA {
			parser.Avancar()
			if err := parser.Confirmar(lexer.TOKEN_VARIAVEL); err != nil {
				return err
			}
		}
	}
	return nil
}

func (parser *Parser) ChamadaFuncao() error {
	if err := parser.Confirmar(lexer.TOKEN_VARIAVEL); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES); err != nil {
		return err
	}
	if err := parser.Expressao(); err != nil {
		return err
	}
	for parser.lookahead.Tag == lexer.TOKEN_VIRGULA {
		parser.Avancar()
		if err := parser.Expressao(); err != nil {
			return err
		}
	}
	return parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES)
}

// for (condicao):
//
//	bloco
func (parser *Parser) For() error {
	if err := parser.Confirmar(lexer.TOKEN_FOR); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES); err != nil {
		return err
	}
	if err := parser.Expressao(); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_FECHA_PARENTESES); err != nil {
		return err
	}
	if err := parser.Confirmar(lexer.TOKEN_DOIS_PONTOS); err != nil {
		return err
	}

	return parser.Bloco()
}

// read()
func (parser *Parser) Read() error {
	if err := parser.Confirmar(lexer.TOKEN_READ); err != nil {
		return err
	}
	return parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES)
}

// input()
func (parser *Parser) Input() error {
	if err := parser.Confirmar(lexer.TOKEN_INPUT); err != nil {
		return err
	}
	return parser.Confirmar(lexer.TOKEN_ABRE_PARENTESES)
}
